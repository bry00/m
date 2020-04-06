package tv

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/bry00/m/utl"
	"github.com/bry00/m/view"
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

type TextArea struct {
	*tview.Box
	view          *View
	firstLine     int
	firstColumn   int
	width         int
	height        int
	rulerPosition int
	foundLine     int
	foundStart    int
	foundEnd      int
	pointedLine   int
	showRuler     bool
	showNumbers   bool
}

func newTextArea(view *View) *TextArea {
	return &TextArea{
		Box:           tview.NewBox(),
		view:          view,
		firstLine:     0,
		firstColumn:   0,
		width:         0,
		height:        0,
		rulerPosition: -1,
		pointedLine:   -1,
		foundLine:     -1,
		foundStart:    -1,
		foundEnd:      -1,
		showRuler:     false,
		showNumbers:   false,
	}
}

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
				if n%100 == 0 {
					line.WriteString(strconv.Itoa(n / 100 % 100))
				} else {
					//line.WriteString(" ")
					line.WriteRune(tcell.RuneVLine)
				}
			case 1:
				if digit == 0 {
					line.WriteString(strconv.Itoa(n / 10 % 10))
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
	if t.rulerPosition < 0 {
		t.rulerPosition = t.height / 2
	}
	lines := t.view.ctl.NoOfLines()
	if t.view.ctl.DataReady() {
		lines = utl.Min(t.height, lines)
	} else {
		lines = t.height
	}
	if t.rulerPosition > lines {
		t.rulerPosition = lines
	}
	return t.rulerPosition
}

const aReverse = "[::r]"
const aNormal  = "[::-]"

func numberString(n int, width int) string {
	num :=strconv.Itoa(n)
	l := len(num)
	if l > width {
		num = num[l-width:]
		l = width
	}
	leading := width - len(num)
	var result strings.Builder
	if leading > 0 {
		result.WriteString("[::d]")
		result.WriteString(strings.Repeat("0", leading))
	}
	result.WriteString("[::b]")
	result.WriteString(num)
	return result.String()
}

func (t *TextArea) Draw(screen tcell.Screen) {
	if t.view.ctl != nil {
		conf := t.view.ctl.GetConfig()
		numbersColor := tcell.GetColor(conf.NumbersColor)
		arrowLeft := fmt.Sprintf("[%s::%s]%c", conf.SideArrowsColor, conf.SideArrowsArttrs, conf.SideArrowLeft)
		arrowRight := fmt.Sprintf("[%s::%s]%c", conf.SideArrowsColor, conf.SideArrowsArttrs, conf.SideArrowRight)
		tabSpaces := strings.Repeat(" ", conf.SpacesPerTab)
		t.Box.Draw(screen)
		xBase, yTop, width, height := t.GetInnerRect()
		showRuler := t.showRuler && height > rulerHeight
		if showRuler {
			height = height - rulerHeight
		}
		t.height = height
		rulerIndex := t.getRulerPosition()

		textWidth := width - 2

		var xLeft int
		if t.showNumbers {
			textWidth = textWidth - nummbersWidth
			xLeft = xBase + nummbersWidth
		} else {
			xLeft = xBase
		}

		if textWidth < 0 {
			textWidth = 0
		}

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
					lineIndex := t.firstLine + i
					if t.showNumbers {
						//tview.Print(screen, fmt.Sprintf("[::b]%*d",
						//	nummbersWidth, lineIndex+1),
						//	xBase, y, nummbersWidth, tview.AlignLeft, numbersColor)
						tview.Print(screen, numberString(lineIndex+1, nummbersWidth),
							xBase, y, nummbersWidth, tview.AlignLeft, numbersColor)
					}
					r := []rune(tview.Escape(strings.Replace(line, "\t", tabSpaces, -1)))
					if t.firstColumn > len(r) {
						line = ""
					} else {
						if lineIndex == t.foundLine && t.foundStart >= 0 && t.foundEnd > t.firstColumn {
							var str strings.Builder
							if t.firstColumn < t.foundStart {
								str.WriteString(string(r[t.firstColumn: t.foundStart]))
							}
							if lineIndex == t.pointedLine {
								str.WriteString(aNormal)
							} else {
								str.WriteString(aReverse)
							}
							str.WriteString(string(r[t.foundStart:t.foundEnd]))
							if lineIndex == t.pointedLine {
								str.WriteString(aReverse)
							} else {
								str.WriteString(aNormal)
							}
							str.WriteString(string(r[t.foundEnd:]))
							line = str.String()
						} else {
							line = string(r[t.firstColumn:])
						}
					}

					lineLen := tview.TaggedStringWidth(line)

					if lineIndex == t.pointedLine {
						line = fmt.Sprintf("[::r]%-*s", textWidth + (len(line) - lineLen) + 5, line)
					}

					if t.firstColumn > 0 {
						tview.PrintSimple(screen, arrowLeft, xLeft, y)

					}
					tview.Print(screen, line, xLeft+1, y, textWidth, tview.AlignLeft, tcell.ColorWhite)
					if lineLen > textWidth {
						tview.PrintSimple(screen, arrowRight, xLeft+textWidth+1, y)
					}
				}
				iter.IndexIncrement()
			}
			if t.showRuler && !rulerDrawn {
				t.drawRuler(screen, xLeft+1, yTop+i, textWidth)
				t.rulerPosition = i
			}
		}
	}
}


