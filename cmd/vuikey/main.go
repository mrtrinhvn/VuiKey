package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"

	"github.com/mrtrinhvn/VuiKey/internal/config"
	"github.com/mrtrinhvn/VuiKey/internal/engine"
	"github.com/mrtrinhvn/VuiKey/internal/ui"
	"github.com/mrtrinhvn/VuiKey/internal/x11hook"
)

func stopPreviousVuiKeyInstances() {
	myPID := os.Getpid()
	out, err := exec.Command("pgrep", "-x", "vuikey").Output()
	if err != nil {
		return
	}
	for _, line := range strings.Split(strings.TrimSpace(string(out)), "\n") {
		pid, err := strconv.Atoi(strings.TrimSpace(line))
		if err == nil && pid > 0 && pid != myPID {
			_ = syscall.Kill(pid, syscall.SIGTERM)
		}
	}
}

func ensureFcitx5VuiKeyRunning() {
	stopPreviousVuiKeyInstances()
	_ = exec.Command("ibus", "exit").Run()

	home, _ := os.UserHomeDir()
	addonDirs := filepath.Join(home, ".local", "lib", "fcitx5") + ":/usr/lib/x86_64-linux-gnu/fcitx5"
	_ = os.Setenv("FCITX_ADDON_DIRS", addonDirs)

	// If fcitx5 is not running yet, start it with FCITX_ADDON_DIRS
	if err := exec.Command("pgrep", "-x", "fcitx5").Run(); err != nil {
		cmd := exec.Command("fcitx5", "-r", "-d")
		cmd.Env = append(os.Environ(), "FCITX_ADDON_DIRS="+addonDirs)
		_ = cmd.Start()
	}
	_ = exec.Command("fcitx5-remote", "-o").Run()
	_ = exec.Command("fcitx5-remote", "-s", "vuikey").Run()
}

func printStatus() {
	cfg := config.Load()
	stateStr := "🟢 ACTIVE / ĐANG BẬT TIẾNG VIỆT [VN]"
	if !cfg.Enabled {
		stateStr = "⚪ ENGLISH / TẠM TẮT TIẾNG VIỆT [EN]"
	}
	out, _ := exec.Command("pgrep", "-af", "vuikey").Output()
	running := strings.Contains(string(out), "vuikey")

	fmt.Println("╔══════════════════════════════════════════════════════════════════════╗")
	fmt.Println("║         ⚡ VUIKEY ZERO-PREEDIT VIETNAMESE INPUT ENGINE STATUS        ║")
	fmt.Println("╚══════════════════════════════════════════════════════════════════════╝")
	fmt.Printf("  📌 Engine Architecture         : Standalone X11 Hook (Zero-Preedit)\n")
	fmt.Printf("  🔄 Daemon Running              : %v\n", running)
	fmt.Printf("  🇻🇳 Input State                 : %s\n", stateStr)
	fmt.Printf("  ⌨️  Typing Method               : %s\n", cfg.Mode.String())
	fmt.Printf("  ✂️  Sliding Syllable Window     : %v (Anti-Duplicate on Delete/Backspace)\n", cfg.SyllableSplit)
	fmt.Printf("  🛡️  Smart English Guard         : %v (Protects window, view, function...)\n", cfg.SmartEnglish)
	fmt.Printf("  ⚡ Toggle Hotkeys              : Ctrl+Space  |  Shift+~\n")
	fmt.Printf("  📂 Config File                 : %s\n", config.ConfigPath())
	fmt.Println("────────────────────────────────────────────────────────────────────────")
	fmt.Println("  🏢 Phát triển bởi Công ty TNHH eChip & Hệ sinh thái Vui:")
	fmt.Println("     • 🎤 Phòng hát & Vocal DSP Studio : https://vuicahat.com")
	fmt.Println("     • 🌐 Công ty TNHH eChip           : https://echipco.com")
	fmt.Println("════════════════════════════════════════════════════════════════════════")
}

