#include <fcitx/addonfactory.h>
#include <fcitx/addonmanager.h>
#include <fcitx/event.h>
#include <fcitx/inputcontext.h>
#include <fcitx/inputmethodengine.h>
#include <fcitx/inputpanel.h>
#include <fcitx/instance.h>
#include <fcitx-utils/event.h>
#include <fcitx-utils/utf8.h>
#include <cstdlib>
#include <memory>
#include <string>
#include "libvuikey_core.h"

namespace fcitx {

class VuiKeyEngine : public InputMethodEngineV2 {
public:
    VuiKeyEngine(Instance *instance) : instance_(instance) {
        VuiCoreReloadConfig();

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
                flushPendingCommit();
                auto &fe = static_cast<InputContextEvent &>(event);
                commitAndClear(fe.inputContext(), 0);
            });

        resetWatcher_ = instance_->watchEvent(
            EventType::InputContextReset,
            EventWatcherPhase::PreInputMethod,
            [this](Event &event) {
                flushPendingCommit();
                auto &re = static_cast<InputContextEvent &>(event);
                auto *ic = re.inputContext();
                char *rem = VuiCoreOnClientReset();
                if (rem) free(rem);
                const auto &rect = ic->cursorRect();
                VuiCoreUpdateCursorRect(rect.left(), rect.top(), rect.height());
                ic->inputPanel().reset();
                ic->updatePreedit();
                ic->updateUserInterface(UserInterfaceComponent::InputPanel);
            });

        cursorRectWatcher_ = instance_->watchEvent(
            EventType::InputContextCursorRectChanged,
            EventWatcherPhase::PreInputMethod,
            [this](Event &event) {
                auto &ce = static_cast<InputContextEvent &>(event);
                auto *ic = ce.inputContext();
                const auto &rect = ic->cursorRect();
                VuiCoreUpdateCursorRect(rect.left(), rect.top(), rect.height());
            });
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
        flushPendingCommit();
        auto *ic = event.inputContext();
        commitAndClear(ic, 0);
        reset(entry, event);
    }

    void reset(const InputMethodEntry &, InputContextEvent &event) override {
        flushPendingCommit();
        auto *ic = event.inputContext();
        char *rem = VuiCoreOnClientReset();
        if (rem) free(rem);
        ic->inputPanel().reset();
        ic->updatePreedit();
        ic->updateUserInterface(UserInterfaceComponent::InputPanel);
    }

    void keyEvent(const InputMethodEntry &, KeyEvent &) override {
        // Handled once at EventWatcherPhase::PreInputMethod by keyWatcher_
    }

