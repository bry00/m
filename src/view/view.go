package view

import (
	"github.com/bry00/m/buffers"
	"github.com/bry00/m/config"
)

type Action int

const (
	ActionUnknown Action = iota
	ActionQuit
	ActionScrollUp
	ActionScrollDown
	ActionPageUp
	ActionPageDown
	ActionTop
	ActionBottom
	ActionScrollLeft
	ActionScrollFastLeft
	ActionScrollRight
	ActionScrollFastRight
	ActionHome
	ActionEnd
	ActionSearch
	ActionFindFirst
	ActionFindNext
	ActionFindPrevious
	ActionGotoLine
	ActionFlipNumbers
	ActionFlipRuler
	ActionMoveRulerUp
	ActionMoveRulerDown
	ActionReset
	ActionShortcuts
)

var actionNames = []string{
	"unknown",
	"quit",
	"scroll up",
	"scroll down",
	"page up",
	"page down",
	"top",
	"bottom",
	"scroll left",
	"scroll fast left",
	"scroll right",
	"scroll fast right",
	"home",
	"end",
	"search",
	"find first",
	"find next",
	"find previous",
	"go to line",
	"flip numbers",
	"flip ruler",
	"move ruler up",
	"move ruler down",
	"reset",
	"show shortcuts",
}

func (action Action) Count() int {
	return len(actionNames)
}

func (action Action) String() string {
	a := int(action)
	if a < 0 || a >= len(actionNames) {
		return actionNames[0]
	}
	return actionNames[a]
}

type AppStatus int

const (
	StatusUnknown AppStatus = iota
	StatusReady
	StatusReading
	StatusReceivingData
)

func (status AppStatus) String() string {
	names := []string{
		"unknown",
		"ready",
		"reading",
		"receiving",
	}
	a := int(status)
	if a < 0 || a >= len(names) {
		return names[0]
	}
	return names[a]
}

func (status AppStatus) Display() string {
	names := []string{
		"",
		"READY",
		"Reading...",
		"Receiving...",
	}
	a := int(status)
	if a < 0 || a >= len(names) {
		return names[0]
	}
	return names[a]
}

type TheViewController interface {
	DoAction(action Action)
	NoOfLines() int
	GetConfig() *config.Config
	GetFileNameTitle() string
	GetDataIterator(firstRow int) (*buffers.LineIndex, bool)
	DataReady() bool
	SetSearchText(text string, regex bool, ignoreCase bool)
	SetPointedLine(lineNo int)
}

type TheStatusBar interface {
	Reset()
	Message(format string, a ...interface{})
	Status(status AppStatus)
	SafeMessage(format string, a ...interface{})
	SafeStatus(status AppStatus)
}

type TheView interface {
	AreNumbersShown() bool
	DisplayAt(left int, top int)
	GetDisplayRect() (int, int, int, int)
	GetKeyShortcuts() map[Action][]string
	GetRulerPosition() int
	GetStatusBar() TheStatusBar
	IsRulerShown() bool
	Prepare()
	Refresh()
	SetController(ctl TheViewController)
	SetRulerPosition(index int)
	Show()
	ShowGotoLineDialog()
	ShowLine(lineIndex int)
	ShowNumbers(show bool)
	ShowRuler(show bool)
	ShowSearchDialog()
	ShowSearchResult(lineIndex int, start int, end int)
	ShowShortcuts()
	StopApplication()
}
