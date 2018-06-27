package services

import (
	"errors"
	"testing"

	"github.com/aufaitio/listener/app"
	"github.com/aufaitio/listener/models"
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
	repository, err := s.Create(nil, &models.Repository{
		Name: "ddd",
	})
	if assert.Nil(t, err) && assert.NotNil(t, repository) {
		assert.Equal(t, 4, repository.ID)
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
	repository, err := s.Update(nil, 2, &models.Repository{
		Name: "ddd",
	})
	if assert.Nil(t, err) && assert.NotNil(t, repository) {
		assert.Equal(t, 2, repository.ID)
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
		assert.Equal(t, 2, repository.ID)
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

func newMockRepositoryDAO() repositoryDAO {
	return &mockRepositoryDAO{
		records: []models.Repository{
			{ID: 1, Name: "aaa"},
			{ID: 2, Name: "bbb"},
			{ID: 3, Name: "ccc"},
		},
	}
}

type mockRepositoryDAO struct {
	records []models.Repository
}

func (m *mockRepositoryDAO) Get(rs app.RequestScope, id int) (*models.Repository, error) {
	for _, record := range m.records {
		if record.ID == id {
			return &record, nil
		}
	}
	return nil, errors.New("not found")
}

func (m *mockRepositoryDAO) Query(rs app.RequestScope, offset, limit int) ([]models.Repository, error) {
	return m.records[offset : offset+limit], nil
}

func (m *mockRepositoryDAO) Count(rs app.RequestScope) (int, error) {
	return len(m.records), nil
}

func (m *mockRepositoryDAO) Create(rs app.RequestScope, repository *models.Repository) error {
	if repository.ID != 0 {
		return errors.New("Id cannot be set")
	}
	repository.ID = len(m.records) + 1
	m.records = append(m.records, *repository)
	return nil
}

func (m *mockRepositoryDAO) Update(rs app.RequestScope, id int, repository *models.Repository) error {
	repository.ID = id
	for i, record := range m.records {
		if record.ID == id {
			m.records[i] = *repository
			return nil
		}
	}
	return errors.New("not found")
}

func (m *mockRepositoryDAO) Delete(rs app.RequestScope, id int) error {
	for i, record := range m.records {
		if record.ID == id {
			m.records = append(m.records[:i], m.records[i+1:]...)
			return nil
		}
	}
	return errors.New("not found")
}
