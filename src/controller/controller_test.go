package controller

import (
	"bufio"
	"github.com/bry00/m/buffers"
	"github.com/bry00/m/config"
	"github.com/bry00/m/view"
	"io"
	"net/http"
	"os"
	"path"
	"testing"
)

const testdataDir = "../../testdata"
const testFileName = "testfile.txt"
const testDataURL = "https://wolnelektury.pl/media/book/txt/pan-tadeusz.txt"
const testBlockSize = 4096
const testTotalSize = 16384

type DummyTestStatusBar struct {
}

type DummyTestView struct {
	ctl           view.TheViewController
	showNumbers   bool
	showRuler     bool
	rulerPosition int
	statusBar     *DummyTestStatusBar
}

var testFilePath string
var theBuffer *buffers.BufferedData
var theController *Controller

func (sb *DummyTestStatusBar) Reset()                                      {}
func (sb *DummyTestStatusBar) Message(format string, a ...interface{})     {}
func (sb *DummyTestStatusBar) Status(status view.AppStatus)                {}
func (sb *DummyTestStatusBar) SafeMessage(format string, a ...interface{}) {}
func (sb *DummyTestStatusBar) SafeStatus(status view.AppStatus)            {}

func (v *DummyTestView) ShowSearchResult(lineIndex int, start int, end int) {}

func (v *DummyTestView) ShowLine(lineIndex int) {}

func (v *DummyTestView) ShowNumbers(show bool) {
	v.showNumbers = show
}

func (v *DummyTestView) AreNumbersShown() bool {
	return v.showNumbers
}

func (v *DummyTestView) IsRulerShown() bool {
	return v.showRuler
}

func (v *DummyTestView) ShowRuler(show bool) {
	v.showRuler = show
}

func (v *DummyTestView) GetRulerPosition() int {
	return v.rulerPosition
}

func (v *DummyTestView) SetRulerPosition(index int) {
	v.rulerPosition = index
}

func (v *DummyTestView) GetStatusBar() view.TheStatusBar {
	return v.statusBar
}

func NewDummyTestView() *DummyTestView {
	return &DummyTestView{
		statusBar: &DummyTestStatusBar{},
	}
}

func (v *DummyTestView) SetController(ctl view.TheViewController) {
	v.ctl = ctl
}

func (v *DummyTestView) StopApplication()            {}
func (v *DummyTestView) DisplayAt(left int, top int) {}

func (v *DummyTestView) GetDisplayRect() (int, int, int, int) {
	return 0, 0, 0, 0
}

func (v *DummyTestView) Refresh()            {}
func (v *DummyTestView) ShowSearchDialog()   {}
func (v *DummyTestView) ShowGotoLineDialog() {}
func (v *DummyTestView) Prepare()            {}
func (v *DummyTestView) Show()               {}
func (v *DummyTestView) ShowShortcuts()      {}

func (v *DummyTestView) GetKeyShortcuts() map[view.Action][]string {
	return make(map[view.Action][]string)
}

func obtainTestFile(dir string, fileName string, url string) string {
	targetPath := path.Join(dir, fileName)
	if _, err := os.Stat(targetPath); os.IsNotExist(err) {
		os.MkdirAll(dir, os.ModePerm)
		if resp, err := http.Get(url); err != nil {
			panic(err)
		} else {
			defer resp.Body.Close()
			if f, err := os.Create(targetPath); err != nil {
				panic(err)
			} else {
				defer f.Close()
				_, err = io.Copy(f, resp.Body)
			}
		}
	}
	return targetPath
}

func TestMain(m *testing.M) {
	theBuffer = buffers.NewBufferedData(testBlockSize, testTotalSize)
	testFilePath = obtainTestFile(testdataDir, testFileName, testDataURL)

	if file, err := os.Open(testFilePath); err != nil {
		panic(err)
	} else {
		defer file.Close()
		reader := bufio.NewReader(file)
		for eof := false; !eof; {
			line, err := reader.ReadString('\n')
			if err != nil {
				if err == io.EOF {
					eof = true
				} else {
					panic(err)
				}
			} else {
				theBuffer.AddLine(line)
			}
		}
	}

	defer theBuffer.Close()
	theController = NewController(testFilePath, "test", theBuffer, NewDummyTestView(), config.NewDefaultConfig(),
		false)
	defer theController.OnExit()
	code := m.Run()

	os.Exit(code)
}

func TestLengthExpandedTabs(t *testing.T) {
	spacesPerTab := 4
	values := []struct {
		Expected int
		Str      string
	}{
		{11, "Ala ma kota"},
		{4, "\t"},
		{14, "\t \t \t"},
		{17, "Ala\tma\tkota"},
	}
	for _, v := range values {
		got := lengthExpandedTabs(v.Str, spacesPerTab)
		if got != v.Expected {
			t.Errorf("lengthExpandedTabs(\"%s\") = %d; want %d", v.Str, got, v.Expected)
		}
	}
}

func TestFileExists(t *testing.T) {
	filename := "controller.go"
	got := fileExists(filename)
	expected := true
	if got != expected {
		t.Errorf("fileExists(\"%s\") = %v; want %v", filename, got, expected)
	}
}

func TestFindNext(t *testing.T) {
	values := []struct {
		Regex        bool
		ExpectedLine int
		ExpectedText string
		StartLine    int
		SearchFor    string
	}{
		{false, 18, "    Litwo! Ojczyzno moja! ty jesteś jak zdrowie:", 1, "Litwo! Ojczyzno moja!"},
		{false, 302, "Okazały budową, poważny ogromem,", 1, "poważny ogromem"},
		{true, 314, "Tłumacząc, że gotyckiej są architektury;", 100, "gotyck[a-z]*"},
		{true, 3070, "Także z drzewa, gotyckiej naśladowstwo sztuki.", 315, "gotyck[a-z]*"},
	}
	for _, v := range values {
		theController.SetSearchText(v.SearchFor, v.Regex, false)
		foundLine, _, _, foundLineText, err := theController.findNext(v.StartLine, 0)
		if err != nil {
			t.Error(err)
		} else {
			if foundLine != v.ExpectedLine || foundLineText != v.ExpectedText {
				t.Errorf("findNext(\"%s\") => \"%s\", %d; want \"%s\", %d",
					v.SearchFor, foundLineText, foundLine, v.ExpectedText, v.ExpectedLine)
			}
		}
	}
}

func TestFindPrevious(t *testing.T) {
	values := []struct {
		Regex        bool
		ExpectedLine int
		ExpectedText string
		StartLine    int
		SearchFor    string
	}{
		{false, 18, "    Litwo! Ojczyzno moja! ty jesteś jak zdrowie:", 100, "Litwo! Ojczyzno moja!"},
		{false, 302, "Okazały budową, poważny ogromem,", 1000, "poważny ogromem"},
		{true, 314, "Tłumacząc, że gotyckiej są architektury;", 1000, "gotyck[a-z]*"},
		{true, 3070, "Także z drzewa, gotyckiej naśladowstwo sztuki.", 3100, "gotyck[a-z]*"},
	}
	for _, v := range values {
		theController.SetSearchText(v.SearchFor, v.Regex, false)
		foundLine, _, _, foundLineText, err := theController.findPrevious(v.StartLine, 0)
		if err != nil {
			t.Error(err)
		} else {
			if foundLine != v.ExpectedLine || foundLineText != v.ExpectedText {
				t.Errorf("findNext(\"%s\") => \"%s\", %d; want \"%s\", %d",
					v.SearchFor, foundLineText, foundLine, v.ExpectedText, v.ExpectedLine)
			}
		}
	}
}
