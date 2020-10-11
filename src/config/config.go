package config

import (
	"github.com/bry00/m/buffers"
	"github.com/bry00/m/utl"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
)

const ConfigFile = "config.yaml"

type CnfDataBuffer struct {
	BlockSizeLimitMB int `yaml:"blockSizeLimitMB"`
	TotalSizeLimitMB int `yaml:"totalSizeLimitMB"`
}

type CnfView struct {
	SpacesPerTab       int `yaml:"spacesPerTab"`
	ViewRefreshSeconds int `yaml:"viewRefreshSeconds"`
}

type CnfSideArrows struct {
	Left  int    `yaml:"left"`
	Right int    `yaml:"right"`
	Color string `yaml:"color"`
	Attrs string `yaml:"attrs"`
}

type CnfStatusBar struct {
	TextColor string `yaml:"textColor"`
	TextAttrs string `yaml:"textAttrs"`
}

type CnfRuler struct {
	Color string `yaml:"color"`
	Attrs string `yaml:"attrs"`
}

type CnfNumbers struct {
	Color string `yaml:"color"`
}

type CnfHelp struct {
	BackgroundColor string `yaml:"backgroundColor"`
	ForegroundColor string `yaml:"foregroundColor"`
	BorderColor     string `yaml:"borderColor"`
}

type CnfSearch struct {
	IgnoreCase bool `yaml:"ignoreCase"`
}

type CnfTheme struct {
	PrimitiveBackgroundColor    string
	ContrastBackgroundColor     string
	MoreContrastBackgroundColor string
	BorderColor                 string
	TitleColor                  string
	GraphicsColor               string
	PrimaryTextColor            string
	SecondaryTextColor          string
	TertiaryTextColor           string
	InverseTextColor            string
	ContrastSecondaryTextColor  string
}

type CnfVisual struct {
	SideArrows CnfSideArrows `yaml:"sideArrows"`
	StatusBar  CnfStatusBar  `yaml:"statusBar"`
	Ruler      CnfRuler      `yaml:"ruler"`
	Numbers    CnfNumbers    `yaml:"numbers"`
	Help       CnfHelp       `yaml:"help"`
	Theme      CnfTheme      `yaml:"theme"`
}

type Config struct {
	DataBuffer CnfDataBuffer `yaml:"dataBuffer"`
	Search     CnfSearch     `yaml:"search"`
	View       CnfView       `yaml:"view"`
	Visual     CnfVisual     `yaml:"visual"`
}

func NewDefaultConfig() *Config {
	return &Config{
		DataBuffer: CnfDataBuffer{
			BlockSizeLimitMB: buffers.DefaultBlockSizeLimit / buffers.MB,
			TotalSizeLimitMB: buffers.DefaultTotalSizeLimit / buffers.MB,
		},
		Search: CnfSearch{
			IgnoreCase: true,
		},
		View: CnfView{
			SpacesPerTab:       4,
			ViewRefreshSeconds: 5,
		},
		Visual: CnfVisual{
			SideArrows: CnfSideArrows{
				Left:  '\u25C0',
				Right: '\u25B6',
				// Left:      '\u2B05',
				// Right:     '\u2B95',
				// Left:      '\u2B45',
				// Right:     '\u2B46',
				// Left:      '\u276E',
				// Right:     '\u276F',
				Color: "orange",
				Attrs: "b",
			},
			Ruler: CnfRuler{
				Color: "gold",
				Attrs: "rb",
			},
			StatusBar: CnfStatusBar{
				TextColor: "gold",
				TextAttrs: "",
			},
			Numbers: CnfNumbers{
				Color: "gold",
			},
			Help: CnfHelp{
				BackgroundColor: "beige",
				ForegroundColor: "darkGreen",
				BorderColor:     "darkGreen",
			},
			Theme: CnfTheme{
				PrimitiveBackgroundColor:    "#001000",
				ContrastBackgroundColor:     "maroon",
				MoreContrastBackgroundColor: "goldenRod",
				BorderColor:                 "moccasin",
				TitleColor:                  "gold",
				GraphicsColor:               "white",
				PrimaryTextColor:            "lightYellow",
				SecondaryTextColor:          "lemonChiffon",
				TertiaryTextColor:           "khaki",
				InverseTextColor:            "seaGreen",
				ContrastSecondaryTextColor:  "yellow",
			},
		},
	}
}

func GetConfig(prog string) *Config {
	result := NewDefaultConfig() // Default config
	configFile := GetConfigFileName(prog)
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		if cnf, err := yaml.Marshal(result); err != nil {
			log.Printf("Cannot encode default configuration: %v", err)
		} else {
			if err := ioutil.WriteFile(configFile, cnf, 0600); err != nil {
				log.Printf("Cannot save default configuration: %v", err)
			}
		}
	} else {
		if cnf, err := ioutil.ReadFile(configFile); err != nil {
			log.Printf("Error reading configuration file: %v", err)
		} else {
			if err := yaml.Unmarshal(cnf, result); err != nil {
				log.Fatalf("Error decoding configuration file: %v", err)
			}
		}
	}
	utl.OnEachFieldWithSuffix(result, "Color",
		func(s string) string { return strings.ToLower(strings.TrimSpace(s)) })
	return result
}

func GetConfigDir(prog string) string {
	configDir, err := os.UserConfigDir()
	if err != nil {
		log.Fatal(err)
	}
	result := path.Join(configDir, prog)
	if _, err := os.Stat(result); os.IsNotExist(err) {
		os.MkdirAll(result, 0700)
	}
	return result
}

func GetConfigFileName(prog string) string {
	return path.Join(GetConfigDir(prog), ConfigFile)
}
