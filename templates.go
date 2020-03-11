package main

import (
	"html/template"
	"path"
	"sync"
)

type templateName string

const (
	List   templateName = "list.html"
	Single templateName = "single.html"
	Edit   templateName = "edit.html"
)

var templates = []templateName{
	List, Single, Edit,
}

func createTemplates() map[templateName]*template.Template {
	out := make(map[templateName]*template.Template, len(templates))

	for _, tmplName := range templates {
		out[tmplName] = template.Must(
			template.New("MyTemplate").ParseFiles(path.Join("templates", string(tmplName))))
	}

	return out
}

var mu sync.Mutex

func getTemplate(mapka map[templateName]*template.Template, key templateName) *template.Template {
	mu.Lock()
	defer mu.Unlock()

	val, ok := mapka[key]
	if !ok {
		return nil
	}

	return val
}
