package tv

import (
	"fmt"
	"github.com/bry00/m/utl"
	"github.com/bry00/m/view"
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

const helpLabelAttribute = "b"
const helpLabelActions = "Actions"
const helpLabelShortcuts = "Keys"
const headerHeight = 3
const footerHeight = 1

const hPadding = 3

type shortcutDescription struct {
	name      string
	shortcuts []string
}

type ShortcutsWindow struct {
	*tview.Box
	view            *View
	descs           []shortcutDescription
	maxLabelLen     int
	maxDescLen      int
	pivot           int
	rows            int
	height          int
	foregroundColor tcell.Color
}

func newShortcutsWindow(shortcuts map[view.Action][]string, v *View, screeWidth int, screenHeight int) (win *ShortcutsWindow, width int, height int) {
	cnf := v.ctl.GetConfig()

	win = &ShortcutsWindow{
		Box:             tview.NewBox(),
		view:            v,
		maxLabelLen:     len(helpLabelActions),
		maxDescLen:      len(helpLabelShortcuts),
		pivot:           0,
		rows:            0,
		foregroundColor: tcell.GetColor(cnf.Visual.Help.ForegroundColor),
	}

	actions := view.ActionUnknown.Count()
	for k := view.ActionUnknown + 1; int(k) < actions; k++ {
		if v, exists := shortcuts[k]; exists {
			actionName := k.String()
			win.descs = append(win.descs, shortcutDescription{actionName, v})
			win.maxLabelLen = utl.MaxInt(win.maxLabelLen, len(actionName))
			for _, s := range v {
				win.maxDescLen = utl.MaxInt(win.maxDescLen, len(s))
				win.rows++
			}
			win.rows++
		}
	}

	width = utl.MinInt(1+hPadding+win.maxLabelLen+hPadding+win.maxDescLen+1+hPadding, screeWidth)
	height = utl.MinInt(headerHeight+win.rows+footerHeight, screenHeight-5)
	v.shortcutWindow = win
	win.SetBorder(true)

	win.SetBorderColor(tcell.GetColor(cnf.Visual.Help.BorderColor))
	win.SetBorderAttributes(tcell.AttrNone)
	win.SetBackgroundColor(tcell.GetColor(cnf.Visual.Help.BackgroundColor))

	return
}

func (w *ShortcutsWindow) Draw(screen tcell.Screen) {
	w.Box.Draw(screen)

	def := tcell.StyleDefault
	background := def.Background(w.GetBackgroundColor())
	border := background.Foreground(w.GetBorderColor()) | tcell.Style(w.GetBorderAttributes())

	left, top, width, height := w.GetInnerRect()

	screen.SetContent(left-1, top+1, tview.BoxDrawingsVerticalDoubleAndRightSingle, nil, border)
	for i := 0; i < width; i++ {
		screen.SetContent(left+i, top+1, tview.BoxDrawingsLightHorizontal, nil, border)
	}
	screen.SetContent(left+width, top+1, tview.BoxDrawingsVerticalDoubleAndLeftSingle, nil, border)

	xAction := left + hPadding
	xDesc := left + width - hPadding - w.maxDescLen

	const helpLabelFormat = "[::%s]%s"
	tview.Print(screen, fmt.Sprintf(helpLabelFormat, helpLabelAttribute, helpLabelActions), xAction, top, w.maxLabelLen, tview.AlignLeft, w.foregroundColor)
	tview.Print(screen, fmt.Sprintf(helpLabelFormat, helpLabelAttribute, helpLabelShortcuts), xDesc, top, w.maxDescLen, tview.AlignLeft, w.foregroundColor)

	w.height = height - headerHeight - footerHeight
	row := 0
	for n := 0; n < len(w.descs) && row < w.height; n++ {
		d := w.getDesc(n + w.pivot)
		for j := 0; j < len(d.shortcuts) && row < w.height; j++ {
			if j == 0 {
				tview.Print(screen, d.name, xAction, top+headerHeight+row, w.maxLabelLen, tview.AlignLeft, w.foregroundColor)
			}
			tview.Print(screen, d.shortcuts[j], xDesc, top+headerHeight+row, w.maxLabelLen, tview.AlignLeft, w.foregroundColor)
			row++
		}
		row++
	}
}

func (w *ShortcutsWindow) getDesc(index int) shortcutDescription {
	length := len(w.descs)
	if index < 0 {
		index = length - ((-index-1)%length + 1)
	} else {
		index %= length
	}
	return w.descs[index]
}

func (w *ShortcutsWindow) nextPivot(start int, delta int) int {
	row := 0
	result := 0
	for n := 0; row < w.height; n++ {
		d := w.getDesc(n*delta + start)
		row += len(d.shortcuts) + 1
		result = n

	}
	return start + result*delta
}

func (w *ShortcutsWindow) Display() {
	w.view.pages.ShowPage(pageShortcuts)
}

func (w *ShortcutsWindow) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return w.WrapInputHandler(func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
		switch event.Key() {
		case tcell.KeyUp:
			w.pivot--
		case tcell.KeyEnter:
			fallthrough
		case tcell.KeyDown:
			w.pivot++
		case tcell.KeyHome:
			w.pivot = 0
		case tcell.KeyEnd:
			w.pivot = w.nextPivot(len(w.descs)-1, -1)
		case tcell.KeyCtrlSpace:
			fallthrough
		case tcell.KeyCtrlB:
			fallthrough
		case tcell.KeyF7:
			w.pivot = w.nextPivot(w.pivot, -1)
		case tcell.KeyPgUp:
			if event.Modifiers()&tcell.ModCtrl != 0 {
				w.pivot = 0
			} else {
				w.pivot = w.nextPivot(w.pivot, -1)
			}
		case tcell.KeyCtrlF:
			fallthrough
		case tcell.KeyF8:
			w.pivot = w.nextPivot(w.pivot, 1)
		case tcell.KeyPgDn:
			if event.Modifiers()&tcell.ModCtrl != 0 {
				w.pivot = w.nextPivot(len(w.descs)-1, -1)
			} else {
				w.pivot = w.nextPivot(w.pivot, 1)
			}
		case tcell.KeyRune:
			switch event.Rune() {
			case ' ':
				w.pivot = w.nextPivot(w.pivot, 1)
			}
		case tcell.KeyEscape:
			w.view.pages.SwitchToPage(pageMain)
			return
		}
	})
}
