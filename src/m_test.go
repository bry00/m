package main

import (
	"github.com/bry00/m/buffers"
	"path"
	"strings"
	"testing"
)

func TestGetProg(t *testing.T) {
	values := []struct {
		Expected string
		Str      string
	}{
		{"m", path.Join("home", "users", "someone", "bin", "m")},
		{"m", path.Join(".", "m")},
		{"m", "m"},
		{"m", "m.exe"},
		{"m", path.Join("home", "users", "someone", "bin", "m.exe")},
		{"m", path.Join("..", "..", "bin", "m")},
	}
	for _, v := range values {
		args := []string{v.Str}
		got := getProg(args)
		if got != v.Expected {
			t.Errorf("getProg(\"%s\") = \"%s\"; want \"%s\"", v.Str, got, v.Expected)
		}
	}
}

func TestComposeFilename(t *testing.T) {
	values := []struct {
		Expected string
		Args     []string
	}{
		{"m", []string{"m"}},
		{"m.exe", []string{"m.exe"}},
		{"the file", []string{"the", "file"}},
		{"/usr/local/data/very important file.txt", []string{"/usr/local/data/very", "important", "file.txt"}},
		{"/usr/local/data/very-important-file.txt", []string{"/usr/local/data/very-important-file.txt"}},
	}
	for _, v := range values {
		got := composeFileName(v.Args)
		if got != v.Expected {
			t.Errorf("composeFileName(\"%s\") = \"%s\"; want \"%s\"", strings.Join(v.Args, "\", \""), got, v.Expected)
		}
	}
}

func TestCheckDefaultValue(t *testing.T) {
	values := []struct {
		Expected int
		Val      int
		Config   int
		Default  int
	}{
		{buffers.DefaultBlockSizeLimit, -1, -1, buffers.DefaultBlockSizeLimit},
		{buffers.DefaultBlockSizeLimit, 0, 0, buffers.DefaultBlockSizeLimit},
		{4 * buffers.KB, 0, 4 * buffers.KB, buffers.DefaultBlockSizeLimit},
		{10 * buffers.MB, 0, 10 * buffers.MB, buffers.DefaultBlockSizeLimit},
	}
	for _, v := range values {
		got := v.Val
		checkDefaultValue(&got, v.Config, v.Default)
		if got != v.Expected {
			t.Errorf("checkDefaultValue(%d, %d, %d) = %d; want %d", v.Val, v.Config, v.Default, got, v.Expected)
		}
	}
}
