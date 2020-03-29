package controller

import (
	"bufio"
	"github.com/bry00/m/buffers"
	"github.com/bry00/m/config"
	"github.com/bry00/m/view"
	"log"
	"os"
	"path/filepath"
	"strings"
	"unicode/utf8"
)

type Controller struct {
	fileName       string
	conf          *config.Config
	maxLineLength  int
	view           view.TheView
	data          *buffers.BufferedData
}

func NewController(fileName string, data *buffers.BufferedData, view view.TheView, conf *config.Config) *Controller {
	var (
		filePath string
		err      error
	)
	if len(fileName) > 0 {
		filePath, err = filepath.Abs(fileName)
		if err != nil {
			log.Fatal(err)
		}
		if !fileExists(filePath) {
			log.Fatalf("File \"%s\" does not exist!\n", filePath)
		}
	}
	result := &Controller {
		fileName:      filePath,
		conf:          conf,
		maxLineLength: 0,
		data:          data,
		view:          view,
	}
	view.SetController(result)
	return result
}

func (ctl *Controller) GetConfig() *config.Config {
	return ctl.conf
}

func (ctl *Controller) Run() {
	go ctl.readFile()
	ctl.view.Show()
}

func (ctl *Controller) OnExit() {
	if ctl.data != nil {
		ctl.data.Close()
		ctl.data = nil
	}
}

func (ctl *Controller) NoOfLines() int {
	return ctl.data.Len()
}

func (ctl *Controller) GetFileNameTitle() string {
	if len(ctl.fileName) == 0 {
		return "<<stdin>>"
	} else {
		if home, err := os.UserHomeDir(); err == nil && strings.HasPrefix(ctl.fileName, home) {
			return "~" + strings.TrimPrefix(ctl.fileName, home)
		}
		return ctl.fileName
	}

	return ctl.fileName
}


func (ctl *Controller) GetDataIterator(firstRow int) (*buffers.LineIndex, bool) {
	result := ctl.data.NewLineIndexer()
	if result.IndexSet(firstRow, false) {
		return result, true
	} else {
		return nil, false
	}
}

func (ctl *Controller) DoAction(action view.Action) {
	lines := ctl.data.Len()
	left, top, width, height := ctl.view.GetDisplayRect()
	switch action {
	case view.ActionScrollUp:
		top -= 1
	case view.ActionScrollDown:
		top += 1
	case view.ActionTop:
		top = 0
	case view.ActionBottom:
		top = lines - height
	case view.ActionHome:
		left = 0
	case view.ActionEnd:
		left = ctl.maxLineLength - width + 1
	case view.ActionPageUp:
		top -= height
	case view.ActionPageDown:
		top += height
	case view.ActionScrollLeft:
		left += 1
	case view.ActionScrollRight:
		left -= 1
	case view.ActionFlipRuler:
		ctl.view.ShowRuler(!ctl.view.IsRulerShown())
	case view.ActionMoveRulerUp:
		if ctl.view.IsRulerShown() {
			ruler := ctl.view.GetRulerPosition() - 1
			if ruler < 0 {
				ruler = 0
			}
			ctl.view.SetRulerPosition(ruler)
		}
	case view.ActionMoveRulerDown:
		if ctl.view.IsRulerShown() {
			ruler := ctl.view.GetRulerPosition() + 1
			if ruler > height {
				ruler = height
			}
			ctl.view.SetRulerPosition(ruler)
		}
	case view.ActionFlipNumbers:
		ctl.view.ShowNumbers(!ctl.view.AreNumbersShown())
	case view.ActionQuit:
		ctl.view.StopApplication()
		return
	default:
		return
	}
	if top >= lines - height {
		top = lines - height
	}
	if top < 0 {
		top = 0
	}
	if left > ctl.maxLineLength - width + 1 {
		left = ctl.maxLineLength - width + 1
	}
	if left < 0 {
		left = 0
	}
	ctl.view.DisplayAt(left, top)
}


func (ctl *Controller)readFile() {
	var (
		file *os.File
		err  error
	)
	_, _, _, height := ctl.view.GetDisplayRect()

	if len(ctl.fileName) > 0 {
		if file, err = os.Open(ctl.fileName); err != nil {
			log.Fatal(err)
		} else {
			defer file.Close()
		}
	} else {
		file = os.Stdin
	}
	ctl.data = buffers.NewBufferedDataDefault()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		currentLength := lengthExpandedTabs(line, ctl.conf.SpacesPerTab)
		ctl.data.AddLine(line)
		if currentLength > ctl.maxLineLength {
			ctl.maxLineLength = currentLength
		}
		if ctl.data.Len() < height {
			ctl.view.Refresh()
		}
	}
	ctl.view.Refresh()
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func lengthExpandedTabs(line string, tabSpaces int) int {
	total := utf8.RuneCountInString(line)
	tabs := strings.Count(line, "\t")
	return (total - tabs) + tabs * tabSpaces
}