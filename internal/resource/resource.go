package resource

import (
	"context"

	"github.com/ellezio/itinera/internal/db"
)

type ResourceService struct {
	store *db.Queries
}

func NewResourceService(queries *db.Queries) *ResourceService {
	return &ResourceService{queries}
}

func (rs *ResourceService) Create(title, source string, status int64, tags []int64) (db.Resource, error) {
	ctx := context.Background()

	resource, err := rs.store.CreateResource(ctx, db.CreateResourceParams{
		Title:      title,
		Source:     source,
		SourceType: "http",
		StatusID:   status,
	})

	for _, t := range tags {
		_ = rs.store.SetTag(ctx, db.SetTagParams{
			ResourceID: resource.ID,
			TagID:      t,
		})
	}

	return resource, err
}

type FullResource struct {
	Resource db.Resource
	Status   db.Status
	Tags     []db.Tag
	Notes    []db.Note
}

func (rs *ResourceService) GetAll() ([]FullResource, error) {
	ctx := context.Background()

	dbResources, err := rs.store.GetResources(ctx)
	if err != nil {
		return nil, err
	}

	ids := make([]int64, len(dbResources))
	for i, r := range dbResources {
		ids[i] = r.Resource.ID
	}

	dbTags, err := rs.store.GetResourcesTags(ctx, ids)
	if err != nil {
		return nil, err
	}

	tags := make(map[int64][]db.Tag)
	for _, t := range dbTags {
		tags[t.ResourceID] = append(tags[t.ResourceID], t.Tag)
	}

	dbNotes, err := rs.store.GetResourcesNotes(ctx, db.GetResourcesNotesParams{Resources: ids, EntityType: "resource"})
	if err != nil {
		return nil, err
	}
	notes := make(map[int64][]db.Note)
	for _, n := range dbNotes {
		notes[n.EntityID] = append(notes[n.EntityID], n)
	}

	resources := make([]FullResource, len(dbResources))
	for i, r := range dbResources {
		resources[i] = FullResource{
			Resource: r.Resource,
			Status:   r.Status,
			Tags:     tags[r.Resource.ID],
			Notes:    notes[r.Resource.ID],
		}
	}
	return resources, nil
}

func (rs *ResourceService) Delete(id int64) error {
	return rs.store.DeleteResource(context.Background(), id)
}

func (rs *ResourceService) GetTags() ([]db.Tag, error) {
	return rs.store.GetTags(context.Background())
}

func (rs *ResourceService) GetStatuses() ([]db.Status, error) {
	return rs.store.GetStatuses(context.Background())
}

func (rs *ResourceService) SetStatus(resourcesID int64, statusID int64) error {
	return rs.store.SetStatus(context.Background(), db.SetStatusParams{ID: resourcesID, StatusID: statusID})
}

func (rs *ResourceService) Edit(id int64, title, source string, status int64, tags []int64) (db.Resource, error) {
	ctx := context.Background()

	resource, err := rs.store.UpdateResource(ctx, db.UpdateResourceParams{
		ID:         id,
		Title:      title,
		Source:     source,
		SourceType: "http",
		StatusID:   status,
	})

	_ = rs.store.ClearTags(ctx, id)
	for _, t := range tags {
		_ = rs.store.SetTag(ctx, db.SetTagParams{
			ResourceID: resource.ID,
			TagID:      t,
		})
	}

	return resource, err
}

func (rs *ResourceService) GetResource(id int64) (FullResource, error) {
	ctx := context.Background()
	resourceRow, err := rs.store.GetResource(ctx, id)
	if err != nil {
		return FullResource{}, err
	}

	tags, err := rs.store.GetResourceTags(ctx, id)
	if err != nil {
		return FullResource{}, err
	}

	notes, err := rs.store.GetNotes(ctx, db.GetNotesParams{EntityID: id, EntityType: "resource"})

	resource := FullResource{
		Resource: resourceRow.Resource,
		Status:   resourceRow.Status,
		Tags:     tags,
		Notes:    notes,
	}

	return resource, nil
}

func (rs *ResourceService) AddNote(title, content string, resourceID int64) (db.Note, error) {
	return rs.store.CreateNote(
		context.Background(),
		db.CreateNoteParams{
			Title:      title,
			Content:    content,
			EntityID:   resourceID,
			EntityType: "resource",
		},
	)
}

func (rs *ResourceService) UpdateNote(id int64, title, content string) (db.Note, error) {
	return rs.store.UpdateNote(
		context.Background(),
		db.UpdateNoteParams{
			ID:      id,
			Title:   title,
			Content: content,
		},
	)
}

func (rs *ResourceService) GetNote(id int64) (db.Note, error) {
	return rs.store.GetNote(context.Background(), id)
}

func (rs *ResourceService) DeleteNote(id int64) error {
	return rs.store.DeleteNote(context.Background(), id)
}

func (rs *ResourceService) MakeFullResource(rsrc db.Resource) (FullResource, error) {
	ctx := context.Background()

	status, err := rs.store.GetStatus(ctx, rsrc.StatusID)
	if err != nil {
		return FullResource{}, err
	}

	tags, err := rs.store.GetResourceTags(ctx, rsrc.ID)
	if err != nil {
		return FullResource{}, err
	}

	notes, err := rs.store.GetNotes(ctx, db.GetNotesParams{EntityID: rsrc.ID, EntityType: "resource"})

	resource := FullResource{
		Resource: rsrc,
		Status:   status,
		Tags:     tags,
		Notes:    notes,
	}

	return resource, nil
}

func (rs *ResourceService) CreateTag(name, color string) (db.Tag, error) {
	return rs.store.CreateTag(context.Background(), db.CreateTagParams{
		Name:  name,
		Color: color,
	})
}

func (rs *ResourceService) CreateStatus(name, color string) (db.Status, error) {
	return rs.store.CreateStatus(context.Background(), db.CreateStatusParams{
		Name:  name,
		Color: color,
	})
}

func (rs *ResourceService) GetStatus(id int64) (db.Status, error) {
	return rs.store.GetStatus(context.Background(), id)
}

func (rs *ResourceService) GetTag(id int64) (db.Tag, error) {
	return rs.store.GetTag(context.Background(), id)
}

func (rs *ResourceService) UpdateStatus(id int64, name, color string) (db.Status, error) {
	return rs.store.UpdateStatus(context.Background(), db.UpdateStatusParams{
		ID:    id,
		Name:  name,
		Color: color,
	})
}

func (rs *ResourceService) UpdateTag(id int64, name, color string) (db.Tag, error) {
	return rs.store.UpdateTag(context.Background(), db.UpdateTagParams{
		ID:    id,
		Name:  name,
		Color: color,
	})
}

func (rs *ResourceService) DeleteStatus(id int64) error {
	return rs.store.DeleteStatus(context.Background(), id)
}

func (rs *ResourceService) DeleteTag(id int64) error {
	return rs.store.DeleteTag(context.Background(), id)
}
