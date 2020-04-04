package view

import (
	"github.com/bry00/m/buffers"
	"github.com/bry00/m/config"
)

type Action int
const (
	ActionUnknown Action = iota
	ActionShowHelp
	ActionQuit
	ActionScrollUp
	ActionScrollDown
	ActionScrollLeft
	ActionScrollRight
	ActionPageUp
	ActionPageDown
	ActionTop
	ActionBottom
	ActionHome
	ActionEnd
	ActionFlipRuler
	ActionMoveRulerUp
	ActionMoveRulerDown
	ActionFlipNumbers
	ActionSearch
	ActionFindNext
	ActionFindPrevious
	ActionGotoLine
	ActionReset
)

func (action Action)String() string {
	names := []string{
		"unknown",
		"showHelp",
		"quit",
		"scrollUp",
		"scrollDown",
		"scrollLeft",
		"scrollRight",
		"pageUp",
		"pageDown",
		"top",
		"bottom",
		"home",
		"end",
		"flipRuler",
		"moveRulerUp",
		"moveRulerDown",
		"flipNumbers",
		"search",
		"findNext",
		"findPrevious",
		"gotoLine",
		"reset",
	}
	a := int(action)
	if a < 0 || a >= len(names) {
		return names[0]
	}
	return names[a]
}


type AppStatus int
const (
	StatusUnknown AppStatus = iota
	StatusReady
	StatusReading
	StatusReceivingData
)

func (status AppStatus)String() string {
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

func (status AppStatus)Display() string {
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
	Show();
	StopApplication();
	DisplayAt(left int, top int);
	GetDisplayRect() (int, int, int, int)
	SetController(ctl TheViewController)
	Refresh()
	GetStatusBar() TheStatusBar
	ShowRuler(show bool)
	IsRulerShown() bool
	GetRulerPosition() int
	SetRulerPosition(index int)
	ShowNumbers(show bool)
	AreNumbersShown() bool
	ShowSearchDialog()
	ShowGotoLineDialog()
	ShowLine(lineIndex int)
	ShowSearchResult(lineIndex int, start int, end int)
}
