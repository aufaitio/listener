package services

import (
	"github.com/quantumew/data-access"
	"github.com/quantumew/data-access/models"
	"github.com/quantumew/listener/app"
)

// RepositoryService provides services related with repositories.
type RepositoryService struct {
	dao access.RepositoryDAO
}

// NewRepositoryService creates a new RepositoryService with the given repository DAO.
func NewRepositoryService(dao access.RepositoryDAO) *RepositoryService {
	return &RepositoryService{dao}
}

// Get returns the repository with the specified the repository name.
func (s *RepositoryService) Get(rs app.RequestScope, name string) (*models.Repository, error) {
	return s.dao.Get(rs.DB(), name)
}

// Create creates a new repository.
func (s *RepositoryService) Create(rs app.RequestScope, model *models.Repository) (*models.Repository, error) {
	if err := model.Validate(); err != nil {
		return nil, err
	}
	if err := s.dao.Create(rs.DB(), model); err != nil {
		return nil, err
	}
	return s.dao.Get(rs.DB(), model.ID)
}

// Update updates the repository with the specified name.
func (s *RepositoryService) Update(rs app.RequestScope, name string, model *models.Repository) (*models.Repository, error) {
	if err := model.Validate(); err != nil {
		return nil, err
	}
	if err := s.dao.Update(rs.DB(), name, model); err != nil {
		return nil, err
	}
	return s.dao.Get(rs.DB(), name)
}

// Patch bulk update of repositories
func (s *RepositoryService) Patch(rs app.RequestScope, repoList []*models.Repository) ([]*models.Repository, error) {
	for _, model := range repoList {
		if err := model.Validate(); err != nil {
			return nil, err
		}
	}

	if err := s.dao.Patch(rs.DB(), repoList); err != nil {
		return nil, err
	}

	var repoNameList []string

	for _, repo := range repoList {
		append(repoNameList, repoList.Name)
	}

	return s.dao.QueryByName(rs.DB(), repoNameList)
}

// Delete deletes the repository with the specified name.
func (s *RepositoryService) Delete(rs app.RequestScope, name string) (*models.Repository, error) {
	repository, err := s.dao.Get(rs.DB(), name)
	if err != nil {
		return nil, err
	}
	err = s.dao.Delete(rs.DB(), name)
	return repository, err
}

// Count returns the number of repositories.
func (s *RepositoryService) Count(rs app.RequestScope) (int64, error) {
	return s.dao.Count(rs.DB())
}

// Query returns the repositories with the specified offset and limit.
func (s *RepositoryService) Query(rs app.RequestScope, offset, limit int) ([]*models.Repository, error) {
	return s.dao.Query(rs.DB(), offset, limit)
}
