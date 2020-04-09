package config

import "strings"

type Config struct {
	SpacesPerTab        int
	SideArrowLeft       rune
	SideArrowRight      rune
	SideArrowsColor     string
	SideArrowsArttrs    string
	StatusBarTextColor  string
	StatusBarTextAttrs  string
	RulerColor          string
	RulerAttrs          string
	NumbersColor        string
	HelpBackgroundColor string
	HelpForegroundColor string
	HelpBorderColor     string
	ViewRefreshSeconds  int
}

func NewDefaultConfig() *Config {
	return &Config{
		SpacesPerTab:       4,
		ViewRefreshSeconds: 5,
		SideArrowLeft:      '\u25C0',
		SideArrowRight:     '\u25B6',
		//SideArrowLeft:      '\u2B05',
		//SideArrowRight:     '\u2B95',
		//SideArrowLeft:      '\u2B45',
		//SideArrowRight:     '\u2B46',
		//SideArrowLeft:      '\u276E',
		//SideArrowRight:     '\u276F',
		SideArrowsColor:    n("orange"),
		SideArrowsArttrs:   "b",
		RulerColor:         n("gold"),
		RulerAttrs:         "rb",
		StatusBarTextColor: n("gold"),
		StatusBarTextAttrs: "",
		NumbersColor:       n("gold"),
		HelpBackgroundColor: n("AliceBlue"),
		HelpForegroundColor: n("MidnightBlue"),
		HelpBorderColor:     n("MidnightBlue"),
	}

}

func n(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}
