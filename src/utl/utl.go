package utl

import (
	"fmt"
	"reflect"
	"strings"
	"unicode/utf8"
)

func MinInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func MaxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func IsEmptyString(str string) bool {
	return len(strings.TrimSpace(str)) == 0
}

func OnEachStringField(structure interface{}, filter func(string)bool, transform func(string)string) int {
	result := 0
	theStruct := reflect.ValueOf(structure)
	stype := theStruct.Elem()
	meta := reflect.Indirect(theStruct).Type()
	for i := 0; i < stype.NumField(); i++ {
		f := stype.Field(i)
		if f.Kind() == reflect.Struct {
			if f.CanAddr() {
				result += OnEachStringField(f.Addr().Interface(), filter, transform)
			}
		} else {
			if f.Kind() == reflect.String && f.CanSet() {
				if filter(meta.Field(i).Name) {
					f.SetString(transform(f.String()))
					result++
				}
			}
		}
	}
	return result
}

func OnEachFieldWithSuffix(structure interface{}, suffix string, transform func(string)string) int {
	return OnEachStringField(structure, func(fname string) bool {return strings.HasSuffix(fname, suffix)}, transform)

}

func RemoveBackspaces(str string) string {
	result := make([]rune, len(str))
	length := 0
	for _, r := range str {
		if r == '\b' {
			if length > 0 {
				length--
			}
		} else {
			result[length] = r
			length++
		}
	}
	return string(result[:length])
}

func SubstrFrom(str string, start int) string {
	i := 0
	for r := range str {
		if i == start {
			return str[r:]
		}
		i++
	}
	return ""
}

func Substr(str string, start int, end int) string {
	start_idx := 0
	i := 0
	for r := range str {
		if i == start {
			start_idx = r
		}
		if i == end {
			return str[start_idx:r]
		}
		i++
	}
	return str[start_idx:]
}

func R2x(str string, index int) int {
	i := 0
	for r := range str {
		if i == index {
			return r
		}
		i++
	}
	panic(fmt.Sprintf("Wrong runes index %d versus length of %d", index, utf8.RuneCountInString(str)))
}

func CountRunesAtIndex(str string, index int) int {
	return utf8.RuneCountInString(str[0:index])
}