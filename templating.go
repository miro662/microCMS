package microcms

import (
	"html/template"
	"os"
	"path/filepath"
)

// Templates holds all templates that can be used by microCMS
var Templates map[string]*template.Template

// templatingFunctions holds all functions that can be used with templatingFunctions
var templatingFunctions template.FuncMap

func init() {
	Templates = make(map[string]*template.Template)
	templatingFunctions = template.FuncMap{
		"page": PageByRoute,
	}

	// Walk trought /templates directory and load all templates
	filepath.Walk("templates", parseTemplates)

}

func parseTemplates(path string, info os.FileInfo, err error) error {
	if !info.IsDir() {
		// Load template
		tmpl := template.New(info.Name())
		tmpl = tmpl.Funcs(templatingFunctions)
		tmpl, err := tmpl.ParseFiles(path)
		if err != nil {
			panic("Cannot parse template " + path)
		}
		Templates[info.Name()] = tmpl
	}
	return nil
}
