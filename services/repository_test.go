package services

import (
	"errors"
	"testing"

	"github.com/aufaitio/data-access"
	"github.com/aufaitio/data-access/models"
	"github.com/stretchr/testify/assert"
)

func TestNewRepositoryService(t *testing.T) {
	dao := newMockRepositoryDAO()
	s := NewRepositoryService(dao)
	assert.Equal(t, dao, s.dao)
}

func TestRepositoryService_Get(t *testing.T) {
	s := NewRepositoryService(newMockRepositoryDAO())
	repository, err := s.Get(nil, 1)
	if assert.Nil(t, err) && assert.NotNil(t, repository) {
		assert.Equal(t, "aaa", repository.Name)
	}

	repository, err = s.Get(nil, 100)
	assert.NotNil(t, err)
}

func TestRepositoryService_Create(t *testing.T) {
	s := NewRepositoryService(newMockRepositoryDAO())
	repository, err := s.Create(nil, createRepository("ddd", "testing", "1.1.1", "1.2.3"))
	if assert.Nil(t, err) && assert.NotNil(t, repository) {
		assert.Equal(t, int64(4), repository.ID)
		assert.Equal(t, "ddd", repository.Name)
	}

	// dao error
	_, err = s.Create(nil, &models.Repository{
		ID:   100,
		Name: "ddd",
	})
	assert.NotNil(t, err)

	// validation error
	_, err = s.Create(nil, &models.Repository{
		Name: "",
	})
	assert.NotNil(t, err)
}

func TestRepositoryService_Update(t *testing.T) {
	s := NewRepositoryService(newMockRepositoryDAO())
	repository, err := s.Update(nil, 2, createRepository("ddd", "a", "1.2.4", "1.2.3"))
	if assert.Nil(t, err) && assert.NotNil(t, repository) {
		assert.Equal(t, int64(2), repository.ID)
		assert.Equal(t, "ddd", repository.Name)
	}

	// dao error
	_, err = s.Update(nil, 100, &models.Repository{
		Name: "ddd",
	})
	assert.NotNil(t, err)

	// validation error
	_, err = s.Update(nil, 2, &models.Repository{
		Name: "",
	})
	assert.NotNil(t, err)
}

func TestRepositoryService_Delete(t *testing.T) {
	s := NewRepositoryService(newMockRepositoryDAO())
	repository, err := s.Delete(nil, 2)
	if assert.Nil(t, err) && assert.NotNil(t, repository) {
		assert.Equal(t, int64(2), repository.ID)
		assert.Equal(t, "bbb", repository.Name)
	}

	_, err = s.Delete(nil, 2)
	assert.NotNil(t, err)
}

func TestRepositoryService_Query(t *testing.T) {
	s := NewRepositoryService(newMockRepositoryDAO())
	result, err := s.Query(nil, 1, 2)
	if assert.Nil(t, err) {
		assert.Equal(t, 2, len(result))
	}
}

func createRepository(name string, depName string, depVersion string, installed string) *models.Repository {
	return &models.Repository{
		Name:   name,
		Config: models.Config{Branch: "master", Remote: "stuff"},
		Dependencies: []models.Dependency{
			models.Dependency{Name: depName, Semver: depVersion, Installed: installed},
		},
	}
}

func newMockRepositoryDAO() access.RepositoryDAO {
	repositoryList := []*models.Repository{
		createRepository("aaa", "test", "1.2.3", "1.0.0"),
		createRepository("bbb", "test", "2.2.3", "2.2.3"),
		createRepository("ccc", "test", "3.2.3", "3.2.3"),
	}

	for i, repository := range repositoryList {
		repository.ID = int64(i) + 1
	}
	return &mockRepositoryDAO{records: repositoryList}
}

type mockRepositoryDAO struct {
	records []*models.Repository
}

func (m *mockRepositoryDAO) Get(rs access.Scope, id int64) (*models.Repository, error) {
	for _, record := range m.records {
		if record.ID == id {
			return record, nil
		}
	}
	return nil, errors.New("not found")
}

func (m *mockRepositoryDAO) Query(rs access.Scope, offset, limit int) ([]*models.Repository, error) {
	return m.records[offset : offset+limit], nil
}

func (m *mockRepositoryDAO) QueryByDependency(rs access.Scope, dependencyName string) ([]*models.Repository, error) {
	return []*models.Repository{}, nil
}

func (m *mockRepositoryDAO) Count(rs access.Scope) (int64, error) {
	return int64(len(m.records)), nil
}

func (m *mockRepositoryDAO) Create(rs access.Scope, repository *models.Repository) error {
	if repository.ID != 0 {
		return errors.New("Id cannot be set")
	}
	repository.ID = int64(len(m.records) + 1)

	m.records = append(m.records, repository)
	return nil
}

func (m *mockRepositoryDAO) Update(rs access.Scope, id int64, repository *models.Repository) error {
	repository.ID = id
	for i, record := range m.records {
		if record.ID == id {
			m.records[i] = repository
			return nil
		}
	}
	return errors.New("not found")
}

func (m *mockRepositoryDAO) Delete(rs access.Scope, id int64) error {
	for i, record := range m.records {
		if record.ID == id {
			m.records = append(m.records[:i], m.records[i+1:]...)
			return nil
		}
	}
	return errors.New("not found")
}
