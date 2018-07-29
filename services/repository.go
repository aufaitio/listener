package services

import (
	"github.com/aufaitio/data-access"
	"github.com/aufaitio/data-access/models"
	"github.com/aufaitio/listener/app"
)

// RepositoryService provides services related with repositories.
type RepositoryService struct {
	dao access.RepositoryDAO
}

// NewRepositoryService creates a new RepositoryService with the given repository DAO.
func NewRepositoryService(dao access.RepositoryDAO) *RepositoryService {
	return &RepositoryService{dao}
}

// Get returns the repository with the specified the repository ID.
func (s *RepositoryService) Get(rs app.RequestScope, id int64) (*models.Repository, error) {
	return s.dao.Get(rs, id)
}

// Create creates a new repository.
func (s *RepositoryService) Create(rs app.RequestScope, model *models.Repository) (*models.Repository, error) {
	if err := model.Validate(); err != nil {
		return nil, err
	}
	if err := s.dao.Create(rs, model); err != nil {
		return nil, err
	}
	return s.dao.Get(rs, model.ID)
}

// Update updates the repository with the specified ID.
func (s *RepositoryService) Update(rs app.RequestScope, id int64, model *models.Repository) (*models.Repository, error) {
	if err := model.Validate(); err != nil {
		return nil, err
	}
	if err := s.dao.Update(rs, id, model); err != nil {
		return nil, err
	}
	return s.dao.Get(rs, id)
}

// Delete deletes the repository with the specified ID.
func (s *RepositoryService) Delete(rs app.RequestScope, id int64) (*models.Repository, error) {
	repository, err := s.dao.Get(rs, id)
	if err != nil {
		return nil, err
	}
	err = s.dao.Delete(rs, id)
	return repository, err
}

// Count returns the number of repositories.
func (s *RepositoryService) Count(rs app.RequestScope) (int64, error) {
	return s.dao.Count(rs)
}

// Query returns the repositories with the specified offset and limit.
func (s *RepositoryService) Query(rs app.RequestScope, offset, limit int) ([]*models.Repository, error) {
	return s.dao.Query(rs, offset, limit)
}
