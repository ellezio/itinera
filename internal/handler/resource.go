package handler

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/ellezio/itinera/internal/resource"
	"github.com/ellezio/itinera/web/templates/components"
	"github.com/ellezio/itinera/web/templates/pages"
)

type ResourceHandler struct {
	resources *resource.ResourceService
}

func NewResourceHandler(resourceService *resource.ResourceService) *ResourceHandler {
	return &ResourceHandler{resourceService}
}

func (rh *ResourceHandler) Page(w http.ResponseWriter, r *http.Request) {
	resources, err := rh.resources.GetAll()
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	tags, err := rh.resources.GetTags()
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	statuses, err := rh.resources.GetStatuses()
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	pages.Resources(resources, tags, statuses).Render(r.Context(), w)
}

func (rh *ResourceHandler) Create(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		slog.Error(err.Error())
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	title := r.FormValue("title")
	source := r.FormValue("source")
	status := r.FormValue("status")
	tags := r.Form["tags"]

	statusID, _ := strconv.ParseInt(status, 10, 64)

	tagIDs := make([]int64, 0)
	for _, t := range tags {
		tid, _ := strconv.ParseInt(t, 10, 64)
		tagIDs = append(tagIDs, tid)
	}

	if err := rh.resources.Create(title, source, statusID, tagIDs); err != nil {
		slog.Error(err.Error())
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	resources, err := rh.resources.GetAll()
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	statuses, err := rh.resources.GetStatuses()
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	components.ResourceList(resources, statuses).Render(r.Context(), w)
}

func (rh *ResourceHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	rid, _ := strconv.ParseInt(id, 10, 64)

	if err := rh.resources.Delete(rid); err != nil {
		slog.Error(err.Error())
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
}

func (rh *ResourceHandler) ChangeStatus(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	rid, _ := strconv.ParseInt(id, 10, 64)

	status := r.FormValue("status")
	statusID, _ := strconv.ParseInt(status, 10, 64)

	if err := rh.resources.SetStatus(rid, statusID); err != nil {
		slog.Error(err.Error())
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
