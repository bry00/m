package utl

import (
	"strings"
	"testing"
)

type TripleInt struct {
	a int
	b int
	c int
}

func TestMinInt(t *testing.T) {
	values := []TripleInt{
		{2, 2, 3},
		{-3, 2, -3},
		{2, 100, 2},
		{200, 1200, 200},
	}
	for _, v := range values {
		got := MinInt(v.b, v.c)
		if got != v.a {
			t.Errorf("MinInt(%d, %d) = %d; want %d", v.b, v.c, got, v.a)
		}
	}
}

func TestMaxInt(t *testing.T) {
	values := []TripleInt{
		{3, 2, 3},
		{2, 2, -3},
		{100, 100, 2},
		{1200, 1200, 200},
	}
	for _, v := range values {
		got := MaxInt(v.b, v.c)
		if got != v.a {
			t.Errorf("MaxInt(%d, %d) = %d; want %d", v.b, v.c, got, v.a)
		}
	}
}

func TestCountRunesAtIndex(t *testing.T) {
	values := []struct {
		Expected int
		Str      string
		Index    int
	}{
		{7, "Ala ma kota", 7},
		{7, "Chrząszcz", 8},
		{8, "półciężarówka", 12},
		{5, "ĄąŚśćĆŻżŹź", 10},
	}
	for _, v := range values {
		got := CountRunesAtIndex(v.Str, v.Index)
		if got != v.Expected {
			t.Errorf("CountRunesAtIndex(\"%s\", %d) = %d; want %d", v.Str, v.Index, got, v.Expected)
		}
	}
}

func TestR2x(t *testing.T) {
	values := []struct {
		Expected int
		Str      string
		Index    int
	}{
		{7, "Ala ma kota", 7},
		{9, "Chrząszcz", 8},
		{17, "półciężarówka", 12},
		{18, "ĄąŚśćĆŻżŹź", 9},
	}
	for _, v := range values {
		got := R2x(v.Str, v.Index)
		if got != v.Expected {
			t.Errorf("R2x(\"%s\", %d) = %d; want %d", v.Str, v.Index, got, v.Expected)
		}
	}
}

func TestIsEmptyString(t *testing.T) {
	values := []struct {
		Expected bool
		Str      string
	}{
		{false, "Ala ma kota"},
		{false, "\tok"},
		{false, " true "},
		{true, ""},
		{true, "   "},
		{true, " \n\t \t \n  \r   "},
	}
	for _, v := range values {
		got := IsEmptyString(v.Str)
		if got != v.Expected {
			t.Errorf("IsEmptyString(\"%s\") = %v; want %v", v.Str, got, v.Expected)
		}
	}
}

func TestRemoveBackspaces(t *testing.T) {
	values := []struct {
		Expected string
		Str      string
	}{
		{"Ala ma kota", "Ala ma kota"},
		{"Alamakota", "Ala \bma \bkota"},
		{"One, two", "Oo\bn\bne\be, t\btw\bwo\bo"},
	}
	for _, v := range values {
		got := RemoveBackspaces(v.Str)
		if got != v.Expected {
			t.Errorf("RemoveBackspaces(\"%s\") = \"%s\"; want \"%s\"", v.Str, got, v.Expected)
		}
	}
}

func TestOnEachFieldWithSuffix(t *testing.T) {
	type TestStruct struct {
		LeftText  string
		RightText string
		Name      string
		NewText   string
	}
	v := TestStruct{
		LeftText:  "on the left   ",
		RightText: "   on the right",
		Name:      "neutral",
		NewText:   "    brand new   ",
	}
	OnEachFieldWithSuffix(&v, "Text",
		func(s string) string { return strings.ToUpper(strings.TrimSpace(s)) })

	results := []struct {
		Expected string
		Got      string
	}{
		{"ON THE LEFT", v.LeftText},
		{"ON THE RIGHT", v.RightText},
		{"neutral", v.Name},
		{"BRAND NEW", v.NewText},
	}
	for _, v := range results {
		if v.Got != v.Expected {
			t.Errorf("OnEachFieldWithSuffix ==> \"%s\"; want \"%s\"", v.Got, v.Expected)
		}
	}
}
