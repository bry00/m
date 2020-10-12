package view

import "testing"

func TestActionsNamesDefinition(t *testing.T) {
	expected := lastAction + 1
	got := ActionUnknown.Count()

	if got != expected {
		t.Errorf("Number of action names = %d; but should be %d", got, expected)
	}
}

func TestActionString(t *testing.T) {
	values := []struct {
		ActionInt    int
		ActionString string
	}{
		{-1, "unknown"},
		{0, "unknown"},
		{1000, "unknown"},
	}
	for _, v := range values {
		got := Action(v.ActionInt).String()
		expected := v.ActionString
		if got != expected {
			t.Errorf("Action(%d).String() => \"%s\"; want \"%s\"", v.ActionInt, got, expected)
		}
	}
}

func TestActionNames(t *testing.T) {
	for i := 0; i < ActionUnknown.Count(); i++ {
		got := Action(i).String()
		expected := actionNames[i]
		if got != expected {
			t.Errorf("Action(%d).String() => \"%s\"; want \"%s\"", i, got, expected)
		}
	}
}

func TestAppStatusString(t *testing.T) {
	values := []struct {
		AppStatusInt    int
		AppStatusString string
	}{
		{-1, "unknown"},
		{0, "unknown"},
		{1, "ready"},
		{2, "reading"},
		{3, "receiving"},
		{4, "unknown"},
	}
	for _, v := range values {
		got := AppStatus(v.AppStatusInt).String()
		expected := v.AppStatusString
		if got != expected {
			t.Errorf("AppStatus(%d).String() => \"%s\"; want \"%s\"", v.AppStatusInt, got, expected)
		}
	}
}

func TestAppStatusDisplay(t *testing.T) {
	values := []struct {
		AppStatusInt    int
		AppStatusString string
	}{
		{-1, ""},
		{0, ""},
		{1, "READY"},
		{2, "Reading..."},
		{3, "Receiving..."},
		{4, ""},
	}
	for _, v := range values {
		got := AppStatus(v.AppStatusInt).Display()
		expected := v.AppStatusString
		if got != expected {
			t.Errorf("AppStatus(%d).Display() => \"%s\"; want \"%s\"", v.AppStatusInt, got, expected)
		}
	}
}
