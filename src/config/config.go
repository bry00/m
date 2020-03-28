package config

type Config struct {
	SpacesPerTab        int
	SideArrowLeft       rune
	SideArrowRight      rune
	SideArrowsColor     string
	SideArrowsArttrs    string
	StatusBarTextColor  string
	RulerColor          string
	RulerAttrs          string
}

func NewDefaultConfig() *Config {
	return &Config{
		SpacesPerTab:       4,
		SideArrowLeft:      '\u25C0',
		SideArrowRight:     '\u25B6',
		//SideArrowLeft:      '\u2B05',
		//SideArrowRight:     '\u2B95',
		//SideArrowLeft:      '\u2B45',
		//SideArrowRight:     '\u2B46',
		//SideArrowLeft:      '\u276E',
		//SideArrowRight:     '\u276F',
		SideArrowsColor:    "orange",
		SideArrowsArttrs:   "b",
		RulerColor:         "gold",
		RulerAttrs:         "rb",
		StatusBarTextColor: "orange",
	}

}