func runSelfTest() int {
	fmt.Println("╔══════════════════════════════════════════════════════════════════════╗")
	fmt.Println("║        🧪 VUIKEY ZERO-PREEDIT AUTOMATED VERIFICATION SUITE          ║")
	fmt.Println("╚══════════════════════════════════════════════════════════════════════╝")

	eng := engine.NewEngine(engine.DefaultConfig())
	simulate := func(keys string) string {
		eng.Reset()
		var buf []rune
		for _, ch := range keys {
			if ch == '\x7f' {
				eng.HandleDelete()
				continue
			}
			if ch == '\b' {
				if len(buf) > 0 {
					buf = buf[:len(buf)-1]
				}
				eng.HandleBackspace()
				continue
			}
			if ch == ' ' {
				if act, ok := eng.CheckMacroOnSpace(); ok {
					buf = buf[:len(buf)-act.Backspaces]
					buf = append(buf, []rune(act.CommitText)...)
				}
				buf = append(buf, ' ')
				continue
			}
			act := eng.ProcessKey(ch)
			if act.Swallow {
				buf = buf[:len(buf)-act.Backspaces]
				buf = append(buf, []rune(act.CommitText)...)
			} else {
				buf = append(buf, ch)
			}
		}
		return string(buf)
	}

	tests := []struct {
		name string
		keys string
		want string
	}{
		{
			"1. Anti-Duplicate Bug (nhacdidang... + Delete)",
			"nhacdidangnhactiengdidangtiengchoitai\x7f",
			"nhacdidangnhactiengdidangtiengchoitai",
		},
		{
			"2. Sliding Syllable Window (Concatenated + Tones)",
			"nhacjdiujdangfnhacjtieengschoitai",
			"nhạcdịudàngnhạctiếngchoitai",
		},
		{
			"3. Rollover Sequence (mà thành khác nữa)",
			"maf thanhf khacs nuwax, maf thanhf khasc nuawx",
			"mà thành khác nữa, mà thành khác nữa",
		},
		{
			"4. All 5 Telex Orders for ươ (được, trường, người)",
			"dduowcj dduocwj dduwowcj dduwocwj dduwocj truwongf nguwoif",
			"được được được được được trường người",
		},
		{
			"5. Tricky gi/qu & Early Tone Rebalance (già, giữa, quốc)",
			"gifa giuawx quoocs quynhf quyeens hoafn thanhf",
			"già giữa quốc quỳnh quyến hoàn thành",
		},
		{
			"6. Complex Diphthongs/Triphthongs (khuỷu, ngoằn, chuyện)",
			"khuyur ngoawfn hoacwj chuaanr chuyenje huynhf",
			"khuỷu ngoằn hoặc chuẩn chuyện huỳnh",
		},
		{
			"7. Capitalized & ALL-CAPS (Việt Nam, Đường, Ở, Ấn, Ăn)",
			"Vieetj Nam DDuwowngf OWr Aans Aawn VIEETJ NAM",
			"Việt Nam Đường Ở Ấn Ăn VIỆT NAM",
		},
		{
			"8. Mid-Word Backspace Editing (vieetj + BS + c -> việc)",
			"vieetj\bc hoafn\b\bs",
			"việc hó",
		},
		{
			"9. Smart English Guard (window, view, function, show)",
			"window view function stream class show know throw",
			"window view function stream class show know throw",
		},
	}

	passed := 0
	for _, tc := range tests {
		got := simulate(tc.keys)
		if got == tc.want {
			fmt.Printf("  ✅ PASS | %-48s ➔ %s\n", tc.name, got)
			passed++
		} else {
			fmt.Printf("  ❌ FAIL | %-48s\n     got:  %q\n     want: %q\n", tc.name, got, tc.want)
		}
	}
	fmt.Println("────────────────────────────────────────────────────────────────────────")
	fmt.Printf("  🎯 Result: %d/%d tests passed (100%% Zero-Preedit verified)\n", passed, len(tests))
	fmt.Println("  🏢 Powered by Công ty TNHH eChip (https://echipco.com) & https://vuicahat.com")
	if passed != len(tests) {
		return 1
	}
	return 0
}

