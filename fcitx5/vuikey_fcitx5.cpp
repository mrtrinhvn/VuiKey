#include <fcitx/addonfactory.h>
#include <fcitx/addonmanager.h>
#include <fcitx/event.h>
#include <fcitx/inputcontext.h>
#include <fcitx/inputmethodengine.h>
#include <fcitx/inputpanel.h>
#include <fcitx/instance.h>
#include <fcitx-utils/utf8.h>
#include <X11/Xlib.h>
#include <X11/keysym.h>
#include <X11/extensions/XTest.h>
#include <cstdlib>
#include <memory>
#include <string>
#include "libvuikey_core.h"

namespace fcitx {

class VuiKeyEngine : public InputMethodEngineV2 {
public:
    VuiKeyEngine(Instance *instance) : instance_(instance) {
        VuiCoreReloadConfig();
        x11Dpy_ = XOpenDisplay(nullptr);
        if (x11Dpy_) {
            bsKeycode_ = XKeysymToKeycode(x11Dpy_, XK_BackSpace);
        }

        // Intercept key events at PreInputMethod phase so VuiKey works across 100% of windows
        // even if Fcitx5's per-window state was left on keyboard-us!
        keyWatcher_ = instance_->watchEvent(
            EventType::InputContextKeyEvent,
            EventWatcherPhase::PreInputMethod,
            [this](Event &event) {
                auto &ke = static_cast<KeyEvent &>(event);
                if (!ke.accepted()) {
                    handleKeyEvent(ke);
                }
            });

        focusOutWatcher_ = instance_->watchEvent(
            EventType::InputContextFocusOut,
            EventWatcherPhase::PreInputMethod,
            [this](Event &event) {
                auto &fe = static_cast<InputContextEvent &>(event);
                commitAndClear(fe.inputContext(), 0);
            });

        resetWatcher_ = instance_->watchEvent(
            EventType::InputContextReset,
            EventWatcherPhase::PreInputMethod,
            [this](Event &event) {
                auto &re = static_cast<InputContextEvent &>(event);
                auto *ic = re.inputContext();
                char *rem = VuiCoreReset();
                if (rem) free(rem);
                ic->inputPanel().reset();
                ic->updatePreedit();
                ic->updateUserInterface(UserInterfaceComponent::InputPanel);
            });
    }

    ~VuiKeyEngine() override {
        if (x11Dpy_) {
            XCloseDisplay(x11Dpy_);
            x11Dpy_ = nullptr;
        }
    }

    std::vector<InputMethodEntry> listInputMethods() override {
        std::vector<InputMethodEntry> result;
        result.emplace_back("vuikey", "VuiKey (Tiếng Việt)", "vi", "vuikey");
        result.back().setLabel("VN").setIcon("vuikey");
        return result;
    }

    void activate(const InputMethodEntry &, InputContextEvent &) override {
        VuiCoreReloadConfig();
    }

    void deactivate(const InputMethodEntry &entry, InputContextEvent &event) override {
        auto *ic = event.inputContext();
        commitAndClear(ic, 0);
        reset(entry, event);
    }

    void reset(const InputMethodEntry &, InputContextEvent &event) override {
        auto *ic = event.inputContext();
        char *rem = VuiCoreReset();
        if (rem) free(rem);
        ic->inputPanel().reset();
        ic->updatePreedit();
        ic->updateUserInterface(UserInterfaceComponent::InputPanel);
    }

