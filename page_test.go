package mircocms

import (
	"bytes"
	"database/sql"
	"html/template"
	"testing"
)

type staticPageData struct {
	Title   string
	Content string
}

var testPage = Page{
	template: sql.NullString{String: "test", Valid: true},
	data: staticPageData{
		Title:   "Title",
		Content: "Content",
	},
}

//TestSimplePage tests rendering of simple page
func TestSimplePage(t *testing.T) {
	Templates["test"], _ = template.New("test").Parse("<h1>{{.Title}}</h1><p>{{.Content}}</p>")
	buf := new(bytes.Buffer)
	err := testPage.Render(buf)
	if err != nil {
		t.Errorf("Rendering error: %v", err.Error())
	}
	str := buf.String()
	target := "<h1>Title</h1><p>Content</p>"
	if str != target {
		t.Errorf("Wrong result\nExpected: %v\nGot: %v", target, str)
	}
}
