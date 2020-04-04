package tv

import (
	"fmt"

	"github.com/bry00/m/view"
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

type StatusBar struct {
	*tview.Box
	view   *View
	status view.AppStatus
	text   *string
}

func newStatusBar(view *View) *StatusBar {
	r := &StatusBar{
		Box:  tview.NewBox(),
		view: view,
	}
	r.SetBorder(false)
	r.SetBorderPadding(0, 0, 0, 0)
	return r
}

func (sb *StatusBar) Reset() {
	sb.text = nil
}

func (sb *StatusBar) Message(format string, a ...interface{}) {
		text := fmt.Sprintf(format, a...)
		sb.text = &text
}

func (sb *StatusBar) SafeMessage(format string, a ...interface{}) {
	sb.view.app.QueueUpdateDraw(func() {
		sb.Message(format, a...)
	})
}

func (sb *StatusBar) Status(status view.AppStatus) {
		sb.status = status
}

func (sb *StatusBar) SafeStatus(status view.AppStatus) {
	sb.view.app.QueueUpdateDraw(func() {
		sb.Status(status)
	})
}

func (sb *StatusBar) Draw(screen tcell.Screen) {
	var text string
	sb.Box.Draw(screen)
	leftColumn, topRow, width, height := sb.view.GetDisplayRect()
	x, y, width, _ := sb.GetInnerRect()
	conf := sb.view.ctl.GetConfig()
	color := tcell.GetColor(conf.StatusBarTextColor)
	statusLabel := sb.status.Display()
	statusLabelWidth := tview.TaggedStringWidth(statusLabel)
	if statusLabelWidth > 0 {
		xLabel := width - (statusLabelWidth + 1)
		tview.Print(screen, statusLabel, xLabel, y, statusLabelWidth, tview.AlignLeft, color)
		width -= statusLabelWidth + 1
	}
	if sb.text != nil {
		text = *sb.text
	} else {
		bottomRow := topRow + height
		topRow++
		leftColumn++
		totalRows := sb.view.ctl.NoOfLines()
		text = fmt.Sprintf("[::%s]%d:%d - %d / %d", conf.StatusBarTextAttrs, topRow, leftColumn, bottomRow, totalRows)
	}
	tview.Print(screen, text, x+1, y, width, tview.AlignLeft, color)

}
