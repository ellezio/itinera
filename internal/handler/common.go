package handler

import (
	"net/http"
	"net/url"
	"slices"

	"github.com/ellezio/itinera/internal/db"
	resourceView "github.com/ellezio/itinera/web/templates/resource"
)

func redirect(w http.ResponseWriter, r *http.Request, url string, status int) {
	if r.Header.Get("Hx-Request") == "true" {
		w.Header().Set("Hx-Redirect", url)
		w.WriteHeader(status)
	} else {
		http.Redirect(w, r, url, status)
	}
}

func prepareFilters(query url.Values, tags []db.Tag, statuses []db.Status) ([]resourceView.Filter[db.Tag], []resourceView.Filter[db.Status]) {
	ftags := make([]resourceView.Filter[db.Tag], len(tags))
	for i, t := range tags {
		ftags[i] = resourceView.Filter[db.Tag]{Data: t, Selected: slices.Contains(query["tag"], t.Name)}
	}

	fstatuses := make([]resourceView.Filter[db.Status], len(statuses))
	for i, s := range statuses {
		fstatuses[i] = resourceView.Filter[db.Status]{Data: s, Selected: slices.Contains(query["status"], s.Name)}
	}

	return ftags, fstatuses
}
