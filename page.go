package microcms

import (
	"database/sql"
)

// Page describes single page
type Page struct {
	id       int
	name     string
	parent   sql.NullInt64
	template sql.NullString
	Data     interface{}
}
