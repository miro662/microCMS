package mircocms

import "testing"

type testdata struct {
}

//TestDefaultTemplate tests if there is default template
func TestDefaultTemplate(t *testing.T) {
	tmpl := Templates["default.html"]
	if tmpl == nil {
		t.Errorf("Cannot find default template")
	}
}
