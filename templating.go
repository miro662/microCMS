package microcms

import "html/template"

// Template hold loaded templates
var Template *template.Template

// templatingFunctions holds all functions that can be used with templatingFunctions
var templatingFunctions = template.FuncMap{
	"page": PageByRoute,
}

func init() {
	var err error
	// Create new template
	Template, err = template.ParseGlob("templates/*.html")
	if err != nil {
		panic(err)
	}
	Template = Template.Funcs(templatingFunctions)
}
