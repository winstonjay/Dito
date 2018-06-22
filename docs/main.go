package main

import (
	"dito/src/eval"
	"html/template"
	"os"
)

type page struct {
	Title    string
	Builtins map[string]map[string]string
}

func main() {
	b := make(map[string]map[string]string)
	for name, fn := range eval.Builtins {
		b[name] = fn.Doc()
	}
	p := page{Title: "Builtins", Builtins: b}
	t, _ := template.ParseFiles("templates/builtins.html")
	t.Execute(os.Stdout, p)
}
