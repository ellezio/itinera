package main

import (
	"database/sql"
	_ "embed"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/ellezio/itinera/internal/db"
	"github.com/ellezio/itinera/internal/handler"
	"github.com/ellezio/itinera/internal/resource"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	port := flag.String("port", "8080", "")
	dsn := flag.String("dsn", "file::memory:", "")
	flag.Parse()

	sqldb, err := sql.Open("sqlite3", *dsn+"?_fk=on")
	if err != nil {
		log.Fatal(err)
	}

	ddl, _ := os.ReadFile("internal/db/schema/schema.sql")
	if _, err := sqldb.Exec(string(ddl)); err != nil {
		log.Fatal(err)
	}

	dml := `
	INSERT INTO statuses VALUES (1, 'pending'), (2, 'inprogress'), (3, 'done');
	INSERT INTO tags VALUES (1, 'go'), (2, 'rust'), (3, 'c');
	`
	if _, err := sqldb.Exec(dml); err != nil {
		log.Fatal(err)
	}

	queries := db.New(sqldb)

	resourceService := resource.NewResourceService(queries)
	resourceHandler := handler.NewResourceHandler(resourceService)

	mux := http.NewServeMux()
	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.Dir("web/static"))))
	mux.HandleFunc("GET /", resourceHandler.Page)
	mux.HandleFunc("GET /resources", resourceHandler.Page)
	mux.HandleFunc("POST /resources", resourceHandler.Create)
	mux.HandleFunc("DELETE /resources/{id}", resourceHandler.Delete)
	mux.HandleFunc("POST /resources/{id}/status", resourceHandler.ChangeStatus)
	mux.HandleFunc("GET /resources/{id}/edit", resourceHandler.EditPage)
	mux.HandleFunc("POST /resources/{id}", resourceHandler.Edit)

	fmt.Printf("listening on %s\n", *port)
	if err := http.ListenAndServe(":"+*port, mux); err != nil {
		log.Fatal(err)
	}
}
