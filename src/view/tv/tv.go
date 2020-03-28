package tv

import (
	"fmt"
	"github.com/bry00/m/utl"
	"github.com/bry00/m/view"
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
	"log"
	"strconv"
	"strings"
)

type StatusBar struct {
	*tview.Box
	view   *View
    text   string
}

func newStatusBar(view *View) *StatusBar {
	r := &StatusBar {
		Box:  tview.NewBox(),
		view: view,

	}
	r.SetBorder(false)
	r.SetBorderPadding(0, 0, 0, 0)
	return r
}

func (sb *StatusBar) Reset() {
	sb.text = ""
}

func (sb *StatusBar) Message(format string, a ...interface{}) {
	sb.text = fmt.Sprintf(format, a)
}

func (sb *StatusBar) Draw(screen tcell.Screen) {
	var text string
	sb.Box.Draw(screen)
	ctl := sb.view.ctl.GetConfig()
	leftColumn, topRow, width, height := sb.view.GetDisplayRect()
	x, y, width, _ := sb.GetInnerRect()
	if len(sb.text) > 0 {
		text = sb.text
	} else {
		bottomRow := topRow + height
		topRow++
		leftColumn++
		totalRows := sb.view.ctl.NoOfLines()
		text = fmt.Sprintf("[%s]%d:%d - %d / %d", ctl.StatusBarTextColor, topRow, leftColumn, bottomRow, totalRows)
	}
	tview.Print(screen, text, x+1, y, width, tview.AlignLeft, tcell.ColorWhite)

}


type View struct {
	app   *tview.Application
	ctl    view.TheViewController
	text  *TextArea
	statusBar *StatusBar
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
	view.text.rulerIndex = index
}

type TextArea struct {
	*tview.Box
	view        *View
	firstLine   int
	firstColumn int
	width       int
	height      int
	rulerIndex  int
	showRuler   bool
}

func newTextArea(view *View) *TextArea {
	return &TextArea{
		Box:         tview.NewBox(),
		view:        view,
		firstLine:   0,
		firstColumn: 0,
		width:       0,
		height:      0,
		rulerIndex: -1,
		showRuler:   false,
	}
}

func (view *View)GetStatusBar() view.TheStatusBar {
	return view.statusBar
}

const rulerHeight = 3

func (t *TextArea) drawRuler(screen tcell.Screen, x int, y int, textWidth int) {
	var (
		line strings.Builder
	)
	config := t.view.ctl.GetConfig()
	color := tcell.GetColor(config.RulerColor)
	attr := fmt.Sprintf("[::%s]", config.RulerAttrs)
	line.Grow(textWidth + len(attr))


	for j := 0; j < 3; j++ {
		line.WriteString(attr)
		for c := 0; c < textWidth; c++ {
			n := (c + t.firstColumn + 1)
			digit := n % 10
			switch j {
			case 0:
				if n % 100 == 0 {
					line.WriteString( strconv.Itoa(n / 100 % 100))
				} else {
					//line.WriteString(" ")
					line.WriteRune(tcell.RuneVLine)
				}
			case 1:
				if digit == 0 {
					line.WriteString( strconv.Itoa(n / 10 % 10))
				} else {
					//line.WriteString(" ")
					line.WriteRune(tcell.RunePlus)
				}
			case 2:
				line.WriteString(strconv.Itoa(digit))
			}
		}
		tview.Print(screen, line.String(), x, y+j, textWidth, tview.AlignLeft, color)
		line.Reset()
	}
}

func (t *TextArea) getRulerPosition() int {
	if t.rulerIndex < 0 {
		t.rulerIndex = t.height / 2
	}
	max := utl.Min(t.height, t.view.ctl.NoOfLines())
	if t.rulerIndex > max {
		t.rulerIndex = max
	}
	return t.rulerIndex
}

