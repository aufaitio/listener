package apis

import (
	"github.com/go-ozzo/ozzo-routing"
	"github.com/quantumew/data-access/models"
	"github.com/quantumew/listener/app"
)

type (
	// jobService specifies the interface for the repository service needed by jobResource.
	jobService interface {
		Get(rs app.RequestScope, name string) (*models.Job, error)
		Query(rs app.RequestScope, offset, limit int) ([]*models.Job, error)
		Count(rs app.RequestScope) (int64, error)
		Create(rs app.RequestScope, model *models.Job) (*models.Job, error)
		CreateJobsFromHook(rs app.RequestScope, hook *models.NpmHook) ([]*models.Job, error)
		Update(rs app.RequestScope, name string, model *models.Job) (*models.Job, error)
		Delete(rs app.RequestScope, name string) (*models.Job, error)
	}

	// jobResource defines the handlers for the CRUD APIs.
	jobResource struct {
		service    jobService
		repService repositoryService
	}
)

// ServeJobResource sets up the routing of repository endpoints and the corresponding handlers.
func ServeJobResource(rg *routing.RouteGroup, service jobService, repService repositoryService) {
	r := &jobResource{service, repService}
	// Some of these routes are probably pointless but building it like a standard REST service
	rg.Get("/jobs/<name>", r.get)
	rg.Get("/jobs", r.query)
	rg.Post("/jobs", r.create)
	rg.Put("/jobs/<name>", r.update)
	rg.Delete("/jobs/name>", r.delete)
}

func (r *jobResource) get(c *routing.Context) error {
	response, err := r.service.Get(app.GetRequestScope(c), c.Param("name"))
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
	var model models.NpmHook
	if err := c.Read(&model); err != nil {
		return err
	}

	response, err := r.service.CreateJobsFromHook(app.GetRequestScope(c), &model)

	if err != nil {
		return err
	}

	return c.Write(response)
}

func (r *jobResource) update(c *routing.Context) error {
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

func (r *jobResource) delete(c *routing.Context) error {
	response, err := r.service.Delete(app.GetRequestScope(c), c.Param("name"))
	if err != nil {
		return err
	}

	return c.Write(response)
}
