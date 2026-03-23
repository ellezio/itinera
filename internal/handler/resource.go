package handler

import (
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/ellezio/itinera/internal/db"
	"github.com/ellezio/itinera/internal/resource"
	"github.com/ellezio/itinera/web/templates/pages"
	resourceView "github.com/ellezio/itinera/web/templates/resource"
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

	resourceView.Page(resources, tags, statuses).Render(r.Context(), w)
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

	// resources, err := rh.resources.GetAll()
	// if err != nil {
	// 	slog.Error(err.Error())
	// 	http.Error(w, "", http.StatusInternalServerError)
	// 	return
	// }
	//
	// statuses, err := rh.resources.GetStatuses()
	// if err != nil {
	// 	slog.Error(err.Error())
	// 	http.Error(w, "", http.StatusInternalServerError)
	// 	return
	// }

	redirect(w, r, "/resources", http.StatusFound)
	// components.ResourceList(resources, statuses).Render(r.Context(), w)
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

func (rh *ResourceHandler) EditPage(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	rid, _ := strconv.ParseInt(id, 10, 64)

	if r.Header.Get("Hx-Request") == "true" {
		w.Header().Add("Hx-Redirect", fmt.Sprintf("/resources/%d/edit", rid))
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
	if rid == 0 {
		pages.Resources(tags, statuses).Render(r.Context(), w)
		return
	}

	resource, err := rh.resources.GetResource(rid)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Redirect(w, r, "/resources", http.StatusFound)
			return
		}

		slog.Error(err.Error())
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	pages.ResourceEdit(resource, tags, statuses).Render(r.Context(), w)
}

func (rh *ResourceHandler) Edit(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	rid, _ := strconv.ParseInt(id, 10, 64)

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

	if err := rh.resources.Edit(rid, title, source, statusID, tagIDs); err != nil {
		slog.Error(err.Error())
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	w.Header().Add("Hx-Redirect", "/resources")
	w.WriteHeader(http.StatusFound)
}

func (rh *ResourceHandler) Info(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	rid, _ := strconv.ParseInt(id, 10, 64)

	rsrc, err := rh.resources.GetResource(rid)
	if err != nil {
		slog.Error(err.Error())
		redirect(w, r, "/resources", http.StatusFound)
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
	pages.Resource(rsrc, tags, statuses)
	resourceView.ResourceInfoPane(rsrc).Render(r.Context(), w)
}

func (rh *ResourceHandler) ResourceNoteEditBox(w http.ResponseWriter, r *http.Request) {
	noteID_str := r.PathValue("note_id")
	noteID, _ := strconv.ParseInt(noteID_str, 10, 64)

	resourceID_str := r.PathValue("resource_id")
	resourceID, _ := strconv.ParseInt(resourceID_str, 10, 64)

	note := db.Note{EntityType: "resource", EntityID: resourceID}
	if noteID > 0 {
		var err error
		note, err = rh.resources.GetNote(noteID)
		if err != nil {
			slog.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	resourceView.Note(note, true).Render(r.Context(), w)
}

func (rh *ResourceHandler) GetResourceNote(w http.ResponseWriter, r *http.Request) {
	noteID_str := r.PathValue("note_id")
	noteID, _ := strconv.ParseInt(noteID_str, 10, 64)

	note, err := rh.resources.GetNote(noteID)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	resourceView.Note(note, false).Render(r.Context(), w)
}

func (rh *ResourceHandler) EditResourceNote(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	noteID_str := r.PathValue("note_id")
	noteID, _ := strconv.ParseInt(noteID_str, 10, 64)

	resourceID_str := r.PathValue("resource_id")
	resourceID, _ := strconv.ParseInt(resourceID_str, 10, 64)

	title := r.FormValue("title")
	content := r.FormValue("content")

	var note db.Note
	if noteID > 0 {
		note, _ = rh.resources.UpdateNote(noteID, title, content)
	} else {
		note, _ = rh.resources.AddNote(title, content, resourceID)
	}
	resourceView.Note(note, false).Render(r.Context(), w)
}

func (rh *ResourceHandler) DeleteResourceNote(w http.ResponseWriter, r *http.Request) {
	noteID_str := r.PathValue("note_id")
	noteID, _ := strconv.ParseInt(noteID_str, 10, 64)

	err := rh.resources.DeleteNote(noteID)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