private:
    void flushPendingCommit() {
        pendingCommitTimer_.reset();
        if (!pendingCommitText_.empty()) {
            std::string toCommit = pendingCommitText_;
            pendingCommitText_.clear();
            if (auto *ic = pendingCommitIC_.get()) {
                ic->commitString(toCommit);
            }
        }
    }

    void handleKeyEvent(KeyEvent &keyEvent) {
        if (keyEvent.isRelease()) {
            return;
        }

        flushPendingCommit();

        auto *ic = keyEvent.inputContext();

        // Bypass VuiKey completely in Password or Sensitive input fields!
        if (ic->capabilityFlags().testAny(CapabilityFlag::PasswordOrSensitive)) {
            char *rem = VuiCoreRecordNavigation(0);
            if (rem) free(rem);
            return;
        }

        const auto &rect = ic->cursorRect();
        VuiCoreUpdateCursorRect(rect.left(), rect.top(), rect.height());

        Key key = keyEvent.key();
        Key raw = keyEvent.rawKey();

        // Toggle [VN] / [EN] directly in ~/.config/vuikey/config.json on Ctrl+Space or Super+Space
        if (key.check(Key(FcitxKey_space, KeyState::Ctrl)) ||
            key.check(Key(FcitxKey_space, KeyState::Super))) {
            commitAndClear(ic, 0);
            VuiCoreToggleEnabled();
            keyEvent.filterAndAccept();
            return;
        }

        // If Ctrl, Alt, or Super is pressed (e.g. Ctrl+A, Ctrl+Backspace, Ctrl+V, Ctrl+Z),
        // commit any active syllable and clear committedLine so selections/clipboard ops never leave stale state!
        if (key.states().testAny(KeyStates{KeyState::Ctrl_Alt, KeyState::Super})) {
            char *rem = VuiCoreRecordNavigation(0);
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

        KeySym sym = raw.sym();

        // 1. Handle BackSpace inside active syllable (0.0ms, zero race condition)
        if (sym == FcitxKey_BackSpace) {
            if (key.states().test(KeyState::Shift)) {
                commitAndClear(ic, 0);
                return;
            }
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

        // 2. Handle Delete, Return, Tab, Escape, and Cursor Movement keys:
        if (sym == FcitxKey_Delete || sym == FcitxKey_KP_Delete ||
            sym == FcitxKey_Return || sym == FcitxKey_KP_Enter ||
            sym == FcitxKey_Tab || sym == FcitxKey_ISO_Left_Tab ||
            sym == FcitxKey_Escape || key.isCursorMove()) {
            int navType = 0;
            // If Shift is held (e.g. Shift+Left/Right/Home/End text selection), clear committedLine (navType = 0)
            if (!key.states().test(KeyState::Shift)) {
                if (sym == FcitxKey_Left || sym == FcitxKey_KP_Left) navType = 1;
                else if (sym == FcitxKey_Right || sym == FcitxKey_KP_Right) navType = 2;
                else if (sym == FcitxKey_Home || sym == FcitxKey_KP_Home) navType = 3;
                else if (sym == FcitxKey_End || sym == FcitxKey_KP_End) navType = 4;
                else if (sym == FcitxKey_Delete || sym == FcitxKey_KP_Delete) navType = 5;
            }

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
            // 3a. If active syllable is empty (user moved cursor via arrow keys, mouse click, or backspace to the end of a word),
            // check if `ucs4` is a tone/diacritic modifier that validly transforms the word right before the cursor!
            bool hasValidSur = ic->capabilityFlags().test(CapabilityFlag::SurroundingText) &&
                               ic->surroundingText().isValid();
            const char *surCStr = hasValidSur ? ic->surroundingText().text().c_str() : nullptr;
            unsigned int curPos = hasValidSur ? ic->surroundingText().cursor() : 0;
            unsigned int anchorPos = hasValidSur ? ic->surroundingText().anchor() : 0;

            int surBs = 0;
            char *cSurCommit = nullptr;
            if (VuiCoreTryModifySurrounding(hasValidSur ? 1 : 0, const_cast<char *>(surCStr), curPos, anchorPos, ucs4, &surBs, &cSurCommit)) {
                std::string surCommitStr = cSurCommit ? cSurCommit : "";
                if (cSurCommit) free(cSurCommit);

                if (hasValidSur) {
                    if (surBs > 0) {
                        ic->deleteSurroundingText(-surBs, surBs);
                    }
                    if (!surCommitStr.empty()) {
                        ic->commitString(surCommitStr);
                    }
                } else {
                    // For apps without SurroundingText (e.g. Chromium/Electron on X11):
                    // Forward `surBs` BackSpace keys via Fcitx5's native `ic->forwardKey` (tagged with FCITX_FORWARD_MASK
                    // by im-fcitx5.so / im-ibus.so so they never loop back to Fcitx5 and never touch XTestFakeKeyEvent),
                    // then commit `surCommitStr` on Fcitx5's event loop after the client's GDK queue processes the BackSpaces!
                    for (int i = 0; i < surBs; i++) {
                        ic->forwardKey(Key(FcitxKey_BackSpace));
                    }
                    if (!surCommitStr.empty()) {
                        pendingCommitIC_ = ic->watch();
                        pendingCommitText_ += surCommitStr;
                        pendingCommitTimer_ = instance_->eventLoop().addTimeEvent(
                            CLOCK_MONOTONIC,
                            now(CLOCK_MONOTONIC) + 8000,
                            1,
                            [this](EventSourceTime *, uint64_t) {
                                flushPendingCommit();
                                return true;
                            });
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
    TrackableObjectReference<InputContext> pendingCommitIC_;
    std::string pendingCommitText_;
    std::unique_ptr<EventSourceTime> pendingCommitTimer_;
    std::unique_ptr<HandlerTableEntry<EventHandler>> keyWatcher_;
    std::unique_ptr<HandlerTableEntry<EventHandler>> focusOutWatcher_;
    std::unique_ptr<HandlerTableEntry<EventHandler>> resetWatcher_;
    std::unique_ptr<HandlerTableEntry<EventHandler>> cursorRectWatcher_;
};

class VuiKeyFactory : public AddonFactory {
public:
    AddonInstance *create(AddonManager *manager) override {
        return new VuiKeyEngine(manager->instance());
    }
};

} // namespace fcitx

FCITX_ADDON_FACTORY(fcitx::VuiKeyFactory);