    void keyEvent(const InputMethodEntry &, KeyEvent &) override {
        // Handled once at EventWatcherPhase::PreInputMethod by keyWatcher_
    }

private:
    void handleKeyEvent(KeyEvent &keyEvent) {
        auto *ic = keyEvent.inputContext();
        Key key = keyEvent.key();
        Key raw = keyEvent.rawKey();
        KeySym sym = raw.sym();

        // XTest BackSpace + TriggerKey Round-Trip Barrier:
        // 1. Let the `pendingBackspaces_` synthetic BackSpace KeyPress/KeyRelease events pass straight through
        //    unaccepted so the client application deletes `surBs` committed characters before the cursor.
        // 2. Right after those BackSpaces finish, the synthetic `pendingTriggerSym_` KeyPress arrives here:
        //    we swallow it and commit `pendingCommit_` at the exact moment all BackSpaces have taken effect!
        if (pendingBackspaces_ > 0 && sym == FcitxKey_BackSpace) {
            if (!keyEvent.isRelease()) {
                pendingBackspaces_--;
                if (pendingBackspaces_ == 0 && pendingTriggerKc_ != 0 && x11Dpy_) {
                    // Only inject the trigger key AFTER the final BackSpace KeyPress has been seen by Fcitx5,
                    // so the client's GDK/X11 loop processes all BackSpace replies before seeing TriggerKey!
                    XTestFakeKeyEvent(x11Dpy_, pendingTriggerKc_, True, CurrentTime);
                    XTestFakeKeyEvent(x11Dpy_, pendingTriggerKc_, False, CurrentTime);
                    XFlush(x11Dpy_);
                    pendingTriggerKc_ = 0;
                }
            }
            return; // Pass BackSpace to the client app so it deletes 1 character
        }
        if (pendingTriggerSym_ != FcitxKey_None && sym == pendingTriggerSym_) {
            if (!keyEvent.isRelease()) {
                pendingTriggerSym_ = FcitxKey_None;
                std::string toCommit = pendingCommit_;
                pendingCommit_.clear();
                if (!toCommit.empty()) {
                    ic->commitString(toCommit);
                }
            }
            keyEvent.filterAndAccept();
            return;
        }

        if (keyEvent.isRelease()) {
            return;
        }

        // Bypass VuiKey completely in Password or Sensitive input fields!
        if (ic->capabilityFlags().testAny(CapabilityFlag::PasswordOrSensitive)) {
            char *rem = VuiCoreRecordNavigation(0);
            if (rem) free(rem);
            return;
        }

        // Toggle [VN] / [EN] directly in ~/.config/vuikey/config.json on Ctrl+Space or Super+Space
        if (key.check(Key(FcitxKey_space, KeyState::Ctrl)) ||
            key.check(Key(FcitxKey_space, KeyState::Super))) {
            commitAndClear(ic, 0);
            VuiCoreToggleEnabled();
            keyEvent.filterAndAccept();
            return;
        }

        // If Ctrl, Alt, or Super is pressed, commit any active syllable first
        if (key.states().testAny(KeyStates{KeyState::Ctrl_Alt, KeyState::Super})) {
            commitAndClear(ic, 0);
            return;
        }

        // 1. Handle BackSpace inside active syllable (0.0ms, zero race condition)
        if (sym == FcitxKey_BackSpace) {
            int handled = 0;
            char *newPreedit = VuiCoreHandleBackspace(&handled);
            std::string pre = newPreedit ? newPreedit : "";
            if (newPreedit) free(newPreedit);
            if (handled) {
                updateClientPreedit(ic, pre);
                keyEvent.filterAndAccept();
            }
            return;
        }

        // 2. Handle Delete, Return, Tab, Escape, and Cursor Movement keys
        if (sym == FcitxKey_Delete || sym == FcitxKey_KP_Delete ||
            sym == FcitxKey_Return || sym == FcitxKey_KP_Enter ||
            sym == FcitxKey_Tab || sym == FcitxKey_ISO_Left_Tab ||
            sym == FcitxKey_Escape || key.isCursorMove()) {
            int navType = 0;
            if (sym == FcitxKey_Left || sym == FcitxKey_KP_Left) navType = 1;
            else if (sym == FcitxKey_Right || sym == FcitxKey_KP_Right) navType = 2;
            else if (sym == FcitxKey_Home || sym == FcitxKey_KP_Home) navType = 3;
            else if (sym == FcitxKey_End || sym == FcitxKey_KP_End) navType = 4;
            else if (sym == FcitxKey_Delete || sym == FcitxKey_KP_Delete) navType = 5;

            char *rem = VuiCoreRecordNavigation(navType);
            std::string syl = rem ? rem : "";
            if (rem) free(rem);
            if (!syl.empty()) {
                ic->inputPanel().reset();
                ic->updatePreedit();
                ic->commitString(syl);
                ic->updateUserInterface(UserInterfaceComponent::InputPanel);
                ic->forwardKey(raw);
                keyEvent.filterAndAccept();
            }
            return;
        }

        // Ignore pure modifier keys (Shift, CapsLock, etc.)
        if (key.isModifier()) {
            return;
        }

        // 3. Printable ASCII character
        uint32_t ucs4 = Key::keySymToUnicode(sym);
        if (ucs4 >= 0x20 && ucs4 <= 0x7e) {
            // 3a. If active syllable is empty (user moved cursor or backspaced to the end of an existing word),
            // check if `ucs4` is a tone/diacritic modifier that can validly transform the word right before the cursor!
            bool hasValidSur = ic->capabilityFlags().test(CapabilityFlag::SurroundingText) &&
                               ic->surroundingText().isValid();
            const char *surCStr = hasValidSur ? ic->surroundingText().text().c_str() : nullptr;
            unsigned int curPos = hasValidSur ? ic->surroundingText().cursor() : 0;

            int surBs = 0;
            char *cSurCommit = nullptr;
            if (VuiCoreTryModifySurrounding(const_cast<char *>(surCStr), curPos, ucs4, &surBs, &cSurCommit)) {
                std::string surCommitStr = cSurCommit ? cSurCommit : "";
                if (cSurCommit) free(cSurCommit);

                if (x11Dpy_ && bsKeycode_ != 0 && surBs > 0) {
                    // On X11, use the deterministic XTest [BackSpace * surBs + TriggerKey] Round-Trip Barrier!
                    // Works across 100% of X11 apps (Chromium, Electron, Antigravity, VS Code, Chrome, GTK, Qt)
                    // even when the app does not implement GtkIMContext delete-surrounding!
                    unsigned int trigKc = XKeysymToKeycode(x11Dpy_, (KeySym)sym);
                    if (trigKc != 0) {
                        pendingBackspaces_ = surBs;
                        pendingTriggerSym_ = sym;
                        pendingTriggerKc_ = trigKc;
                        pendingCommit_ = surCommitStr;
                        for (int i = 0; i < surBs; i++) {
                            XTestFakeKeyEvent(x11Dpy_, bsKeycode_, True, CurrentTime);
                            XTestFakeKeyEvent(x11Dpy_, bsKeycode_, False, CurrentTime);
                        }
                        XFlush(x11Dpy_);
                    }
                } else if (hasValidSur) {
                    if (surBs > 0) {
                        ic->deleteSurroundingText(-surBs, surBs);
                    }
                    if (!surCommitStr.empty()) {
                        ic->commitString(surCommitStr);
                    }
                }
                keyEvent.filterAndAccept();
                return;
            }

            char *cCommit = nullptr;
            char *cPreedit = nullptr;
            int handled = VuiCoreProcessKey(ucs4, &cCommit, &cPreedit);
            std::string commitStr = cCommit ? cCommit : "";
            std::string preeditStr = cPreedit ? cPreedit : "";
            if (cCommit) free(cCommit);
            if (cPreedit) free(cPreedit);

            if (handled) {
                if (!commitStr.empty() && preeditStr.empty()) {
                    ic->inputPanel().reset();
                    ic->updatePreedit();
                    ic->commitString(commitStr);
                    ic->updateUserInterface(UserInterfaceComponent::InputPanel);
                } else {
                    if (!commitStr.empty()) {
                        ic->commitString(commitStr);
                    }
                    updateClientPreedit(ic, preeditStr);
                }
                keyEvent.filterAndAccept();
            }
            return;
        }
    }

