package ui

/*
#cgo pkg-config: gtk+-3.0 ayatana-appindicator3-0.1
#include <gtk/gtk.h>
#include <libayatana-appindicator/app-indicator.h>
#include <stdlib.h>

extern void goOnToggleClicked();
extern void goOnModeChanged(int mode);
extern void goOnFeatureToggled(int feature, int active);
extern void goOnNextSampleClicked();
extern void goOnOpenURLClicked(int target);
extern void goOnQuitClicked();

static AppIndicator *indicator = NULL;
static GtkWidget *win = NULL;
static GtkWidget *status_btn = NULL;
static GtkWidget *mode_combo = NULL;
static GtkWidget *sample_cat_lbl = NULL;
static GtkWidget *sample_quote_lbl = NULL;
static GtkWidget *sample_hint_lbl = NULL;
static GtkWidget *sample_link_btn = NULL;
static GtkWidget *test_entry = NULL;
static GtkWidget *stats_lbl = NULL;

static void on_toggle_btn(GtkWidget *w, gpointer data) {
    goOnToggleClicked();
}

static void on_mode_combo_changed(GtkComboBox *cb, gpointer data) {
    int idx = gtk_combo_box_get_active(cb);
    goOnModeChanged(idx);
}

static void on_split_toggled(GtkToggleButton *tb, gpointer data) {
    goOnFeatureToggled(1, gtk_toggle_button_get_active(tb));
}

static void on_eng_toggled(GtkToggleButton *tb, gpointer data) {
    goOnFeatureToggled(2, gtk_toggle_button_get_active(tb));
}

static void on_next_sample_btn(GtkWidget *w, gpointer data) {
    goOnNextSampleClicked();
}

static void on_clear_test_btn(GtkWidget *w, gpointer data) {
    if (test_entry) {
        gtk_entry_set_text(GTK_ENTRY(test_entry), "");
        gtk_widget_grab_focus(test_entry);
    }
}

static void on_open_url_0(GtkWidget *w, gpointer data) { goOnOpenURLClicked(0); }
static void on_open_url_1(GtkWidget *w, gpointer data) { goOnOpenURLClicked(1); }
static void on_open_url_2(GtkWidget *w, gpointer data) { goOnOpenURLClicked(2); }

static void on_show_window(GtkMenuItem *item, gpointer data) {
    if (win) {
        gtk_widget_show_all(win);
        gtk_window_present(GTK_WINDOW(win));
    }
}

static void on_tray_vn_telex(GtkMenuItem *item, gpointer data) {
    goOnModeChanged(0);
}

static void on_tray_vn_vni(GtkMenuItem *item, gpointer data) {
    goOnModeChanged(1);
}

static void on_tray_en(GtkMenuItem *item, gpointer data) {
    goOnModeChanged(99);
}

static void on_quit_menu(GtkMenuItem *item, gpointer data) {
    goOnQuitClicked();
    gtk_main_quit();
}

static gboolean on_win_delete(GtkWidget *w, GdkEvent *e, gpointer data) {
    gtk_widget_hide(w);
    return TRUE;
}

static void apply_dark_css() {
    GtkCssProvider *provider = gtk_css_provider_new();
    const char *css =
        "window { background-color: #0f172a; color: #f8fafc; font-family: 'Inter', 'Ubuntu', sans-serif; }"
        ".card { background-color: #1e293b; border-radius: 12px; padding: 14px; border: 1px solid #334155; }"
        ".sandbox-card { background-color: #172554; border-radius: 12px; padding: 14px; border: 1px solid #1d4ed8; }"
        ".eco-card { background-color: #1e1b4b; border-radius: 10px; padding: 12px; border: 1px solid #4338ca; }"
        ".title-lbl { font-size: 18px; font-weight: bold; color: #38bdf8; }"
        ".subtitle-lbl { font-size: 12px; color: #94a3b8; }"
        ".quote-lbl { font-size: 14px; font-weight: bold; color: #fde047; }"
        ".hint-lbl { font-size: 12px; color: #cbd5e1; }"
        ".btn-eco { background: #2563eb; color: #ffffff; font-weight: bold; border-radius: 8px; padding: 6px 12px; }"
        "entry { background-color: #090d16; color: #38bdf8; font-size: 15px; padding: 8px; border-radius: 8px; border: 1px solid #3b82f6; }";
    gtk_css_provider_load_from_data(provider, css, -1, NULL);
    gtk_style_context_add_provider_for_screen(
        gdk_screen_get_default(),
        GTK_STYLE_PROVIDER(provider),
        GTK_STYLE_PROVIDER_PRIORITY_APPLICATION
    );
    g_object_unref(provider);
}

static void init_ui(const char *icon_dir, int show_window_initially, int enabled, int mode, int split_on, int eng_on) {
    gtk_init(0, NULL);
    apply_dark_css();

    // 1. System Tray AppIndicator with custom high-contrast [VN] / [EN] badge icons
    indicator = app_indicator_new_with_path(
        "vuikey-indicator",
        enabled ? "vuikey-vn" : "vuikey-en",
        APP_INDICATOR_CATEGORY_APPLICATION_STATUS,
        icon_dir
    );
    app_indicator_set_status(indicator, APP_INDICATOR_STATUS_ACTIVE);
    app_indicator_set_title(indicator, "VuiKey Vietnamese Input");
    app_indicator_set_label(indicator, enabled ? "VN" : "EN", "VN");

    GtkWidget *menu = gtk_menu_new();
    GtkWidget *item_show = gtk_menu_item_new_with_label("⚙️ Mở Bảng Điều Khiển VuiKey Studio...");
    GtkWidget *item_vn_telex = gtk_menu_item_new_with_label("🇻🇳 Bật Tiếng Việt — Kiểu gõ Telex (Zero-Preedit)");
    GtkWidget *item_vn_vni = gtk_menu_item_new_with_label("🇻🇳 Bật Tiếng Việt — Kiểu gõ VNI (Zero-Preedit)");
    GtkWidget *item_en = gtk_menu_item_new_with_label("🇺🇸 Tạm Tắt Tiếng Việt — Chế độ Tiếng Anh (Ctrl+Space)");
    GtkWidget *sep1 = gtk_separator_menu_item_new();
    GtkWidget *item_vch = gtk_menu_item_new_with_label("🎤 Phòng hát & Thu âm trực tuyến — VuiCaHat.com");
    GtkWidget *item_echip = gtk_menu_item_new_with_label("🏢 Công ty TNHH eChip — echipco.com");
    GtkWidget *sep2 = gtk_separator_menu_item_new();
    GtkWidget *item_quit = gtk_menu_item_new_with_label("❌ Thoát VuiKey");

    g_signal_connect(item_show, "activate", G_CALLBACK(on_show_window), NULL);
    g_signal_connect(item_vn_telex, "activate", G_CALLBACK(on_tray_vn_telex), NULL);
    g_signal_connect(item_vn_vni, "activate", G_CALLBACK(on_tray_vn_vni), NULL);
    g_signal_connect(item_en, "activate", G_CALLBACK(on_tray_en), NULL);
    g_signal_connect(item_vch, "activate", G_CALLBACK(on_open_url_1), NULL);
    g_signal_connect(item_echip, "activate", G_CALLBACK(on_open_url_2), NULL);
    g_signal_connect(item_quit, "activate", G_CALLBACK(on_quit_menu), NULL);

    gtk_menu_shell_append(GTK_MENU_SHELL(menu), item_show);
    gtk_menu_shell_append(GTK_MENU_SHELL(menu), item_vn_telex);
    gtk_menu_shell_append(GTK_MENU_SHELL(menu), item_vn_vni);
    gtk_menu_shell_append(GTK_MENU_SHELL(menu), item_en);
    gtk_menu_shell_append(GTK_MENU_SHELL(menu), sep1);
    gtk_menu_shell_append(GTK_MENU_SHELL(menu), item_vch);
    gtk_menu_shell_append(GTK_MENU_SHELL(menu), item_echip);
    gtk_menu_shell_append(GTK_MENU_SHELL(menu), sep2);
    gtk_menu_shell_append(GTK_MENU_SHELL(menu), item_quit);
    gtk_widget_show_all(menu);
    app_indicator_set_menu(indicator, GTK_MENU(menu));

    // 2. Main Control Studio Window
    win = gtk_window_new(GTK_WINDOW_TOPLEVEL);
    gtk_window_set_title(GTK_WINDOW(win), "VuiKey Studio — Bộ Gõ Tiếng Việt Zero-Preedit");
    gtk_window_set_default_size(GTK_WINDOW(win), 650, 520);
    gtk_window_set_position(GTK_WINDOW(win), GTK_WIN_POS_CENTER);
    g_signal_connect(win, "delete-event", G_CALLBACK(on_win_delete), NULL);

    GtkWidget *root_box = gtk_box_new(GTK_ORIENTATION_VERTICAL, 12);
    gtk_container_set_border_width(GTK_CONTAINER(root_box), 16);
    gtk_container_add(GTK_CONTAINER(win), root_box);

    // Header
    GtkWidget *hdr_box = gtk_box_new(GTK_ORIENTATION_VERTICAL, 4);
    GtkWidget *title = gtk_label_new("⚡ VuiKey — Bộ Gõ Tiếng Việt Zero-Preedit (Không Gạch Chân)");
    gtk_style_context_add_class(gtk_widget_get_style_context(title), "title-lbl");
    gtk_widget_set_halign(title, GTK_ALIGN_START);
    GtkWidget *subtitle = gtk_label_new("Triệt tiêu 100% lỗi nhân đôi chữ khi bấm Delete/Backspace • Tự động tách âm tiết chuỗi dính liền");
    gtk_style_context_add_class(gtk_widget_get_style_context(subtitle), "subtitle-lbl");
    gtk_widget_set_halign(subtitle, GTK_ALIGN_START);
    gtk_box_pack_start(GTK_BOX(hdr_box), title, FALSE, FALSE, 0);
    gtk_box_pack_start(GTK_BOX(hdr_box), subtitle, FALSE, FALSE, 0);
    gtk_box_pack_start(GTK_BOX(root_box), hdr_box, FALSE, FALSE, 0);

    // Card 1: Core Engine Controls
    GtkWidget *card1 = gtk_box_new(GTK_ORIENTATION_VERTICAL, 10);
    gtk_style_context_add_class(gtk_widget_get_style_context(card1), "card");

    GtkWidget *row1 = gtk_box_new(GTK_ORIENTATION_HORIZONTAL, 12);
    status_btn = gtk_button_new_with_label(enabled ? "🟢 ĐANG BẬT TIẾNG VIỆT [VN] (Bấm hoặc Ctrl+Space để đổi)" : "⚪ ĐANG Ở CHẾ ĐỘ TIẾNG ANH [EN] (Bấm để bật Tiếng Việt)");
    g_signal_connect(status_btn, "clicked", G_CALLBACK(on_toggle_btn), NULL);

    mode_combo = gtk_combo_box_text_new();
    gtk_combo_box_text_append_text(GTK_COMBO_BOX_TEXT(mode_combo), "Kiểu gõ: Telex (Chuẩn)");
    gtk_combo_box_text_append_text(GTK_COMBO_BOX_TEXT(mode_combo), "Kiểu gõ: VNI (Số 1-9)");
    gtk_combo_box_text_append_text(GTK_COMBO_BOX_TEXT(mode_combo), "Kiểu gõ: Telex + VNI kết hợp");
    gtk_combo_box_set_active(GTK_COMBO_BOX(mode_combo), mode);
    g_signal_connect(mode_combo, "changed", G_CALLBACK(on_mode_combo_changed), NULL);

    gtk_box_pack_start(GTK_BOX(row1), status_btn, TRUE, TRUE, 0);
    gtk_box_pack_start(GTK_BOX(row1), mode_combo, FALSE, FALSE, 0);
    gtk_box_pack_start(GTK_BOX(card1), row1, FALSE, FALSE, 0);

    GtkWidget *chk_split = gtk_check_button_new_with_label("✔ Cửa sổ âm tiết trượt (Sliding Syllable Window): Tự tách nhac|di|dang|tieng..., chống nhân đôi chữ 100%");
    gtk_toggle_button_set_active(GTK_TOGGLE_BUTTON(chk_split), split_on);
    g_signal_connect(chk_split, "toggled", G_CALLBACK(on_split_toggled), NULL);

    GtkWidget *chk_eng = gtk_check_button_new_with_label("✔ Bảo vệ từ tiếng Anh & lập trình (Smart English Guard: giữ nguyên window, view, function, stream...)");
    gtk_toggle_button_set_active(GTK_TOGGLE_BUTTON(chk_eng), eng_on);
    g_signal_connect(chk_eng, "toggled", G_CALLBACK(on_eng_toggled), NULL);

    stats_lbl = gtk_label_new("📊 Phiên hiện tại: Sẵn sàng gõ Tiếng Việt mượt mà (Zero-Preedit)");
    gtk_style_context_add_class(gtk_widget_get_style_context(stats_lbl), "subtitle-lbl");
    gtk_widget_set_halign(stats_lbl, GTK_ALIGN_START);

    gtk_box_pack_start(GTK_BOX(card1), chk_split, FALSE, FALSE, 0);
    gtk_box_pack_start(GTK_BOX(card1), chk_eng, FALSE, FALSE, 0);
    gtk_box_pack_start(GTK_BOX(card1), stats_lbl, FALSE, FALSE, 0);
    gtk_box_pack_start(GTK_BOX(root_box), card1, FALSE, FALSE, 0);

    // Card 2: Interactive Typing Sandbox
    GtkWidget *card2 = gtk_box_new(GTK_ORIENTATION_VERTICAL, 8);
    gtk_style_context_add_class(gtk_widget_get_style_context(card2), "sandbox-card");

    GtkWidget *sb_top = gtk_box_new(GTK_ORIENTATION_HORIZONTAL, 8);
    sample_cat_lbl = gtk_label_new("🎯 Khu Vực Thử Phím Trực Tiếp");
    gtk_widget_set_halign(sample_cat_lbl, GTK_ALIGN_START);
    GtkWidget *next_btn = gtk_button_new_with_label("🎲 Đổi câu mẫu");
    g_signal_connect(next_btn, "clicked", G_CALLBACK(on_next_sample_btn), NULL);
    GtkWidget *clear_btn = gtk_button_new_with_label("🧹 Xóa ô thử");
    g_signal_connect(clear_btn, "clicked", G_CALLBACK(on_clear_test_btn), NULL);
    gtk_box_pack_start(GTK_BOX(sb_top), sample_cat_lbl, TRUE, TRUE, 0);
    gtk_box_pack_start(GTK_BOX(sb_top), clear_btn, FALSE, FALSE, 0);
    gtk_box_pack_start(GTK_BOX(sb_top), next_btn, FALSE, FALSE, 0);

    sample_quote_lbl = gtk_label_new("");
    gtk_style_context_add_class(gtk_widget_get_style_context(sample_quote_lbl), "quote-lbl");
    gtk_widget_set_halign(sample_quote_lbl, GTK_ALIGN_START);
    gtk_label_set_line_wrap(GTK_LABEL(sample_quote_lbl), TRUE);

    sample_hint_lbl = gtk_label_new("");
    gtk_style_context_add_class(gtk_widget_get_style_context(sample_hint_lbl), "hint-lbl");
    gtk_widget_set_halign(sample_hint_lbl, GTK_ALIGN_START);

    test_entry = gtk_entry_new();
    gtk_entry_set_placeholder_text(GTK_ENTRY(test_entry), "Click vào đây và gõ thử Tiếng Việt hoặc chuỗi nhacdidangnhactieng... rồi ấn Delete...");

    sample_link_btn = gtk_button_new_with_label("🎶 Trải nghiệm Phòng hát & Thu âm DSP tại VuiCaHat.com ➔");
    gtk_style_context_add_class(gtk_widget_get_style_context(sample_link_btn), "btn-eco");
    g_signal_connect(sample_link_btn, "clicked", G_CALLBACK(on_open_url_0), NULL);

    gtk_box_pack_start(GTK_BOX(card2), sb_top, FALSE, FALSE, 0);
    gtk_box_pack_start(GTK_BOX(card2), sample_quote_lbl, FALSE, FALSE, 0);
    gtk_box_pack_start(GTK_BOX(card2), sample_hint_lbl, FALSE, FALSE, 0);
    gtk_box_pack_start(GTK_BOX(card2), test_entry, FALSE, FALSE, 2);
    gtk_box_pack_start(GTK_BOX(card2), sample_link_btn, FALSE, FALSE, 0);
    gtk_box_pack_start(GTK_BOX(root_box), card2, FALSE, FALSE, 0);

    // Card 3: Official eChip Co., Ltd. & VuiCaHat.com Attribution Footer
    GtkWidget *card3 = gtk_box_new(GTK_ORIENTATION_VERTICAL, 6);
    gtk_style_context_add_class(gtk_widget_get_style_context(card3), "eco-card");
    GtkWidget *eco_lbl = gtk_label_new("💡 Sản phẩm mã nguồn mở phát triển bởi Công ty TNHH eChip (echipco.com) & dự án VuiCaHat.com:");
    gtk_style_context_add_class(gtk_widget_get_style_context(eco_lbl), "subtitle-lbl");
    gtk_widget_set_halign(eco_lbl, GTK_ALIGN_START);

    GtkWidget *eco_btns = gtk_box_new(GTK_ORIENTATION_HORIZONTAL, 10);
    GtkWidget *b_vch = gtk_button_new_with_label("🎤 VuiCaHat.com — Phòng Hát & Vocal DSP Studio");
    GtkWidget *b_echip = gtk_button_new_with_label("🏢 Công ty TNHH eChip — echipco.com");
    g_signal_connect(b_vch, "clicked", G_CALLBACK(on_open_url_1), NULL);
    g_signal_connect(b_echip, "clicked", G_CALLBACK(on_open_url_2), NULL);
    gtk_box_pack_start(GTK_BOX(eco_btns), b_vch, TRUE, TRUE, 0);
    gtk_box_pack_start(GTK_BOX(eco_btns), b_echip, TRUE, TRUE, 0);

    gtk_box_pack_start(GTK_BOX(card3), eco_lbl, FALSE, FALSE, 0);
    gtk_box_pack_start(GTK_BOX(card3), eco_btns, FALSE, FALSE, 0);
    gtk_box_pack_start(GTK_BOX(root_box), card3, FALSE, FALSE, 0);

    if (show_window_initially) {
        gtk_widget_show_all(win);
        gtk_widget_grab_focus(test_entry);
    }
}

static void update_sample_ui(const char *cat, const char *quote, const char *hint, const char *btn_text) {
    if (sample_cat_lbl) gtk_label_set_text(GTK_LABEL(sample_cat_lbl), cat);
    if (sample_quote_lbl) gtk_label_set_text(GTK_LABEL(sample_quote_lbl), quote);
    if (sample_hint_lbl) gtk_label_set_text(GTK_LABEL(sample_hint_lbl), hint);
    if (sample_link_btn) gtk_button_set_label(GTK_BUTTON(sample_link_btn), btn_text);
}

typedef struct {
    int enabled;
    int mode;
} StateUpdate;

static gboolean idle_update_state(gpointer data) {
    StateUpdate *st = (StateUpdate *)data;
    if (indicator) {
        app_indicator_set_icon_full(indicator, st->enabled ? "vuikey-vn" : "vuikey-en", "VuiKey");
        app_indicator_set_label(indicator, st->enabled ? "VN" : "EN", "VN");
    }
    if (status_btn) {
        gtk_button_set_label(GTK_BUTTON(status_btn),
            st->enabled ? "🟢 ĐANG BẬT TIẾNG VIỆT [VN] (Bấm hoặc Ctrl+Space để đổi)"
                        : "⚪ ĐANG Ở CHẾ ĐỘ TIẾNG ANH [EN] (Bấm để bật Tiếng Việt)");
    }
    free(st);
    return FALSE;
}

static void schedule_state_update(int enabled, int mode) {
    StateUpdate *st = (StateUpdate *)malloc(sizeof(StateUpdate));
    st->enabled = enabled;
    st->mode = mode;
    g_idle_add(idle_update_state, st);
}

static gboolean idle_update_stats(gpointer data) {
    char *txt = (char *)data;
    if (stats_lbl) {
        gtk_label_set_text(GTK_LABEL(stats_lbl), txt);
    }
    free(txt);
    return FALSE;
}

static void schedule_stats_update(const char *txt) {
    g_idle_add(idle_update_stats, strdup(txt));
}
*/
import "C"

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"unsafe"

	"github.com/mrtrinhvn/VuiKey/internal/config"
	"github.com/mrtrinhvn/VuiKey/internal/engine"
	"github.com/mrtrinhvn/VuiKey/internal/x11hook"
)