//func (t *TextArea) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
//	return t.WrapInputHandler(func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
//		//if t.view.ctl != nil {
//		action := view.ActionUnknown
//		t.view.statusBar.Reset()
//		switch event.Key() {
//		case tcell.KeyUp:
//			action = view.ActionScrollUp
//		case tcell.KeyEnter:
//			fallthrough
//		case tcell.KeyDown:
//			action = view.ActionScrollDown
//		case tcell.KeyPgDn:
//			if event.Modifiers()&tcell.ModCtrl != 0 {
//				action = view.ActionBottom
//			} else {
//				action = view.ActionPageDown
//			}
//		case tcell.KeyCtrlF:
//			action = view.ActionPageDown
//		case tcell.KeyCtrlSpace:
//			action = view.ActionPageUp
//		case tcell.KeyPgUp:
//			if event.Modifiers()&tcell.ModCtrl != 0 {
//				action = view.ActionTop
//			} else {
//				action = view.ActionPageUp
//			}
//		case tcell.KeyCtrlB:
//			action = view.ActionPageUp
//		case tcell.KeyCtrlN:
//			action = view.ActionFlipNumbers
//		case tcell.KeyCtrlL:
//			fallthrough
//		case tcell.KeyCtrlG:
//			action = view.ActionGotoLine
//		case tcell.KeyHome:
//			if event.Modifiers()&tcell.ModCtrl != 0 {
//				action = view.ActionTop
//			} else {
//				action = view.ActionHome
//			}
//		case tcell.KeyEnd:
//			if event.Modifiers()&tcell.ModCtrl != 0 {
//				action = view.ActionBottom
//			} else {
//				action = view.ActionEnd
//			}
//		case tcell.KeyLeft:
//			action = view.ActionScrollRight
//		case tcell.KeyRight:
//			action = view.ActionScrollLeft
//		case tcell.KeyRune:
//			switch event.Rune() {
//			case '/':
//				action = view.ActionSearch
//			case ' ':
//				action = view.ActionPageDown
//			case 'q':
//				action = view.ActionQuit
//			case 'n':
//				action = view.ActionFindNext
//			case 'N':
//				action = view.ActionFindPrevious
//			case '-':
//				action = view.ActionMoveRulerUp
//			case '+':
//				action = view.ActionMoveRulerDown
//			case 'r':
//				action = view.ActionFlipRuler
//			case ':':
//				action = view.ActionGotoLine
//			case '\\':
//				action = view.ActionReset
//			case 'g':
//				action = view.ActionTop
//			case 'G':
//				action = view.ActionBottom
//			default:
//				return
//			}
//		case tcell.KeyEscape:
//			action = view.ActionQuit
//		default:
//			return
//		}
//		t.view.ctl.DoAction(action)
//		//}
//	})
//}

func (t *TextArea) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return t.WrapInputHandler(func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
		t.view.statusBar.Reset()
		action := textAreaShortcutMap.mapKeys(event)
		if action != view.ActionUnknown {
			t.view.ctl.DoAction(action)
		}
	})
}

var (
	textAreaShortcuts = []shortcut {
		{key: tcell.KeyUp, action: view.ActionScrollUp},
		{key: tcell.KeyDown, action: view.ActionScrollDown},
		{key: tcell.KeyEnter, action: view.ActionScrollDown},
		{key: tcell.KeyPgDn, action: view.ActionPageDown},
		{key: tcell.KeyCtrlF, action: view.ActionPageDown},
		{key: tcell.KeyPgDn, mod: tcell.ModCtrl, action: view.ActionBottom},
		{key: tcell.KeyCtrlSpace, action: view.ActionPageUp},
		{key: tcell.KeyPgUp, action: view.ActionPageUp},
		{key: tcell.KeyPgUp, mod: tcell.ModCtrl, action: view.ActionTop},
		{key: tcell.KeyCtrlB, action: view.ActionPageUp},
		{key: tcell.KeyCtrlN, action: view.ActionFlipNumbers},
		{key: tcell.KeyCtrlL, action: view.ActionGotoLine},
		{key: tcell.KeyCtrlG, action: view.ActionGotoLine},
		{key: tcell.KeyHome, action: view.ActionHome},
		{key: tcell.KeyHome, mod: tcell.ModCtrl, action: view.ActionTop},
		{key: tcell.KeyEnd, action: view.ActionEnd},
		{key: tcell.KeyEnd, mod: tcell.ModCtrl, action: view.ActionBottom},
		{key: tcell.KeyLeft, action: view.ActionScrollRight},
		{key: tcell.KeyRight, action: view.ActionScrollLeft},

		{r: ' ', action: view.ActionPageDown},
		{r: '/', action: view.ActionSearch},

		{r: 'n', action: view.ActionFindNext},
		{r: 'N', action: view.ActionFindPrevious},
		{r: 'r', action: view.ActionFlipRuler},
		{r: '-', action: view.ActionMoveRulerUp},
		{r: '+', action: view.ActionMoveRulerDown},
		{r: ':', action: view.ActionGotoLine},
		{r: '\\', action: view.ActionReset},
		{r: 'g', action: view.ActionTop},
		{r: 'G', action: view.ActionBottom},


		{r: 'q', action: view.ActionQuit},
		{key: tcell.KeyEscape, action: view.ActionQuit},
	}
	textAreaShortcutMap *shortcutMap
)

func init() {
	textAreaShortcutMap = newShortcutMap(textAreaShortcuts)
}