package controller

import (
	"bufio"
	"github.com/bry00/m/buffers"
	"github.com/bry00/m/config"
	"github.com/bry00/m/view"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"
)

type Controller struct {
	fileName          *string
	title             *string
	conf              *config.Config
	maxLineLength      int
	view               view.TheView
	data              *buffers.BufferedData
	dataReady          bool
	searchString       string
	searchRegex        bool
	searchIgnoreCase   bool
	searchLastRow      int
	searchLastCol      int
	pointedLine        int
}


func NewController(fileName string, data *buffers.BufferedData, view view.TheView, conf *config.Config) *Controller {
	var (
		filePath *string = nil
	)
	if len(fileName) > 0 {
		if absPath, err := filepath.Abs(fileName);  err != nil {
			log.Fatal(err)
		} else {
			if !fileExists(absPath) {
				log.Fatalf("File \"%s\" does not exist!\n", absPath)
			}
			filePath = &absPath
		}
	}
	result := &Controller {
		fileName:       filePath,
		title:          nil,
		conf:           conf,
		maxLineLength:  0,
		data:           data,
		view:           view,
		dataReady:      false,
		searchString:     "",
		searchRegex:      false,
		searchIgnoreCase: false,
		searchLastRow:   -1,
		searchLastCol:   -1,
		pointedLine:   -1,
	}
	view.SetController(result)
	return result
}

func (ctl *Controller) DataReady() bool {
	return ctl.dataReady
}

func (ctl *Controller) GetConfig() *config.Config {
	return ctl.conf
}

func (ctl *Controller) Run() {
	ctl.view.Prepare()
	go ctl.readFile()
	ctl.view.Show()
}

func (ctl *Controller) OnExit() {
	if ctl.data != nil {
		ctl.data.Close()
		ctl.data = nil
	}
}

func (ctl *Controller) NoOfLines() int {
	return ctl.data.Len()
}

func (ctl *Controller) GetFileNameTitle() string {
	if ctl.title == nil {
		if ctl.fileName == nil {
			return "<<stdin>>"
		}
		if home, err := os.UserHomeDir(); err == nil && strings.HasPrefix(*ctl.fileName, home) {
			return "~" + strings.TrimPrefix(*ctl.fileName, home)
		}
		return *ctl.fileName
	}
	return *ctl.title
}


func (ctl *Controller) GetDataIterator(firstRow int) (*buffers.LineIndex, bool) {
	result := ctl.data.NewLineIndexer()
	if result.IndexSet(firstRow, false) {
		return result, true
	} else {
		return nil, false
	}
}

func setFoundStringPosition(left int, top int, width int, height int, foundLine int, foundStart int, foundEnd int) (int, int) {
	if foundLine >= 0 && foundStart >= 0 && foundEnd >= 0 {

		if foundLine < top || foundLine >= top + height {
			top = foundLine - height / 3
		}

		if foundEnd >= left + width {
			left = foundEnd - width
		}
		if foundStart < left {
			left = foundStart
		}

		if left < 0 {
			left = 0
		}
		if top < 0 {
			top = 0
		}
	}
	return left, top
}


