package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/miro662/microcms"

	_ "github.com/lib/pq"
)

func main() {
	var err error
	microcms.Db, err = sql.Open("postgres", "user=miroslav dbname=microcms sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	err = microcms.Schema()
	if err != nil {
		log.Fatal(err)
	}
	http.HandleFunc("/", getHandler(nil))
	http.ListenAndServe(":1488", nil)
}

func getHandler(db *sql.DB) func(res http.ResponseWriter, req *http.Request) {
	return func(res http.ResponseWriter, req *http.Request) {
		route := req.RequestURI
		if route != "/" {
			if route[(len(route)-1):] == "/" {
				route = string(route[:len(route)-1])
			}
		}
		page, err := microcms.PageByRoute(route)
		if err != nil {
			if err == microcms.ErrPageNotFound {
				http.Error(res, "Page not found: "+route, 404)
			} else {
				http.Error(res, err.Error(), 500)
				log.Fatal(err)
			}
		} else {
			err = page.Render(res)
			if err != nil {
				http.Error(res, err.Error(), 500)
				log.Fatal(err)
			}
		}
	}
}
