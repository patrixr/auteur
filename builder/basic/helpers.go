package basic

import (
	"html/template"
)

var helpers = template.FuncMap{
	"len": func(array []any) int {
		return len(array)
	},
	"errorstr": func(e error) string {
		return e.Error()
	},
	"gt": func(a, b int) bool {
		return a > b
	},
}
