package microcms

import (
	"errors"
	"io"
)

// Renderable interface describes something that can be rendered
type Renderable interface {
	// Render renders given object to io.Writer
	Render(io.Writer) error
}

// ErrNotRenderable is an error that should be thrown when page is not renderable
var ErrNotRenderable = errors.New("This page is not renderable")

// ErrTemplateNotFound is an error that should be thrown when template with given name is not found
var ErrTemplateNotFound = errors.New("Template not found")

// Render renders given page to io.Writer
func (page *Page) Render(w io.Writer) error {
	// If page has not assigned template, it cannot be rendered
	if !page.template.Valid {
		return ErrNotRenderable
	}

	// Render assigned template
	err := Template.ExecuteTemplate(w, page.template.String, page)
	return err
}
