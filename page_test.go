package microcms

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"html/template"
	"testing"
)

var testPage = Page{
	template: sql.NullString{String: "test", Valid: true},
	Data:     "",
}

//TestSimplePage tests rendering of simple page
func TestSimplePage(t *testing.T) {
	err := json.Unmarshal([]byte(`{
        "Title": "Title",
        "Content": "Content"
    }`), &testPage.Data)
	if err != nil {
		t.Errorf(err.Error())
	}
	Templates["test"], _ = template.New("test").Parse("<h1>{{.Data.Title}}</h1><p>{{.Data.Content}}</p>")
	buf := new(bytes.Buffer)
	err = testPage.Render(buf)
	if err != nil {
		t.Errorf("Rendering error: %v", err.Error())
	}
	str := buf.String()
	target := "<h1>Title</h1><p>Content</p>"
	if str != target {
		t.Errorf("Wrong result\nExpected: %v\nGot: %v", target, str)
	}
}