func (t *TextArea) Draw(screen tcell.Screen) {
	if t.view.ctl != nil {
		conf := t.view.ctl.GetConfig()
		arrowLeft  := fmt.Sprintf("[%s::%s]%c", conf.SideArrowsColor, conf.SideArrowsArttrs, conf.SideArrowLeft)
		arrowRight := fmt.Sprintf("[%s::%s]%c", conf.SideArrowsColor, conf.SideArrowsArttrs, conf.SideArrowRight)
		tabSpaces := strings.Repeat(" ", conf.SpacesPerTab)
		t.Box.Draw(screen)
		xLeft, yTop, width, height := t.GetInnerRect()
		showRuler := t.showRuler && height > rulerHeight
		if showRuler {
			height = height - rulerHeight
		}
		t.height = height
		rulerIndex := t.getRulerPosition()

		textWidth := width - 2
		t.width = textWidth

		if iter, ok := t.view.ctl.GetDataIterator(t.firstLine); ok {
			var i int
			rulerDrawn := false
			for i = 0; i < height && iter.IndexOK(); i++ {
				y := yTop + i
				if t.showRuler {
					if i == rulerIndex {
						t.drawRuler(screen, xLeft+1, y, textWidth)
						rulerDrawn = true
					}
					if i >= rulerIndex {
						y += rulerHeight
					}
				}
				if line, err := iter.GetLine(); err != nil {
					log.Fatal(err)
				} else {
					r := []rune(tview.Escape(strings.Replace(line, "\t", tabSpaces, -1)))
					if t.firstColumn > len(r) {
						line = ""
					} else {
						line = string(r[t.firstColumn:])
					}
					if t.firstColumn > 0 {
						tview.PrintSimple(screen, arrowLeft, xLeft, y)

					}
					tview.Print(screen, line, xLeft+1, y, textWidth, tview.AlignLeft, tcell.ColorWhite)
					if tview.TaggedStringWidth(line) > textWidth {
						tview.PrintSimple(screen, arrowRight, xLeft+textWidth+1, y)
					}
				}
				iter.IndexIncrement()
			}
			if t.showRuler && !rulerDrawn {
				t.drawRuler(screen, xLeft+1, yTop + i, textWidth)
				t.rulerIndex = i
			}
		}
	}
}

func (t *TextArea) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
		return t.WrapInputHandler(func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
			//if t.view.ctl != nil {
				action := view.ActionUnknown
				switch event.Key() {
				case tcell.KeyUp:
					action = view.ActionScrollUp
				case tcell.KeyDown:
					action = view.ActionScrollDown
				case tcell.KeyPgDn:
					if event.Modifiers()&tcell.ModCtrl != 0 {
						action = view.ActionBottom
					} else {
						action = view.ActionPageDown
					}
				case tcell.KeyCtrlF:
					action = view.ActionPageDown
				case tcell.KeyCtrlSpace:
					action = view.ActionPageUp
				case tcell.KeyPgUp:
					if event.Modifiers()&tcell.ModCtrl != 0 {
						action = view.ActionTop
					} else {
						action = view.ActionPageUp
					}
				case tcell.KeyCtrlB:
					action = view.ActionPageUp
				case tcell.KeyHome:
					if event.Modifiers()&tcell.ModCtrl != 0 {
						action = view.ActionTop
					} else {
						action = view.ActionHome
					}
				case tcell.KeyEnd:
					if event.Modifiers()&tcell.ModCtrl != 0 {
						action = view.ActionBottom
					} else {
						action = view.ActionEnd
					}
				case tcell.KeyLeft:
					action = view.ActionScrollRight
				case tcell.KeyRight:
					action = view.ActionScrollLeft
				case tcell.KeyRune:
					switch event.Rune() {
					case ' ':
						action = view.ActionPageDown
					case 'q':
						action = view.ActionQuit
					case 'r':
						action = view.ActionFlipRuler
					case '-':
						action = view.ActionMoveRulerUp
					case '+':
						action = view.ActionMoveRulerDown
					default:
						return
					}
				case tcell.KeyEscape:
					action = view.ActionQuit
				default:
					return
				}
				t.view.ctl.DoAction(action)
			//}
		})
}


func NewView() *View {
	result := &View {
		app: nil,
		ctl: nil,
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
//func (v *View) Show() {
//	v.text.SetBorder(true).
//		SetBorderAttributes(tcell.AttrBold).
//		SetTitle(" " + v.ctl.GetFileNameTitle() + " ")
//	v.app = tview.NewApplication()
//
//	if err := v.app.SetRoot(v.text, true).Run(); err != nil {
//		panic(err)
//	}
//}

func (v *View) Show() {
	v.text.SetBorder(true).
		SetBorderAttributes(tcell.AttrBold).
		SetTitle(" " + v.ctl.GetFileNameTitle() + " ")

	v.app = tview.NewApplication()
	v.statusBar = newStatusBar(v)

	flex := tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(v.text, 0, 1, true).
			AddItem(v.statusBar, 1, 1, false)

	if err := v.app.SetRoot(flex, true).Run(); err != nil {
		panic(err)
	}
}