    void commitAndClear(InputContext *ic, int navType) {
        char *rem = VuiCoreRecordNavigation(navType);
        std::string syl = rem ? rem : "";
        if (rem) free(rem);
        if (!syl.empty()) {
            ic->inputPanel().reset();
            ic->updatePreedit();
            ic->commitString(syl);
            ic->updateUserInterface(UserInterfaceComponent::InputPanel);
        }
    }

    void updateClientPreedit(InputContext *ic, const std::string &pre) {
        ic->inputPanel().reset();
        if (!pre.empty()) {
            Text preeditText(pre, TextFormatFlag::NoFlag);
            preeditText.setCursor(pre.size());
            if (ic->capabilityFlags().test(CapabilityFlag::Preedit)) {
                ic->inputPanel().setClientPreedit(preeditText);
            } else {
                ic->inputPanel().setPreedit(preeditText);
            }
        }
        ic->updatePreedit();
        ic->updateUserInterface(UserInterfaceComponent::InputPanel);
    }

    Instance *instance_;
    Display *x11Dpy_ = nullptr;
    unsigned int bsKeycode_ = 0;
    int pendingBackspaces_ = 0;
    KeySym pendingTriggerSym_ = FcitxKey_None;
    unsigned int pendingTriggerKc_ = 0;
    std::string pendingCommit_;
    std::unique_ptr<HandlerTableEntry<EventHandler>> keyWatcher_;
    std::unique_ptr<HandlerTableEntry<EventHandler>> focusOutWatcher_;
    std::unique_ptr<HandlerTableEntry<EventHandler>> resetWatcher_;
};

class VuiKeyFactory : public AddonFactory {
public:
    AddonInstance *create(AddonManager *manager) override {
        return new VuiKeyEngine(manager->instance());
    }
};

} // namespace fcitx

FCITX_ADDON_FACTORY(fcitx::VuiKeyFactory);
