package tv

import (
	"github.com/bry00/m/view"
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
	"log"
)

const rulerHeight = 3
const nummbersWidth = 8

const pageMain   = "main"
const pageSearch = "search"
const pageGoToLine = "goto-line"
const pageShortcuts = "shortcuts"

type View struct {
	app            *tview.Application
	ctl             view.TheViewController
	pages          *tview.Pages
	text           *TextArea
	statusBar      *StatusBar
	searchDialog   *SearchDialog
	lineDialog     *LineDialog
	shortcutWindow *ShortcutsWindow
}
func (view *View) ShowSearchResult(lineIndex int, start int, end int) {
	view.text.foundLine = lineIndex
	view.text.foundStart = start
	view.text.foundEnd = end
}

func (view *View) ShowLine(lineIndex int) {
	view.text.pointedLine = lineIndex
}

func (view *View) ShowNumbers(show bool) {
	view.text.showNumbers = show
}

func (view *View) AreNumbersShown() bool {
	return view.text.showNumbers
}

func (view *View) IsRulerShown() bool {
	return view.text.showRuler
}

func (view *View) ShowRuler(show bool) {
	view.text.showRuler = show
}

func (view *View) GetRulerPosition() int {
	if view.text.showRuler {
		return view.text.getRulerPosition()
	}
	return 0
}

func (view *View) SetRulerPosition(index int) {
	view.text.rulerPosition = index
}


func (view *View)GetStatusBar() view.TheStatusBar {
	return view.statusBar
}

func NewView() *View {
	result := &View {
		app: nil,
		ctl: nil,
		pages: tview.NewPages(),
		searchDialog: nil,
		lineDialog:   nil,
		statusBar:    nil,
		text:         nil,
	}
	result.text = newTextArea(result)
	return result
}

func (v *View)SetController(ctl view.TheViewController) {
	v.ctl = ctl
}

func (v *View) StopApplication() {
	v.app.Stop()
}


func (v *View)DisplayAt(left int, top int) {
	v.text.firstLine = top
	v.text.firstColumn = left
}

func (v *View)GetDisplayRect() (int, int, int, int) {
	return v.text.firstColumn, v.text.firstLine, v.text.width, v.text.height
}

func (v *View)Refresh() {
	v.app.Draw()
}


func (v *View)newModal(modal tview.Primitive, width int, height int) tview.Primitive {
	return tview.NewGrid().
		SetColumns(0, width, 0).
		SetRows(0, height, 0).
		AddItem(modal, 1, 1, 1, 1, 0, 0, true)
}

func (view *View) ShowSearchDialog() {
	if view.searchDialog != nil {
		view.searchDialog.Display()
	}
}

func (view *View) ShowGotoLineDialog() {
	if view.lineDialog != nil {
		view.lineDialog.Display()
	}
}

func (view *View)GenDefaultTheme() *tview.Theme {
	t := &view.ctl.GetConfig().Visual.Theme

	return &tview.Theme{
		PrimitiveBackgroundColor:    tcell.GetColor(t.PrimitiveBackgroundColor),
		ContrastBackgroundColor:     tcell.GetColor(t.ContrastBackgroundColor),
		MoreContrastBackgroundColor: tcell.GetColor(t.MoreContrastBackgroundColor),
		BorderColor:                 tcell.GetColor(t.BorderColor),
		TitleColor:                  tcell.GetColor(t.TitleColor),
		GraphicsColor:               tcell.GetColor(t.GraphicsColor),
		PrimaryTextColor:            tcell.GetColor(t.PrimaryTextColor),
		SecondaryTextColor:          tcell.GetColor(t.SecondaryTextColor),
		TertiaryTextColor:           tcell.GetColor(t.TertiaryTextColor),
		InverseTextColor:            tcell.GetColor(t.InverseTextColor),
		ContrastSecondaryTextColor:  tcell.GetColor(t.ContrastSecondaryTextColor),
	}
}

func (v *View) Prepare() {
	tview.Styles = *v.GenDefaultTheme()

	v.text.SetBorderColor(tview.Styles.BorderColor)
	v.text.SetTitleColor(tview.Styles.TitleColor)
	v.text.SetBackgroundColor(tview.Styles.PrimitiveBackgroundColor)

	v.text.SetBorder(true).
		//SetBorderAttributes(tcell.AttrBold).
		SetTitle(" " + v.ctl.GetFileNameTitle() + " ")
	v.statusBar = newStatusBar(v)

	v.app = tview.NewApplication()
	screen, err := tcell.NewScreen()
	if err != nil {
		log.Fatal(err.Error())
	}
	if err = screen.Init(); err != nil {
		log.Fatal(err.Error())
	}
	v.app.SetScreen(screen)
	screenWidth, screenHeight := screen.Size()

	pgMain := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(v.text, 0, 1, true).
		AddItem(v.statusBar, 1, 1, false)

	v.pages.AddPage(pageMain, pgMain, true, true).
		AddPage(pageSearch, v.newModal(newSearchDialog(v, screenWidth)), true, false).
		AddPage(pageGoToLine, v.newModal(newLineDialog(v)), true, false).
		AddPage(pageShortcuts, v.newModal(newShortcutsWindow(v.GetKeyShortcuts(), v, screenWidth, screenHeight)), true, false)

	v.app.EnableMouse(true)
}


func (v *View) Show() {
	if err := v.app.SetRoot(v.pages, true).Run(); err != nil {
		log.Fatal(err.Error())
	}
}

func (view *View) GetKeyShortcuts() map[view.Action][]string {
	return generateActionShortcutNames(textAreaShortcuts)
}

func (view *View) ShowShortcuts() {
	view.pages.ShowPage(pageShortcuts)
}