var (
	activeDaemon *x11hook.Daemon
	sampleIdx    int
)

func ensureTrayIcons() string {
	home, _ := os.UserHomeDir()
	iconDir := filepath.Join(home, ".config", "vuikey", "icons")
	_ = os.MkdirAll(iconDir, 0755)

	vnSVG := `<svg xmlns="http://www.w3.org/2000/svg" width="64" height="64" viewBox="0 0 64 64">
  <rect x="2" y="6" width="60" height="52" rx="14" fill="#10b981" stroke="#059669" stroke-width="3"/>
  <text x="32" y="42" text-anchor="middle" fill="#ffffff" font-family="Ubuntu, Inter, sans-serif" font-size="28" font-weight="900">VN</text>
</svg>`
	enSVG := `<svg xmlns="http://www.w3.org/2000/svg" width="64" height="64" viewBox="0 0 64 64">
  <rect x="2" y="6" width="60" height="52" rx="14" fill="#475569" stroke="#334155" stroke-width="3"/>
  <text x="32" y="42" text-anchor="middle" fill="#f8fafc" font-family="Ubuntu, Inter, sans-serif" font-size="28" font-weight="900">EN</text>
</svg>`
	_ = os.WriteFile(filepath.Join(iconDir, "vuikey-vn.svg"), []byte(vnSVG), 0644)
	_ = os.WriteFile(filepath.Join(iconDir, "vuikey-en.svg"), []byte(enSVG), 0644)
	return iconDir
}

