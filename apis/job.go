package apis

import (
	"strconv"

	"github.com/aufaitio/listener/app"
	"github.com/aufaitio/listener/models"
	"github.com/go-ozzo/ozzo-routing"
)

type (
	// jobService specifies the interface for the repository service needed by jobResource.
	jobService interface {
		Get(rs app.RequestScope, id int64) (*models.Job, error)
		Query(rs app.RequestScope, offset, limit int) ([]*models.Job, error)
		Count(rs app.RequestScope) (int64, error)
		Create(rs app.RequestScope, model *models.Job) (*models.Job, error)
		Update(rs app.RequestScope, id int64, model *models.Job) (*models.Job, error)
		Delete(rs app.RequestScope, id int64) (*models.Job, error)
	}

	// jobResource defines the handlers for the CRUD APIs.
	jobResource struct {
		service jobService
	}
)

// ServeJobResource sets up the routing of repository endpoints and the corresponding handlers.
func ServeJobResource(rg *routing.RouteGroup, service jobService) {
	r := &jobResource{service}
	rg.Get("/jobs/<id>", r.get)
	rg.Get("/jobs", r.query)
	rg.Post("/jobs", r.create)
	rg.Put("/jobs/<id>", r.update)
	rg.Delete("/jobs/<id>", r.delete)
}

func (r *jobResource) get(c *routing.Context) error {
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

func (r *jobResource) query(c *routing.Context) error {
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

func (r *jobResource) create(c *routing.Context) error {
	var model models.Job
	if err := c.Read(&model); err != nil {
		return err
	}
	response, err := r.service.Create(app.GetRequestScope(c), &model)
	if err != nil {
		return err
	}

	return c.Write(response)
}

func (r *jobResource) update(c *routing.Context) error {
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

func (r *jobResource) delete(c *routing.Context) error {
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