func (ctl *Controller) DoAction(action view.Action) {
	lines := ctl.data.Len()
	left, top, width, height := ctl.view.GetDisplayRect()

	switch action {
	case view.ActionReset:
		ctl.pointedLine = -1
		ctl.view.ShowLine(ctl.pointedLine)
		ctl.view.ShowRuler(false)
		ctl.view.ShowNumbers(false)
		ctl.view.ShowSearchResult(-1, -1, -1)
	case view.ActionScrollUp:
		top -= 1
	case view.ActionScrollDown:
		top += 1
	case view.ActionTop:
		top = 0
	case view.ActionBottom:
		top = lines - height
	case view.ActionHome:
		left = 0
	case view.ActionEnd:
		left = ctl.maxLineLength - width + 1
	case view.ActionPageUp:
		top -= height
	case view.ActionPageDown:
		top += height
	case view.ActionScrollLeft:
		left += 1
	case view.ActionScrollRight:
		left -= 1
	case view.ActionFlipRuler:
		ctl.view.ShowRuler(!ctl.view.IsRulerShown())
	case view.ActionMoveRulerUp:
		if ctl.view.IsRulerShown() {
			ruler := ctl.view.GetRulerPosition() - 1
			if ruler < 0 {
				ruler = 0
			}
			ctl.view.SetRulerPosition(ruler)
		}
	case view.ActionMoveRulerDown:
		if ctl.view.IsRulerShown() {
			ruler := ctl.view.GetRulerPosition() + 1
			if ruler > height {
				ruler = height
			}
			ctl.view.SetRulerPosition(ruler)
		}
	case view.ActionFlipNumbers:
		ctl.view.ShowNumbers(!ctl.view.AreNumbersShown())
	case view.ActionSearch:
		ctl.searchLastRow = -1
		ctl.searchLastCol = -1
		ctl.view.ShowSearchResult(-1, -1, -1)
		ctl.view.ShowSearchDialog()
	case view.ActionFindNext:
		if ctl.searchLastRow < 0 {
			ctl.searchLastRow = top
		}
		if ctl.searchLastCol < 0 {
			ctl.searchLastCol = left
		}
		foundLine, foundStart, foundEnd, err := ctl.findNext(ctl.searchLastRow, ctl.searchLastCol)
		if err == nil {
			ctl.searchLastRow = foundLine
			ctl.searchLastCol = foundEnd
			ctl.view.ShowSearchResult(foundLine,foundStart, foundEnd)
			if ctl.searchLastRow >= 0 {
				left, top = setFoundStringPosition(left, top, width, height, foundLine, foundStart, foundEnd)
				ctl.view.GetStatusBar().Message("Found at: %d:%d \"%s\"",
					foundLine + 1, foundStart + 1, ctl.searchString)
			} else {
				ctl.view.GetStatusBar().Message("Cannot find: \"%s\"", ctl.searchString)
				ctl.searchLastRow = 0
				ctl.searchLastCol = 0
			}
		} else {
			ctl.view.GetStatusBar().Message("Wrong search string \"%s\": %s", ctl.searchString, err.Error())
		}
	case view.ActionFindPrevious:
		if ctl.searchLastRow < 0 {
			ctl.searchLastRow = ctl.NoOfLines()
			ctl.searchLastCol = -1
		} else {
			ctl.searchLastCol--
		}
		foundLine, foundStart, foundEnd, err := ctl.findPrevious(ctl.searchLastRow, ctl.searchLastCol)
		if err == nil {
			ctl.searchLastRow = foundLine
			ctl.searchLastCol = foundEnd
			ctl.view.ShowSearchResult(foundLine,foundStart, foundEnd)
			if ctl.searchLastRow >= 0 {
				left, top = setFoundStringPosition(left, top, width, height, foundLine, foundStart, foundEnd)
				ctl.view.GetStatusBar().Message("Previous at: %d:%d \"%s\"",
					foundLine + 1, foundStart + 1, ctl.searchString)
			} else {
				ctl.view.GetStatusBar().Message("Cannot find previous: \"%s\"", ctl.searchString)
				ctl.searchLastRow = ctl.NoOfLines() - 1
				ctl.searchLastCol = -1
			}
		} else {
			ctl.view.GetStatusBar().Message("Wrong search string \"%s\": %s", ctl.searchString, err.Error())
		}
	case view.ActionGotoLine:
		if ctl.pointedLine > 0 {
			if ctl.pointedLine > lines {
				ctl.view.GetStatusBar().Message("Wrong line number: %d", ctl.pointedLine)
			} else {
				lineIndex := ctl.pointedLine - 1
				top = lineIndex - height / 3
				ctl.view.ShowLine(lineIndex)
				ctl.view.GetStatusBar().Message("Line #%d", ctl.pointedLine)
			}
			ctl.pointedLine = -1

		} else {
			ctl.view.ShowGotoLineDialog()
		}
	case view.ActionQuit:
		ctl.view.StopApplication()
		return
	case view.ActionShortcuts:
		ctl.view.ShowShortcuts()
		return
	default:
		return
	}
	if top >= lines - height {
		top = lines - height
	}
	if top < 0 {
		top = 0
	}
	if left > ctl.maxLineLength - width + 1 {
			left = ctl.maxLineLength - width + 1
	}
	if left < 0 {
		left = 0
	}
	ctl.view.DisplayAt(left, top)
}

func (ctl *Controller) SetSearchText(text string, regex bool, ignoreCase bool) {
	ctl.searchString = text
	ctl.searchRegex = regex
	ctl.searchIgnoreCase = ignoreCase
}

