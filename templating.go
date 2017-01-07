package microcms

import (
	"html/template"
	"path"
)

// Template hold loaded templates
var Template *template.Template

// Dir contains main microCMS directory
var Dir string

// templatingFunctions holds all functions that can be used with templatingFunctions
var templatingFunctions = template.FuncMap{
	"page": PageByRoute,
	"root": Root,
}

//LoadTemplates loads templates
func LoadTemplates() {
	var err error
	// Create new template
	Template, err = template.New("").Funcs(templatingFunctions).ParseGlob(path.Join(Dir, "templates/*.html"))
	if err != nil {
		panic(err)
	}
	Template = Template.Funcs(templatingFunctions)
}
