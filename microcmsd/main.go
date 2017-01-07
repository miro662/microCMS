package main

import (
	"log"
	"net/http"
	"os"

	"database/sql"

	"github.com/miro662/microcms"

	"path"

	_ "github.com/lib/pq"
)

var fileServer http.Handler

func main() {
	// Get working directory
	if len(os.Args) == 1 {
		microcms.Dir = "."
	} else {
		microcms.Dir = os.Args[1]
	}

	microcms.LoadTemplates()

	// Get enviroment variables
	addr := os.Getenv("MICROCMSD_ADDR")
	if addr == "" {
		addr = ":8080"
	}
	dbAddr := os.Getenv("MICROCMSD_DB")
	if dbAddr == "" {
		log.Fatal("Enviroment variable MICROCMSD_DB empty; set MICROCMSD_DB to PostgreSQL connection string")
	}

	// Try to connect to database
	var err error
	microcms.Db, err = sql.Open("postgres", dbAddr)
	if err != nil {
		log.Fatalf("Cannot connect to database: %v", err.Error())
	}

	// Apply DB schema
	err = microcms.Schema()
	if err != nil {
		log.Fatalf("Schema eror: %v", err.Error())
	}

	// Apply static files serving
	fs := http.FileServer(http.Dir(path.Join(microcms.Dir, "static")))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Apply main handler
	http.HandleFunc("/", handler)

	// Listen
	http.ListenAndServe(addr, nil)
}

func stripRoute(route string) string {
	if route != "/" {
		if route[(len(route)-1):] == "/" {
			route = string(route[:len(route)-1])
		}
	}
	return route
}

func handler(res http.ResponseWriter, req *http.Request) {
	log.Printf("[%s] %q from %v\n", req.Method, req.URL.String(), req.RemoteAddr)
	// Search for page with given route
	page, err := microcms.PageByRoute(stripRoute(req.RequestURI))
	if err == nil && page != nil {
		// Show page
		err := page.Render(res)
		if err != nil {
			// Error rendering page :(
			errorHandler(res, req, 500)
			log.Printf("Rendering error: %v\n", err.Error())
		}
	} else {
		if err == nil && page == nil {
			// Show 404
			errorHandler(res, req, 404)
		} else {
			errorHandler(res, req, 500)
			log.Printf("HTTP hanlder error: %v\n", err.Error())
		}
	}
}

func errorHandler(res http.ResponseWriter, req *http.Request, status int) {
	// Write status header
	res.WriteHeader(status)
	// If 404, show 404 template
	if status == 404 {
		microcms.Template.ExecuteTemplate(res, "404.html", nil)
	}
}
