package config

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/mrtrinhvn/VuiKey/internal/engine"
)

type LyricSample struct {
	Category   string
	Quote      string
	TelexHint  string
	ActionText string
	ActionURL  string
}

// Curated Typing Sandbox Samples — clean, respectful introduction to VuiCaHat.com and eChip Co., Ltd. (echipco.com)
var LyricSamples = []LyricSample{
	{
		Category:   "🎤 Khu Vực Thử Phím & Kiểm Tra Chống Nhân Đôi Chữ",
		Quote:      "Nhạc dịu dàng, tiếng hát bay cao trong không gian âm nhạc sống động",
		TelexHint:  "Thử gõ chuỗi liền không dấu cách: nhacdidangnhactiengdidangtiengchoitai rồi ấn phím Delete!",
		ActionText: "🎶 Trải nghiệm Phòng hát & Thu âm DSP tại VuiCaHat.com ➔",
		ActionURL:  "https://vuicahat.com",
	},
	{
		Category:   "🏢 Hệ Sinh Thái Sản Phẩm & Dịch Vụ Công Nghệ Số",
		Quote:      "Giải pháp công nghệ tối ưu cho trải nghiệm làm việc và giải trí hiện đại",
		TelexHint:  "Gõ thử: Giair phaps coong ngheej toois uwu cho trair nghieemj...",
		ActionText: "🌐 Khám phá Hệ sinh thái Công ty TNHH eChip (echipco.com) ➔",
		ActionURL:  "https://echipco.com",
	},
}

func ConfigPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", "vuikey", "config.json")
}

func Load() engine.Config {
	cfg := engine.DefaultConfig()
	data, err := os.ReadFile(ConfigPath())
	if err != nil {
		return cfg
	}
	_ = json.Unmarshal(data, &cfg)
	// Clean up any legacy macro entries
	delete(cfg.Macros, "vch")
	delete(cfg.Macros, "qh")
	return cfg
}

func Save(cfg engine.Config) error {
	p := ConfigPath()
	_ = os.MkdirAll(filepath.Dir(p), 0755)
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(p, data, 0644)
}
