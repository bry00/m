package buffers

import (
	"bufio"
	"io"
	"net/http"
	"os"
	"path"
	"strings"
	"testing"
)

const testdataDir = "../../testdata"
const testFileName = "testfile.txt"
const testDataURL = "https://wolnelektury.pl/media/book/txt/pan-tadeusz.txt"
const testBlockSize = 1 * KB
const testTotalSize = 8 * KB

var testFilePath string
var theBuffer *BufferedData

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

func reverse(str string) string {
	s := strings.Split(str, "\n")
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	return strings.Join(s, "\n")
}

func TestMain(m *testing.M) {
	theBuffer = NewBufferedData(testBlockSize, testTotalSize)
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
	code := m.Run()

	os.Exit(code)
}

func TestLen(t *testing.T) {
	expected := 10845
	got := theBuffer.Len()
	if got != expected {
		t.Errorf("BufferedData.Len ==> %d; want %d", got, expected)
	}
}

func TestGetLine(t *testing.T) {
	values := []struct {
		Index    int
		Expected string
	}{
		{10, "Księga pierwsza"},
		{1775, "Na Litwie, chwała Bogu, stare obyczaje."},
		{10750, "Kraj lat dziecinnych! On zawsze zostanie"},
		{5016, "Nawet śpią muchy."},
		{1633, "Siedzieli przeciw sobie mrukliwi i gniewni."},
		{45, "Świeciły się z daleka pobielane ściany,"},
	}
	i := theBuffer.NewLineIndexer()
	for _, v := range values {
		i.IndexSet(v.Index, false)
		if i.IndexOK() {
			got, _ := i.GetLine()
			if got != v.Expected {
				t.Errorf("BufferedData.GetLine ==> \"%s\"; want \"%s\"", got, v.Expected)
			}
		} else {
			t.Errorf("Invalid index")
		}
	}
}

func TestIndexIncrement(t *testing.T) {
	values := []struct {
		Index    int
		Rows     int
		Expected string
	}{
		{10750, 5, "Kraj lat dziecinnych! On zawsze zostanie\nŚwięty i czysty jak pierwsze kochanie,\nNiezaburzony błędów przypomnieniem,\nNiepodkopany nadziei złudzeniem,\nAni zmieniony wypadków strumieniem."},
		{18, 4, "Litwo! Ojczyzno moja! ty jesteś jak zdrowie:\nIle cię trzeba cenić, ten tylko się dowie,\nKto cię stracił. Dziś piękność twą w całej ozdobie\nWidzę i opisuję, bo tęsknię po tobie."},
		{8057, 6, "Znowu wzmaga się burza, ulewa nawalna\nI ciemność gruba, gęsta, prawie dotykalna.\nZnowu deszcz ciszej szumi, grom na chwilę uśnie;\nZnowu wzbudzi się, ryknie i znów wodą chluśnie.\nAż się uspokoiło wszystko; tylko drzewa\nSzumią około domu i szemrze ulewa."},
		{5211, 9, "Czy potrzeba, żebyśmy zaraz w pole wyszli?\nStrzelców zebrać, rzecz łatwa; prochu mam dostatek;\nW plebanii u księdza jest kilka armatek;\nPrzypominam, iż Jankiel mówił, iż u siebie\nMa groty do lanc, że je mogę wziąć w potrzebie;\nTe groty przywiózł w pakach gotowych z Królewca\nPod sekretem; weźmiem je, zaraz zrobim drzewca;\nSzabel nam nie zabraknie; szlachta na koń wsiędzie,\nJa z synowcem na czele, i — jakoś to będzie!»"},
	}
	i := theBuffer.NewLineIndexer()
	for _, v := range values {
		i.IndexSet(v.Index, false)

		lines := []string{}

		for row := 0; row < v.Rows && i.IndexOK(); _, row = i.IndexIncrement(), row+1 {
			if line, err := i.GetLine(); err == nil {
				lines = append(lines, strings.TrimSpace(line))
			}
		}

		got := strings.Join(lines, "\n")
		if got != v.Expected {
			t.Errorf("BufferedData.IndexIncrement/GetLine ==> \"%s\"; want \"%s\"", got, v.Expected)
		}
	}
}

func TestIndexDecrement(t *testing.T) {
	values := []struct {
		Index    int
		Rows     int
		Expected string
	}{
		{10754, 5, "Kraj lat dziecinnych! On zawsze zostanie\nŚwięty i czysty jak pierwsze kochanie,\nNiezaburzony błędów przypomnieniem,\nNiepodkopany nadziei złudzeniem,\nAni zmieniony wypadków strumieniem."},
		{21, 4, "Litwo! Ojczyzno moja! ty jesteś jak zdrowie:\nIle cię trzeba cenić, ten tylko się dowie,\nKto cię stracił. Dziś piękność twą w całej ozdobie\nWidzę i opisuję, bo tęsknię po tobie."},
		{8062, 6, "Znowu wzmaga się burza, ulewa nawalna\nI ciemność gruba, gęsta, prawie dotykalna.\nZnowu deszcz ciszej szumi, grom na chwilę uśnie;\nZnowu wzbudzi się, ryknie i znów wodą chluśnie.\nAż się uspokoiło wszystko; tylko drzewa\nSzumią około domu i szemrze ulewa."},
		{5219, 9, "Czy potrzeba, żebyśmy zaraz w pole wyszli?\nStrzelców zebrać, rzecz łatwa; prochu mam dostatek;\nW plebanii u księdza jest kilka armatek;\nPrzypominam, iż Jankiel mówił, iż u siebie\nMa groty do lanc, że je mogę wziąć w potrzebie;\nTe groty przywiózł w pakach gotowych z Królewca\nPod sekretem; weźmiem je, zaraz zrobim drzewca;\nSzabel nam nie zabraknie; szlachta na koń wsiędzie,\nJa z synowcem na czele, i — jakoś to będzie!»"},
	}
	i := theBuffer.NewLineIndexer()
	for _, v := range values {
		i.IndexSet(v.Index, false)

		lines := []string{}

		for row := 0; row < v.Rows && i.IndexOK(); _, row = i.IndexDecrement(), row+1 {
			if line, err := i.GetLine(); err == nil {
				lines = append(lines, strings.TrimSpace(line))
			}
		}
		expected := reverse(v.Expected)
		got := strings.Join(lines, "\n")
		if got != expected {
			t.Errorf("BufferedData.IndexDecrement/GetLine ==> \"%s\"; want \"%s\"", got, expected)
		}
	}
}
