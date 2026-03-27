package handler

import (
	"database/sql"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/ellezio/itinera/internal/db"
	"github.com/ellezio/itinera/internal/resource"
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

	rsrc, err := rh.resources.Create(title, source, statusID, tagIDs)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	fullRsrc, err := rh.resources.MakeFullResource(rsrc)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	resourceView.ResourceInfoPane(fullRsrc, nil, nil, false).Render(r.Context(), w)
	resourceView.ListResourcesCards([]resource.FullResource{fullRsrc}).Render(r.Context(), w)
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

	var resource resource.FullResource
	if rid > 0 {
		if resource, err = rh.resources.GetResource(rid); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				http.Redirect(w, r, "/resources", http.StatusFound)
				return
			}

			slog.Error(err.Error())
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
	}

	resourceView.ResourceInfoPane(resource, tags, statuses, true).Render(r.Context(), w)
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

	rsrc, err := rh.resources.Edit(rid, title, source, statusID, tagIDs)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	fullRsrc, err := rh.resources.MakeFullResource(rsrc)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	resourceView.ResourceInfoPane(fullRsrc, nil, nil, false).Render(r.Context(), w)
	resourceView.Card(fullRsrc).Render(r.Context(), w)
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
	resourceView.ResourceInfoPane(rsrc, tags, statuses, false).Render(r.Context(), w)
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

func (rh *ResourceHandler) CreateTag(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	name := r.FormValue("name")
	color := r.FormValue("color")

	tag, err := rh.resources.CreateTag(name, color)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resourceView.SideNavTagList([]db.Tag{tag}).Render(r.Context(), w)
}

func (rh *ResourceHandler) CreateStatus(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	name := r.FormValue("name")
	color := r.FormValue("color")

	status, err := rh.resources.CreateStatus(name, color)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resourceView.SideNavStatusList([]db.Status{status}).Render(r.Context(), w)
}

func (rh *ResourceHandler) GetStatusEdit(w http.ResponseWriter, r *http.Request) {
	statusID_str := r.PathValue("id")
	statusID, _ := strconv.ParseInt(statusID_str, 10, 64)

	var status db.Status

	if statusID > 0 {
		var err error
		status, err = rh.resources.GetStatus(statusID)
		if err != nil {
			slog.Error(err.Error())
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
	}

	resourceView.ListItemEdit(status.ID, status.Name, "status", status.Color).Render(r.Context(), w)
}

func (rh *ResourceHandler) GetTagEdit(w http.ResponseWriter, r *http.Request) {
	tagID_str := r.PathValue("id")
	tagID, _ := strconv.ParseInt(tagID_str, 10, 64)

	var tag db.Tag

	if tagID > 0 {
		var err error
		tag, err = rh.resources.GetTag(tagID)
		if err != nil {
			slog.Error(err.Error())
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
	}

	resourceView.ListItemEdit(tag.ID, tag.Name, "tag", tag.Color).Render(r.Context(), w)
}

func (rh *ResourceHandler) GetTag(w http.ResponseWriter, r *http.Request) {
	tagID_str := r.PathValue("id")
	tagID, _ := strconv.ParseInt(tagID_str, 10, 64)

	if tagID <= 0 {
		http.Error(w, "invalid tagID", http.StatusBadRequest)
		return
	}

	tag, err := rh.resources.GetTag(tagID)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	resourceView.ListItem(tag.ID, tag.Name, "#", tag.Color, 1).Render(r.Context(), w)
}

func (rh *ResourceHandler) GetStatus(w http.ResponseWriter, r *http.Request) {
	statusID_str := r.PathValue("id")
	statusID, _ := strconv.ParseInt(statusID_str, 10, 64)

	if statusID <= 0 {
		http.Error(w, "invalid statusID", http.StatusBadRequest)
		return
	}

	status, err := rh.resources.GetStatus(statusID)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	resourceView.ListItem(status.ID, status.Name, "", status.Color, 1).Render(r.Context(), w)
}

func (rh *ResourceHandler) EditTag(w http.ResponseWriter, r *http.Request) {
	tagID_str := r.PathValue("id")
	tagID, _ := strconv.ParseInt(tagID_str, 10, 64)

	if tagID <= 0 {
		http.Error(w, "invalid tagID", http.StatusBadRequest)
		return
	}

	r.ParseForm()
	name := r.FormValue("name")
	color := r.FormValue("color")

	tag, err := rh.resources.UpdateTag(tagID, name, color)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	resourceView.ListItem(tag.ID, tag.Name, "#", tag.Color, 1).Render(r.Context(), w)
}

func (rh *ResourceHandler) EditStatus(w http.ResponseWriter, r *http.Request) {
	statusID_str := r.PathValue("id")
	statusID, _ := strconv.ParseInt(statusID_str, 10, 64)

	if statusID <= 0 {
		http.Error(w, "invalid statusID", http.StatusBadRequest)
		return
	}

	r.ParseForm()
	name := r.FormValue("name")
	color := r.FormValue("color")

	status, err := rh.resources.UpdateStatus(statusID, name, color)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	resourceView.ListItem(status.ID, status.Name, "", status.Color, 1).Render(r.Context(), w)
}

func (rh *ResourceHandler) DeleteStatus(w http.ResponseWriter, r *http.Request) {
	statusID_str := r.PathValue("id")
	statusID, _ := strconv.ParseInt(statusID_str, 10, 64)

	if statusID <= 0 {
		http.Error(w, "invalid statusID", http.StatusBadRequest)
		return
	}

	if err := rh.resources.DeleteStatus(statusID); err != nil {
		slog.Error(err.Error())
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
}

func (rh *ResourceHandler) DeleteTag(w http.ResponseWriter, r *http.Request) {
	tagID_str := r.PathValue("id")
	tagID, _ := strconv.ParseInt(tagID_str, 10, 64)

	if tagID <= 0 {
		http.Error(w, "invalid tagID", http.StatusBadRequest)
		return
	}

	if err := rh.resources.DeleteTag(tagID); err != nil {
		slog.Error(err.Error())
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
}
