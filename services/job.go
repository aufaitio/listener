package services

import (
	"github.com/aufaitio/listener/app"
	"github.com/aufaitio/listener/models"
)

// jobDAO specifies the interface of the job DAO needed by JobService.
type jobDAO interface {
	// Get returns the job with the specified job ID.
	Get(rs app.RequestScope, id int64) (*models.Job, error)
	// Count returns the number of repositories.
	Count(rs app.RequestScope) (int64, error)
	// Query returns the list of repositories with the given offset and limit.
	Query(rs app.RequestScope, offset, limit int) ([]*models.Job, error)
	// Create saves a new job in the storage.
	Create(rs app.RequestScope, job *models.Job) error
	// Update updates the job with given ID in the storage.
	Update(rs app.RequestScope, id int64, job *models.Job) error
	// Delete removes the job with given ID from the storage.
	Delete(rs app.RequestScope, id int64) error
}

// JobService provides services related with repositories.
type JobService struct {
	dao jobDAO
}

// NewJobService creates a new JobService with the given job DAO.
func NewJobService(dao jobDAO) *JobService {
	return &JobService{dao}
}

// Get returns the job with the specified the job ID.
func (s *JobService) Get(rs app.RequestScope, id int64) (*models.Job, error) {
	return s.dao.Get(rs, id)
}

// Create creates a new job.
func (s *JobService) Create(rs app.RequestScope, model *models.Job) (*models.Job, error) {
	if err := model.Validate(); err != nil {
		return nil, err
	}
	if err := s.dao.Create(rs, model); err != nil {
		return nil, err
	}
	return s.dao.Get(rs, model.ID)
}

// Update updates the job with the specified ID.
func (s *JobService) Update(rs app.RequestScope, id int64, model *models.Job) (*models.Job, error) {
	if err := model.Validate(); err != nil {
		return nil, err
	}
	if err := s.dao.Update(rs, id, model); err != nil {
		return nil, err
	}
	return s.dao.Get(rs, id)
}

// Delete deletes the job with the specified ID.
func (s *JobService) Delete(rs app.RequestScope, id int64) (*models.Job, error) {
	job, err := s.dao.Get(rs, id)
	if err != nil {
		return nil, err
	}
	err = s.dao.Delete(rs, id)
	return job, err
}

// Count returns the number of repositories.
func (s *JobService) Count(rs app.RequestScope) (int64, error) {
	return s.dao.Count(rs)
}

// Query returns the repositories with the specified offset and limit.
func (s *JobService) Query(rs app.RequestScope, offset, limit int) ([]*models.Job, error) {
	return s.dao.Query(rs, offset, limit)
}
