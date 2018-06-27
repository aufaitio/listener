package services

import (
	"github.com/aufaitio/listener/app"
	"github.com/aufaitio/listener/models"
)

// repositoryDAO specifies the interface of the repository DAO needed by RepositoryService.
type repositoryDAO interface {
	// Get returns the repository with the specified repository ID.
	Get(rs app.RequestScope, id int64) (*models.Repository, error)
	// Count returns the number of repositories.
	Count(rs app.RequestScope) (int64, error)
	// Query returns the list of repositories with the given offset and limit.
	Query(rs app.RequestScope, offset, limit int) ([]*models.Repository, error)
	// Create saves a new repository in the storage.
	Create(rs app.RequestScope, repository *models.Repository) error
	// Update updates the repository with given ID in the storage.
	Update(rs app.RequestScope, id int64, repository *models.Repository) error
	// Delete removes the repository with given ID from the storage.
	Delete(rs app.RequestScope, id int64) error
}

// RepositoryService provides services related with repositories.
type RepositoryService struct {
	dao repositoryDAO
}

// NewRepositoryService creates a new RepositoryService with the given repository DAO.
func NewRepositoryService(dao repositoryDAO) *RepositoryService {
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