//export goOnToggleClicked
func goOnToggleClicked() {
	if activeDaemon != nil {
		next := activeDaemon.ToggleEnabled()
		_ = config.Save(activeDaemon.Eng.Cfg)
		syncFcitx5State(next)
	}
}

//export goOnModeChanged
func goOnModeChanged(mode C.int) {
	if activeDaemon == nil {
		return
	}
	if mode == 99 {
		activeDaemon.SetEnabled(false)
		_ = config.Save(activeDaemon.Eng.Cfg)
		syncFcitx5State(false)
	} else {
		activeDaemon.SetEnabled(true)
		activeDaemon.SetMode(engine.InputMode(mode))
		_ = config.Save(activeDaemon.Eng.Cfg)
		syncFcitx5State(true)
	}
}

//export goOnFeatureToggled
func goOnFeatureToggled(feature C.int, active C.int) {
	if activeDaemon == nil {
		return
	}
	on := active != 0
	switch int(feature) {
	case 1:
		activeDaemon.Eng.Cfg.SyllableSplit = on
	case 2:
		activeDaemon.Eng.Cfg.SmartEnglish = on
	}
	_ = config.Save(activeDaemon.Eng.Cfg)
	syncFcitx5State(activeDaemon.Eng.Cfg.Enabled)
}

