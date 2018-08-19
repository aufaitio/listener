package apis

import (
	"strconv"

	"github.com/go-ozzo/ozzo-routing"
	"github.com/quantumew/data-access/models"
	"github.com/quantumew/listener/app"
)

type (
	// repositoryService specifies the interface for the repository service needed by repositoryResource.
	repositoryService interface {
		Get(rs app.RequestScope, name string) (*models.Repository, error)
		Query(rs app.RequestScope, offset, limit int) ([]*models.Repository, error)
		Count(rs app.RequestScope) (int64, error)
		Create(rs app.RequestScope, model *models.Repository) (*models.Repository, error)
		Update(rs app.RequestScope, name string, model *models.Repository) (*models.Repository, error)
		Patch(rs app.RequestScope, modelList []*models.Repository) ([]*models.Repository, error)
		Delete(rs app.RequestScope, name string) (*models.Repository, error)
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
	rg.Get("/repositories/<name>", r.get)
	rg.Get("/repositories", r.query)
	rg.Post("/repositories", r.create)
	rg.Put("/repositories/<name>", r.update)
	rg.Patch("/repositories", r.patch)
	rg.Delete("/repositories/<name>", r.delete)
}

func (r *repositoryResource) get(c *routing.Context) error {
	name := c.Param("name")

	response, err := r.service.Get(app.GetRequestScope(c), name)
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
	name := c.Param("name")
	rs := app.GetRequestScope(c)

	model, err := r.service.Get(rs, name)
	if err != nil {
		return err
	}

	if err := c.Read(model); err != nil {
		return err
	}

	response, err := r.service.Update(rs, name, model)
	if err != nil {
		return err
	}

	return c.Write(response)
}

func (r *repositoryResource) patch(c *routing.Context) error {
	rs := app.GetRequestScope(c)
	var repoList []*models.Repository

	if err := c.Read(repoList); err != nil {
		return err
	}

	responseList, err := r.service.Patch(rs, repoList)

	if err != nil {
		return err
	}

	return c.Write(responseList)
}

func (r *repositoryResource) delete(c *routing.Context) error {
	name := c.Param("name")
	response, err := r.service.Delete(app.GetRequestScope(c), name)
	if err != nil {
		return err
	}

	return c.Write(response)
}
