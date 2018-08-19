package services

import (
	"errors"
	"github.com/mongodb/mongo-go-driver/mongo"
	"testing"

	"github.com/quantumew/data-access"
	"github.com/quantumew/data-access/models"
	"github.com/stretchr/testify/assert"
)

func TestNewRepositoryService(t *testing.T) {
	dao := newMockRepositoryDAO()
	s := NewRepositoryService(dao)
	assert.Equal(t, dao, s.dao)
}

func TestRepositoryService_Get(t *testing.T) {
	s := NewRepositoryService(newMockRepositoryDAO())
	repository, err := s.Get(new(MockRequestScope), 1)
	if assert.Nil(t, err) && assert.NotNil(t, repository) {
		assert.Equal(t, "aaa", repository.Name)
	}

	repository, err = s.Get(new(MockRequestScope), 100)
	assert.NotNil(t, err)
}

func TestRepositoryService_Create(t *testing.T) {
	s := NewRepositoryService(newMockRepositoryDAO())
	repository, err := s.Create(new(MockRequestScope), createRepository("ddd", "testing", "1.1.1", "1.2.3"))
	if assert.Nil(t, err) && assert.NotNil(t, repository) {
		assert.Equal(t, int64(4), repository.ID)
		assert.Equal(t, "ddd", repository.Name)
	}

	// dao error
	_, err = s.Create(new(MockRequestScope), &models.Repository{
		ID:   100,
		Name: "ddd",
	})
	assert.NotNil(t, err)

	// validation error
	_, err = s.Create(new(MockRequestScope), &models.Repository{
		Name: "",
	})
	assert.NotNil(t, err)
}

func TestRepositoryService_Update(t *testing.T) {
	s := NewRepositoryService(newMockRepositoryDAO())
	repository, err := s.Update(new(MockRequestScope), 2, createRepository("ddd", "a", "1.2.4", "1.2.3"))
	if assert.Nil(t, err) && assert.NotNil(t, repository) {
		assert.Equal(t, int64(2), repository.ID)
		assert.Equal(t, "ddd", repository.Name)
	}

	// dao error
	_, err = s.Update(new(MockRequestScope), 100, &models.Repository{
		Name: "ddd",
	})
	assert.NotNil(t, err)

	// validation error
	_, err = s.Update(new(MockRequestScope), 2, &models.Repository{
		Name: "",
	})
	assert.NotNil(t, err)
}

func TestRepositoryService_Delete(t *testing.T) {
	s := NewRepositoryService(newMockRepositoryDAO())
	repository, err := s.Delete(new(MockRequestScope), 2)
	if assert.Nil(t, err) && assert.NotNil(t, repository) {
		assert.Equal(t, int64(2), repository.ID)
		assert.Equal(t, "bbb", repository.Name)
	}

	_, err = s.Delete(new(MockRequestScope), 2)
	assert.NotNil(t, err)
}

func TestRepositoryService_Query(t *testing.T) {
	s := NewRepositoryService(newMockRepositoryDAO())
	result, err := s.Query(new(MockRequestScope), 1, 2)
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

func (m *mockRepositoryDAO) Get(db *mongo.Database, id int64) (*models.Repository, error) {
	for _, record := range m.records {
		if record.ID == id {
			return record, nil
		}
	}
	return nil, errors.New("not found")
}

func (m *mockRepositoryDAO) Query(db *mongo.Database, offset, limit int) ([]*models.Repository, error) {
	return m.records[offset : offset+limit], nil
}

func (m *mockRepositoryDAO) QueryByDependency(db *mongo.Database, dependencyName string) ([]*models.Repository, error) {
	return []*models.Repository{}, nil
}

func (m *mockRepositoryDAO) Count(db *mongo.Database) (int64, error) {
	return int64(len(m.records)), nil
}

func (m *mockRepositoryDAO) Create(db *mongo.Database, repository *models.Repository) error {
	if repository.ID != 0 {
		return errors.New("Id cannot be set")
	}
	repository.ID = int64(len(m.records) + 1)

	m.records = append(m.records, repository)
	return nil
}

func (m *mockRepositoryDAO) Update(db *mongo.Database, id int64, repository *models.Repository) error {
	repository.ID = id
	for i, record := range m.records {
		if record.ID == id {
			m.records[i] = repository
			return nil
		}
	}
	return errors.New("not found")
}

func (m *mockRepositoryDAO) Delete(db *mongo.Database, id int64) error {
	for i, record := range m.records {
		if record.ID == id {
			m.records = append(m.records[:i], m.records[i+1:]...)
			return nil
		}
	}
	return errors.New("not found")
}