func syncFcitx5State(enabled bool) {
	if enabled {
		_ = exec.Command("fcitx5-remote", "-o").Run()
		_ = exec.Command("fcitx5-remote", "-s", "vuikey").Run()
	} else {
		_ = exec.Command("fcitx5-remote", "-s", "keyboard-us").Run()
		_ = exec.Command("fcitx5-remote", "-c").Run()
	}
}

//export goOnNextSampleClicked
func goOnNextSampleClicked() {
	sampleIdx = (sampleIdx + 1) % len(config.LyricSamples)
	applySample(sampleIdx)
}

//export goOnOpenURLClicked
func goOnOpenURLClicked(target C.int) {
	var url string
	switch int(target) {
	case 0:
		url = config.LyricSamples[sampleIdx%len(config.LyricSamples)].ActionURL
	case 1:
		url = "https://vuicahat.com"
	case 2:
		url = "https://echipco.com"
	}
	if url != "" {
		_ = exec.Command("xdg-open", url).Start()
	}
}

//export goOnQuitClicked
func goOnQuitClicked() {
	if activeDaemon != nil {
		_ = config.Save(activeDaemon.Eng.Cfg)
	}
}

func applySample(idx int) {
	s := config.LyricSamples[idx%len(config.LyricSamples)]
	cCat := C.CString(s.Category)
	cQuote := C.CString("“" + s.Quote + "”")
	cHint := C.CString(s.TelexHint)
	cBtn := C.CString(s.ActionText)
	defer C.free(unsafe.Pointer(cCat))
	defer C.free(unsafe.Pointer(cQuote))
	defer C.free(unsafe.Pointer(cHint))
	defer C.free(unsafe.Pointer(cBtn))
	C.update_sample_ui(cCat, cQuote, cHint, cBtn)
}

