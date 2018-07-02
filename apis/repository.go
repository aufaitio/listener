package apis

import (
	"strconv"

	"github.com/aufaitio/listener/app"
	"github.com/aufaitio/listener/models"
	"github.com/go-ozzo/ozzo-routing"
)

type (
	// repositoryService specifies the interface for the repository service needed by repositoryResource.
	repositoryService interface {
		Get(rs app.RequestScope, id int64) (*models.Repository, error)
		Query(rs app.RequestScope, offset, limit int) ([]*models.Repository, error)
		Count(rs app.RequestScope) (int64, error)
		Create(rs app.RequestScope, model *models.Repository) (*models.Repository, error)
		Update(rs app.RequestScope, id int64, model *models.Repository) (*models.Repository, error)
		Delete(rs app.RequestScope, id int64) (*models.Repository, error)
	}

	// repositoryResource defines the handlers for the CRUD APIs.
	repositoryResource struct {
		service repositoryService
	}
)

// ServeRepositoryResource sets up the routing of repository endpoints and the corresponding handlers.
func ServeRepositoryResource(rg *routing.RouteGroup, service repositoryService) {
	r := &repositoryResource{service}
	// Some of these routes are probably pointless but building it like a standard REST service
	rg.Get("/repositories/<id>", r.get)
	rg.Get("/repositories", r.query)
	// Sort of a post/put
	rg.Post("/repositories", r.create)
	rg.Put("/repositories/<id>", r.update)
	rg.Delete("/repositories/<id>", r.delete)
}

func (r *repositoryResource) get(c *routing.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return err
	}

	response, err := r.service.Get(app.GetRequestScope(c), id)
	if err != nil {
		return err
	}

	return c.Write(response)
}

func (r *repositoryResource) query(c *routing.Context) error {
	rs := app.GetRequestScope(c)
	count, err := r.service.Count(rs)
	if err != nil {
		return err
	}
	paginatedList := getPaginatedListFromRequest(c, count)
	items, err := r.service.Query(app.GetRequestScope(c), paginatedList.Offset(), paginatedList.Limit())
	if err != nil {
		return err
	}
	paginatedList.Items = items
	return c.Write(paginatedList)
}

func (r *repositoryResource) create(c *routing.Context) error {
	var model models.Repository
	if err := c.Read(&model); err != nil {
		return err
	}
	response, err := r.service.Create(app.GetRequestScope(c), &model)
	if err != nil {
		return err
	}

	return c.Write(response)
}

func (r *repositoryResource) update(c *routing.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return err
	}

	rs := app.GetRequestScope(c)

	model, err := r.service.Get(rs, id)
	if err != nil {
		return err
	}

	if err := c.Read(model); err != nil {
		return err
	}

	response, err := r.service.Update(rs, id, model)
	if err != nil {
		return err
	}

	return c.Write(response)
}

func (r *repositoryResource) delete(c *routing.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return err
	}

	response, err := r.service.Delete(app.GetRequestScope(c), id)
	if err != nil {
		return err
	}

	return c.Write(response)
}
