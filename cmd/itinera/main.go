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
	INSERT OR IGNORE INTO statuses VALUES (1, 'pending', 'gray'), (2, 'inprogress', 'blue'), (3, 'done', 'green');
	INSERT OR IGNORE INTO tags VALUES (1, 'go', 'blue'), (2, 'rust', 'orange'), (3, 'c', 'yellow');
	`
	if _, err := sqldb.Exec(dml); err != nil {
		log.Fatal(err)
	}

	queries := db.New(sqldb)

	resourceService := resource.NewResourceService(queries)
	resourceHandler := handler.NewResourceHandler(resourceService)

	mux := http.NewServeMux()
	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.Dir("web/static"))))

	mux.HandleFunc("GET /", resourceHandler.ResourcesPage)

	mux.HandleFunc("GET /resources", resourceHandler.ResourcesPage)
	mux.HandleFunc("POST /resources", resourceHandler.Create)

	mux.HandleFunc("GET /resources/{id}/edit", resourceHandler.EditPage)

	mux.HandleFunc("POST /resources/{id}/status", resourceHandler.ChangeStatus)

	mux.HandleFunc("GET /resources/{resource_id}/notes/{note_id}/edit", resourceHandler.ResourceNoteEditBox)
	mux.HandleFunc("GET /resources/{resource_id}/notes/{note_id}", resourceHandler.GetResourceNote)
	mux.HandleFunc("POST /resources/{resource_id}/notes/{note_id}", resourceHandler.EditResourceNote)
	mux.HandleFunc("DELETE /resources/{resource_id}/notes/{note_id}", resourceHandler.DeleteResourceNote)

	mux.HandleFunc("GET /resources/{id}", resourceHandler.Info)
	mux.HandleFunc("POST /resources/{id}", resourceHandler.Edit)
	mux.HandleFunc("DELETE /resources/{id}", resourceHandler.Delete)

	mux.HandleFunc("GET /tags/{id}", resourceHandler.GetTag)
	mux.HandleFunc("GET /tags/{id}/edit", resourceHandler.GetTagEdit)
	mux.HandleFunc("POST /tags", resourceHandler.CreateTag)
	mux.HandleFunc("POST /tags/{id}/edit", resourceHandler.EditTag)
	mux.HandleFunc("DELETE /tags/{id}", resourceHandler.DeleteTag)

	mux.HandleFunc("GET /statuses/{id}", resourceHandler.GetStatus)
	mux.HandleFunc("GET /statuses/{id}/edit", resourceHandler.GetStatusEdit)
	mux.HandleFunc("POST /statuses", resourceHandler.CreateStatus)
	mux.HandleFunc("POST /statuses/{id}/edit", resourceHandler.EditStatus)
	mux.HandleFunc("DELETE /statuses/{id}", resourceHandler.DeleteStatus)

	mux.HandleFunc("GET /collections", resourceHandler.CollectionsPage)
	mux.HandleFunc("GET /collections/{collection_id}", resourceHandler.Collection)
	mux.HandleFunc("GET /collections/{collection_id}/edit", resourceHandler.CollectionEdit)
	mux.HandleFunc("POST /collections", resourceHandler.CollectionCreate)

	fmt.Printf("listening on %s\n", *port)
	if err := http.ListenAndServe(":"+*port, mux); err != nil {
		log.Fatal(err)
	}
}
