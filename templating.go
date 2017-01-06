package mircocms

import (
	"html/template"
	"os"
	"path/filepath"
)

// Templates holds all templates that can be used by microCMS
var Templates map[string]*template.Template

func init() {
	Templates = make(map[string]*template.Template)

	// Walk trought /templates directory and load all templates
	filepath.Walk("templates", parseTemplates)
}

func parseTemplates(path string, info os.FileInfo, err error) error {
	if !info.IsDir() {
		// Load template
		tmpl := template.New(info.Name())
		tmpl, err := tmpl.ParseFiles(path)
		if err != nil {
			panic("Cannot parse template " + path)
		}
		Templates[info.Name()] = tmpl
	}
	return nil
}
