package tv

import "testing"

func TestNumberString(t *testing.T) {
	values := []struct {
		Val      int
		Width    int
		Expected string
	}{
		{0, 4, "[::d]000[::b]0"},
		{1, 5, "[::d]0000[::b]1"},
		{74, 6, "[::d]0000[::b]74"},
		{38475, 7, "[::d]00[::b]38475"},
		{754, 5, "[::d]00[::b]754"},
		{847248, 3, "[::b]248"},
	}
	for _, v := range values {
		got := numberString(v.Val, v.Width)
		expected := v.Expected
		if got != expected {
			t.Errorf("numberString(%d, %d) => \"%s\"; want \"%s\"", v.Val, v.Width, got, expected)
		}
	}
}
