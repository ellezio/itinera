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

func (rs *ResourceService) Create(title, source string, status int64, tags []int64) error {
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

	return err
}

type FullResource struct {
	Resource db.Resource
	Status   db.Status
	Tags     []db.Tag
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
	tags := make(map[int64][]db.Tag)
	for _, t := range dbTags {
		tags[t.ResourceID] = append(tags[t.ResourceID], t.Tag)
	}

	resources := make([]FullResource, len(dbResources))
	for i, r := range dbResources {
		resources[i] = FullResource{
			Resource: r.Resource,
			Status:   r.Status,
			Tags:     tags[r.Resource.ID],
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
