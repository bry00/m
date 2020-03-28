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
	}
	a := int(action)
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
}

type TheStatusBar interface {
	Reset()
	Message(format string, a ...interface{})
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
}
