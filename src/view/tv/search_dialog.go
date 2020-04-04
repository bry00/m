package tv

import (
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
	"strings"
)

type SearchDialog struct {
	*tview.Form
	view        *View
	startOfEdit  bool
}

func newSearchDialog(view *View, screenWidth int) (dialog *SearchDialog, width int, height int) {
	width = screenWidth / 3 * 2
	if width < 20 {
		width = 20
	}
	form := tview.NewForm().
		AddInputField("Find:", "", width - 10, nil, nil).
    	AddCheckbox("Ignore Case:", false, nil).
		AddCheckbox("Plain:", false, nil)

	dialog = &SearchDialog{
		Form:  form,
		view:  view,
	}
	cancelFun := func() {
		view.pages.SwitchToPage(pageMain)
	}

	okFun := func() {
		searchText := strings.TrimSpace(dialog.GetSearchText())
		if len(searchText) > 0 {
			//view.statusBar.Message("SEARCH: %s", searchText)
			view.ctl.SetSearchText(searchText, dialog.IsRegexSearch(), dialog.IsIgnoreCaseSearch())
			view.app.QueueEvent(tcell.NewEventKey(tcell.KeyRune, 'n', 0))
		}
		view.pages.SwitchToPage(pageMain)
	}
	form.SetButtonsAlign(tview.AlignRight).
		AddButton("Search", okFun).
		AddButton("Cancel", cancelFun).
		SetCancelFunc(cancelFun)
	form.SetBorder(true)

	searchField := dialog.GetSearchField()
	searchField.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEnter:
			okFun()
			return nil
		case tcell.KeyHome:
			return event
		case tcell.KeyRune:
			if dialog.startOfEdit {
				searchField.SetText("")
				dialog.startOfEdit = false
			}
		default:
			dialog.startOfEdit = false
		}
		return event
	})

	height = 11
	view.searchDialog = dialog
	return
}

func (s *SearchDialog)Display() {
	s.view.pages.ShowPage(pageSearch)
	s.view.app.SetFocus(s.GetSearchField())
	s.startOfEdit = true
	s.view.app.QueueEvent(tcell.NewEventKey(tcell.KeyHome, 0, 0))
}

func (s *SearchDialog)GetSearchField() *tview.InputField {
	return s.GetFormItem(0).(*tview.InputField)
}

func (s *SearchDialog)GetIgnoreCaseCheck() *tview.Checkbox {
	return s.GetFormItem(1).(*tview.Checkbox)
}

func (s *SearchDialog)GetPlainCheck() *tview.Checkbox {
	return s.GetFormItem(2).(*tview.Checkbox)
}

func (s *SearchDialog)GetSearchText() string {
	return s.GetSearchField().GetText()
}

func (s *SearchDialog)IsIgnoreCaseSearch() bool {
	return s.GetIgnoreCaseCheck().IsChecked()
}

func (s *SearchDialog)IsRegexSearch() bool {
	return !s.GetPlainCheck().IsChecked()
}