func (ctl *Controller) findPrevious(startLine int, startColumn int) (int, int, int, error) {
	var (
		err          error
		re           *regexp.Regexp
	)
	tabSpaces := strings.Repeat(" ", ctl.conf.SpacesPerTab)
	searchString := ctl.searchString

	if ctl.searchRegex {
		if ctl.searchIgnoreCase && !strings.HasPrefix(ctl.searchString, "(?i)") {
			searchString = "(?i)" + searchString
		}
		re, err = regexp.Compile(searchString)
		if err != nil {
			return -1, -1, -1, err
		}
	} else {
		if ctl.searchIgnoreCase {
			searchString = strings.ToUpper(searchString)
		}
	}
	search := func(txt string) ([]int) {
		if ctl.searchRegex {
			if f := re.FindAllStringIndex(txt, -1); f != nil {
                return f[len(f)-1]
			}
			return nil;
		}
		if ctl.searchIgnoreCase {
			txt = strings.ToUpper(txt)
		}
		i := strings.LastIndex(txt, searchString)
		if i<0 {
			return nil
		}
		result := make([]int, 2)
		result[0] = i
		result[1] = i + len(searchString)
		return result
	}
	line := -1
	start := -1
	end := -1
	limit := startColumn
	i := ctl.data.NewLineIndexer()
	i.IndexSet(startLine, false)
	for ; i.IndexOK() ; i.IndexDecrement()  {
		if txt, err := i.GetLine(); err == nil {
			txt = strings.Replace(txt, "\t", tabSpaces, -1)
			if limit > 0 {
				txt = txt[0:limit]
				limit = 0
			}
			if found := search(txt); found != nil {
				line = i.Index()
				start = found[0]
				end = found[1]
				break
			}
		}
	}
	return line, start, end, nil
}

func (ctl *Controller) findNext(startLine int, startColumn int) (int, int, int, error) {
	var (
		err          error
		re           *regexp.Regexp
	)
	tabSpaces := strings.Repeat(" ", ctl.conf.SpacesPerTab)
	searchString := ctl.searchString

	if ctl.searchRegex {
		if ctl.searchIgnoreCase && !strings.HasPrefix(ctl.searchString, "(?i)") {
			searchString = "(?i)" + searchString
		}
		re, err = regexp.Compile(searchString)
		if err != nil {
			return -1, -1, -1, err
		}
	} else {
		if ctl.searchIgnoreCase {
			searchString = strings.ToUpper(searchString)
		}
	}
	search := func(txt string) ([]int) {
		if ctl.searchRegex {
			return re.FindStringIndex(txt)
		}
		if ctl.searchIgnoreCase {
			txt = strings.ToUpper(txt)
		}
		i := strings.Index(txt, searchString)
		if i<0 {
			return nil
		}
		result := make([]int, 2)
		result[0] = i
		result[1] = i + len(searchString)
		return result
	}
	line := -1
	start := -1
	end := -1
	offset := startColumn
	i := ctl.data.NewLineIndexer()
	i.IndexSet(startLine, false)
	for ; i.IndexOK() ; i.IndexIncrement()  {
		if txt, err := i.GetLine(); err == nil {
			txt = strings.Replace(txt, "\t", tabSpaces, -1)
			if offset < len(txt) {
				if offset > 0 {
					txt = txt[offset:]
				}
				if found := search(txt); found != nil {
					line = i.Index()
					start = found[0] + offset
					end = found[1] + offset
					break
				}
			}
		}
		offset = 0
	}
	return line, start, end, nil
}


func (ctl *Controller) SetPointedLine(lineNo int) {
	ctl.pointedLine = lineNo
}


func (ctl *Controller)readFile() {
	var (
		file *os.File
		err  error
	)

	if ctl.fileName != nil {
		ctl.view.GetStatusBar().SafeStatus(view.StatusReading)
		if file, err = os.Open(*ctl.fileName); err != nil {
			log.Fatal(err)
		} else {
			defer file.Close()
		}
	} else {
		ctl.view.GetStatusBar().SafeStatus(view.StatusReceivingData)
		file = os.Stdin
	}
	ctl.data = buffers.NewBufferedDataDefault()

	_, _, _, height := ctl.view.GetDisplayRect()

	scanner := bufio.NewScanner(file)
	ctl.dataReady = false
	go func() {
		refreshPeriod := time.Duration(ctl.GetConfig().ViewRefreshSeconds) * time.Second
		for !ctl.dataReady {
			time.Sleep(refreshPeriod)
			if !ctl.dataReady {
				ctl.view.Refresh()
			}
		}
	}()

	for scanner.Scan() {
		line := scanner.Text()
		currentLength := lengthExpandedTabs(line, ctl.conf.SpacesPerTab)
		ctl.data.AddLine(line)
		if currentLength > ctl.maxLineLength {
			ctl.maxLineLength = currentLength
		}
		if ctl.data.Len() <= height {
			ctl.view.Refresh()
		}
	}
	ctl.dataReady = true
	ctl.view.GetStatusBar().SafeStatus(view.StatusReady)
	ctl.view.Refresh()
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func lengthExpandedTabs(line string, tabSpaces int) int {
	total := utf8.RuneCountInString(line)
	tabs := strings.Count(line, "\t")
	return (total - tabs) + tabs * tabSpaces
}