func RunApp(d *x11hook.Daemon, showWindow bool) {
	activeDaemon = d
	cfg := d.Eng.Cfg
	iconDir := ensureTrayIcons()
	cIconDir := C.CString(iconDir)
	defer C.free(unsafe.Pointer(cIconDir))

	d.OnStateChanged = func(enabled bool, mode engine.InputMode) {
		enInt := 0
		if enabled {
			enInt = 1
		}
		C.schedule_state_update(C.int(enInt), C.int(mode))
	}

	d.OnStatsUpdated = func(tk, vt int64) {
		msg := fmt.Sprintf("📊 Phiên hiện tại: %d phím đã gõ • %d dấu Tiếng Việt đã biến đổi mượt mà (Zero-Preedit)", tk, vt)
		cMsg := C.CString(msg)
		defer C.free(unsafe.Pointer(cMsg))
		C.schedule_stats_update(cMsg)
	}

	b2i := func(b bool) C.int {
		if b {
			return 1
		}
		return 0
	}

	C.init_ui(
		cIconDir,
		b2i(showWindow),
		b2i(cfg.Enabled),
		C.int(cfg.Mode),
		b2i(cfg.SyllableSplit),
		b2i(cfg.SmartEnglish),
	)
	applySample(0)
	C.gtk_main()
}
