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
    select distinct on (route) route, page_id
	from routes p
	order by char_length(route);
`

// Db describes database connection used by model
var Db *sql.DB

// ErrPageNotFound is returned when page is not found in database
var ErrPageNotFound = errors.New("Page not found")

// Schema creates database schema
func Schema() error {
	_, err := Db.Exec(schema)
	return err
}

func fromRow(row *sql.Row) (Page, error) {
	var page = Page{}
	var jsonData []byte
	err := row.Scan(&page.id, &page.name, &page.parent, &page.template, &jsonData)
	if err != nil {
		if err == sql.ErrNoRows {
			return Page{}, ErrPageNotFound
		}
		return Page{}, err
	}
	err = json.Unmarshal(jsonData, &page.Data)
	return page, err
}

// PageByID gets page which has given ID
func PageByID(id int) (Page, error) {
	row := Db.QueryRow("SELECT * FROM pages WHERE id = $1", id)
	return fromRow(row)
}

// PageByRoute gets page which has given route
func PageByRoute(route string) (Page, error) {
	row := Db.QueryRow("SELECT p.* FROM pages p JOIN routes_v r ON (p.id = r.page_id) WHERE r.route = $1", route)
	return fromRow(row)
}

// Parent returns page's parent
func (p *Page) Parent() (Page, error) {
	row := Db.QueryRow("SELECT p.* FROM pages c JOIN pages p ON (p.id = c.parent) WHERE c.id = $1", p.id)
	return fromRow(row)
}

// Children returns page's children
func (p *Page) Children() ([]Page, error) {
	rows, err := Db.Query("SELECT * FROM pages WHERE parent = $1", p.id)
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
		err = json.Unmarshal(jsonData, &page.Data)
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

// Route returns page's route
func (p *Page) Route() (string, error) {
	row := Db.QueryRow("SELECT route FROM routes_v WHERE page_id = $1", p.id)
	var route string
	err := row.Scan(&route)
	return route, err
}
