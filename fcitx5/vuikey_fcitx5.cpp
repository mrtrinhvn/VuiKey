#include <fcitx/addonfactory.h>
#include <fcitx/addonmanager.h>
#include <fcitx/inputcontext.h>
#include <fcitx/inputmethodengine.h>
#include <fcitx/inputpanel.h>
#include <fcitx/instance.h>
#include <fcitx-utils/utf8.h>
#include <cstdlib>
#include <string>
#include "libvuikey_core.h"

namespace fcitx {

class VuiKeyEngine : public InputMethodEngineV2 {
public:
    VuiKeyEngine(Instance *instance) : instance_(instance) {
        VuiCoreReloadConfig();
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
        commitAndClear(ic);
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

    void keyEvent(const InputMethodEntry &, KeyEvent &keyEvent) override {
        if (keyEvent.isRelease()) {
            return;
        }

        auto *ic = keyEvent.inputContext();
        Key key = keyEvent.key();
        Key raw = keyEvent.rawKey();

        // If Ctrl, Alt, or Super is pressed, commit any active syllable first
        if (key.states().testAny(KeyStates{KeyState::Ctrl_Alt, KeyState::Super})) {
            commitAndClear(ic);
            return;
        }

        KeySym sym = raw.sym();

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

        // 2. Handle Delete, Return, Tab, Escape, and Cursor Movement keys:
        //    CRITICAL ANTI-DUPLICATE FIX:
        //    When an editing/navigation key is pressed while a syllable is active,
        //    we clear ClientPreedit FIRST, commit the syllable, and serialize the key via forwardKey()
        //    while calling filterAndAccept() so X11 and DBus NEVER race!
        if (sym == FcitxKey_Delete || sym == FcitxKey_KP_Delete ||
            sym == FcitxKey_Return || sym == FcitxKey_KP_Enter ||
            sym == FcitxKey_Tab || sym == FcitxKey_ISO_Left_Tab ||
            sym == FcitxKey_Escape || key.isCursorMove()) {
            char *rem = VuiCoreReset();
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
            char *cCommit = nullptr;
            char *cPreedit = nullptr;
            int handled = VuiCoreProcessKey(ucs4, &cCommit, &cPreedit);
            std::string commitStr = cCommit ? cCommit : "";
            std::string preeditStr = cPreedit ? cPreedit : "";
            if (cCommit) free(cCommit);
            if (cPreedit) free(cPreedit);

            if (handled) {
                if (!commitStr.empty() && preeditStr.empty()) {
                    // Word completed (e.g. Space or punctuation): clear preedit FIRST, then commit atomically!
                    ic->inputPanel().reset();
                    ic->updatePreedit();
                    ic->commitString(commitStr);
                    ic->updateUserInterface(UserInterfaceComponent::InputPanel);
                } else {
                    if (!commitStr.empty()) {
                        // Sliding Syllable Window boundary (e.g. nhac|di|dang): commit finished syllable immediately!
                        ic->commitString(commitStr);
                    }
                    updateClientPreedit(ic, preeditStr);
                }
                keyEvent.filterAndAccept();
            }
            return;
        }
    }

private:
    void commitAndClear(InputContext *ic) {
        char *rem = VuiCoreReset();
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
            // Use TextFormatFlag::NoFlag (no underline!) for a clean UniKey-like feel
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
};

class VuiKeyFactory : public AddonFactory {
public:
    AddonInstance *create(AddonManager *manager) override {
        return new VuiKeyEngine(manager->instance());
    }
};

} // namespace fcitx

FCITX_ADDON_FACTORY(fcitx::VuiKeyFactory);
