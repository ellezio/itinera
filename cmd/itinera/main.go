package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"io"
	"net/http"

	"github.com/a-h/templ"
	"github.com/ellezio/itinera/internal/db"
	"github.com/ellezio/itinera/web/templates/layout"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	port := flag.String("port", "8080", "")
	dsn := flag.String("dsn", "file::memory:", "")
	flag.Parse()

	sqldb, err := sql.Open("sqlite3", *dsn+"?_fk=on")
	if err != nil {
		panic(err)
	}

	_ = db.New(sqldb)

	mux := http.NewServeMux()
	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.Dir("web/static"))))
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		contents := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
			_, err := io.WriteString(w, "<h1>Hello world</h1>")
			return err
		})
		ctx := templ.WithChildren(r.Context(), contents)
		layout.Base("Hello world").Render(ctx, w)
	})

	fmt.Printf("listening on %s\n", *port)
	if err := http.ListenAndServe(":"+*port, mux); err != nil {
		fmt.Println(err)
	}
}
