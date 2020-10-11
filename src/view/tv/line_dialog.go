package tv

import (
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
	"strconv"
	"unicode"
)

type LineDialog struct {
	*tview.Form
	view        *View
	startOfEdit bool
}

func newLineDialog(view *View) (dialog *LineDialog, width int, height int) {
	width = 23
	height = 7

	form := tview.NewForm().
		AddInputField("Line: ", "", width-10, func(textToCheck string, lastChar rune) bool {
			if unicode.IsDigit(lastChar) {
				return true
			}
			return false
		}, nil)

	dialog = &LineDialog{
		Form: form,
		view: view,
	}
	cancelFun := func() {
		view.GetStatusBar().Reset()
		view.pages.SwitchToPage(pageMain)
	}

	okFun := func() {
		lineNo := dialog.GetLineNo()
		lines := view.ctl.NoOfLines()
		if lineNo > lines {
			view.GetStatusBar().Message("There are only %d lines in this file, thus you cannot go to line %d",
				lines, lineNo)
		} else {
			if lineNo > 0 {
				view.ctl.SetPointedLine(lineNo)
				view.app.QueueEvent(tcell.NewEventKey(tcell.KeyRune, ':', 0))
			} else {
				view.ctl.SetPointedLine(-1) // reset
			}
			view.pages.SwitchToPage(pageMain)
		}
	}
	form.
		AddButton("Search", func() {
			dialog.view.app.SetFocus(dialog.GetLineField())
			view.app.QueueEvent(tcell.NewEventKey(tcell.KeyEnter, '\x00', 0))
		}).
		AddButton("Cancel", func() {
			dialog.view.app.SetFocus(dialog.GetLineField())
			view.app.QueueEvent(tcell.NewEventKey(tcell.KeyEscape, '\x00', 0))
		}).
		SetCancelFunc(cancelFun)
	form.SetBorder(true)

	lineField := dialog.GetLineField()
	lineField.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		view.GetStatusBar().Reset()
		switch event.Key() {
		case tcell.KeyEnter:
			okFun()
			return nil
		case tcell.KeyHome:
			return event
		case tcell.KeyRune:
			if dialog.startOfEdit {
				lineField.SetText("")
				dialog.startOfEdit = false
			}
		default:
			dialog.startOfEdit = false
		}
		return event
	})

	view.lineDialog = dialog
	return
}

func (s *LineDialog) Display() {
	s.view.pages.ShowPage(pageGoToLine)
	s.view.app.SetFocus(s.GetLineField())
	s.startOfEdit = true
	s.view.app.QueueEvent(tcell.NewEventKey(tcell.KeyHome, 0, 0))
}

func (s *LineDialog) GetLineField() *tview.InputField {
	return s.GetFormItem(0).(*tview.InputField)
}

func (s *LineDialog) GetLineNo() int {
	if result, err := strconv.Atoi(s.GetLineField().GetText()); err != nil {
		return 0
	} else {
		return result
	}
}
