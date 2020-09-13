package utl

import (
	"reflect"
	"strings"
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
