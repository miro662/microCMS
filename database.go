package microcms

import (
	"database/sql"
	"encoding/json"
	"errors"
)

const schema = `
    create table if not exists pages (
        id serial primary key,
        name varchar not null,
        parent int,
        template varchar,
        data json
    );

    create table if not exists custom_routes (
        route varchar not null,
        page_id int references pages
    );

    create or replace view routes_v as
    with recursive routes(route, page_id) as (
        select * from custom_routes
        union all
        select '/' || name, id
        from pages
        where parent is null
        union all
        select r.route || 
            CASE WHEN r.route = '/' THEN ''
                ELSE '/'
            END
        || p.name, p.id
        from routes r
        join pages p on (r.page_id = p.parent)
    )
    select distinct route, page_id from routes;
`

// ErrPageNotFound is returned when page is not found in database
var ErrPageNotFound = errors.New("Page not found")

// Schema creates database schema
func Schema(db *sql.DB) error {
	_, err := db.Exec(schema)
	return err
}

func fromRow(row *sql.Row) (Page, error) {
	var page = Page{}
	var jsonData []byte
	err := row.Scan(&page.id, &page.name, &page.parent, &page.template, &jsonData)
	if err != nil {
		if err == sql.ErrNoRows {
			return Page{}, ErrPageNotFound
		} else {
			return Page{}, err
		}
	}
	err = json.Unmarshal(jsonData, &page.data)
	return page, err
}

// PageByID gets page which has given ID
func PageByID(id int, db *sql.DB) (Page, error) {
	row := db.QueryRow("SELECT * FROM pages WHERE id = $1", id)
	return fromRow(row)
}

// PageByRoute gets page which has given route
func PageByRoute(route string, db *sql.DB) (Page, error) {
	row := db.QueryRow("SELECT p.* FROM pages p JOIN routes_v r ON (p.id = r.page_id) WHERE r.route = $1", route)
	return fromRow(row)
}

// Parent returns page's parent
func (p *Page) Parent(db *sql.DB) (Page, error) {
	row := db.QueryRow("SELECT p.* FROM pages c JOIN pages p ON (p.id = c.parent) WHERE c.id = $1", p.id)
	return fromRow(row)
}

// Children returns page's children
func (p *Page) Children(db *sql.DB) ([]Page, error) {
	rows, err := db.Query("SELECT * FROM pages WHERE parent = $1", p.id)
	defer rows.Close()
	if err != nil {
		return []Page{}, err
	}

	pages := make([]Page, 0)
	for rows.Next() {
		var page = Page{}
		var jsonData []byte
		err := rows.Scan(&page.id, &page.name, &page.parent, &page.template, &jsonData)
		if err != nil {
			return []Page{}, err
		}
		err = json.Unmarshal(jsonData, &page.data)
		if err != nil {
			return []Page{}, err
		}
		pages = append(pages, page)
	}
	if len(pages) == 0 {
		return []Page{}, ErrPageNotFound
	}
	return pages, nil
}