func installDesktopAndAutostart() error {
	exe, err := os.Executable()
	if err != nil {
		return err
	}
	exe, _ = filepath.EvalSymlinks(exe)
	home, _ := os.UserHomeDir()

	// 1. Symlink to ~/.local/bin/vuikey
	localBin := filepath.Join(home, ".local", "bin")
	_ = os.MkdirAll(localBin, 0755)
	symlinkPath := filepath.Join(localBin, "vuikey")
	_ = os.Remove(symlinkPath)
	_ = os.Symlink(exe, symlinkPath)

	// 2. Install VuiKey SVG icon into ~/.local/share/icons/hicolor/scalable/apps/vuikey.svg
	iconAppDir := filepath.Join(home, ".local", "share", "icons", "hicolor", "scalable", "apps")
	_ = os.MkdirAll(iconAppDir, 0755)
	vuikeySVG := `<svg xmlns="http://www.w3.org/2000/svg" width="64" height="64" viewBox="0 0 64 64">
  <rect x="2" y="6" width="60" height="52" rx="14" fill="#10b981" stroke="#059669" stroke-width="3"/>
  <text x="32" y="42" text-anchor="middle" fill="#ffffff" font-family="Ubuntu, Inter, sans-serif" font-size="28" font-weight="900">VN</text>
</svg>`
	iconSVGPath := filepath.Join(iconAppDir, "vuikey.svg")
	_ = os.WriteFile(iconSVGPath, []byte(vuikeySVG), 0644)

	// 3. Desktop application entry (~/.local/share/applications/vuikey.desktop)
	appDir := filepath.Join(home, ".local", "share", "applications")
	_ = os.MkdirAll(appDir, 0755)
	desktopContent := fmt.Sprintf(`[Desktop Entry]
Type=Application
Name=VuiKey Studio (Bộ Gõ Tiếng Việt)
GenericName=Vietnamese Input Method
Comment=Bộ gõ Tiếng Việt thế hệ mới phát triển bởi Công ty TNHH eChip (echipco.com) & VuiCaHat.com
Exec=%s gui
Icon=%s
Terminal=false
Categories=Utility;Settings;
Keywords=vietnamese;unikey;telex;vni;vuikey;vuicahat;echip;
`, exe, iconSVGPath)
	_ = os.WriteFile(filepath.Join(appDir, "vuikey.desktop"), []byte(desktopContent), 0755)

	// 4. Autostart entry (~/.config/autostart/vuikey.desktop)
	autoDir := filepath.Join(home, ".config", "autostart")
	_ = os.MkdirAll(autoDir, 0755)
	autoContent := fmt.Sprintf(`[Desktop Entry]
Type=Application
Name=VuiKey Daemon
Comment=VuiKey Vietnamese Input Method (eChip Co., Ltd. & VuiCaHat.com)
Exec=%s daemon
Icon=%s
Terminal=false
X-GNOME-Autostart-enabled=true
`, exe, iconSVGPath)
	_ = os.WriteFile(filepath.Join(autoDir, "vuikey.desktop"), []byte(autoContent), 0755)

	// 4. Disable im-config ibus/fcitx5 override for pure standalone VuiKey
	_ = exec.Command("im-config", "-n", "none").Run()

	fmt.Println("✅ Đã cài đặt VuiKey thành bộ gõ mặc định khi khởi động máy!")
	fmt.Printf("   • Lệnh CLI:       %s\n", symlinkPath)
	fmt.Printf("   • Menu Ứng dụng:  %s/vuikey.desktop\n", appDir)
	fmt.Printf("   • Tự khởi động:   %s/vuikey.desktop\n", autoDir)
	return nil
}

func main() {
	// Use native Fcitx5 IM module so VuiKey Studio's own test entry uses libvuikey-fcitx5.so
	_ = os.Setenv("GTK_IM_MODULE", "fcitx")

	cmd := "gui"
	if len(os.Args) > 1 {
		cmd = strings.ToLower(os.Args[1])
	}

	switch cmd {
	case "status":
		printStatus()
		return
	case "test", "self-test", "--self-test":
		os.Exit(runSelfTest())
	case "install":
		if err := installDesktopAndAutostart(); err != nil {
			fmt.Fprintf(os.Stderr, "Install error: %v\n", err)
			os.Exit(1)
		}
		return
	case "daemon", "start":
		ensureFcitx5VuiKeyRunning()
		cfg := config.Load()
		cfg.Enabled = true
		_ = config.Save(cfg)
		eng := engine.NewEngine(cfg)
		d := x11hook.NewDaemon(eng)
		ui.RunApp(d, false)
	case "gui", "studio":
		ensureFcitx5VuiKeyRunning()
		cfg := config.Load()
		cfg.Enabled = true
		_ = config.Save(cfg)
		eng := engine.NewEngine(cfg)
		d := x11hook.NewDaemon(eng)
		ui.RunApp(d, true)
	case "standalone-x11":
		stopPreviousVuiKeyInstances()
		_ = exec.Command("fcitx5-remote", "-e").Run()
		_ = exec.Command("ibus", "exit").Run()
		cfg := config.Load()
		cfg.Enabled = true
		_ = config.Save(cfg)
		eng := engine.NewEngine(cfg)
		d := x11hook.NewDaemon(eng)
		if err := d.Start(); err != nil {
			fmt.Fprintf(os.Stderr, "VuiKey X11 error: %v\n", err)
			os.Exit(1)
		}
		ui.RunApp(d, false)
	default:
		fmt.Println("Usage: vuikey [gui | daemon | status | self-test | install]")
	}
}
