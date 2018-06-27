package services

import (
	"errors"
	"testing"

	"github.com/aufaitio/listener/app"
	"github.com/aufaitio/listener/models"
	"github.com/stretchr/testify/assert"
)

func TestNewJobService(t *testing.T) {
	dao := newMockJobDAO()
	s := NewJobService(dao)
	assert.Equal(t, dao, s.dao)
}

func TestJobService_Get(t *testing.T) {
	s := NewJobService(newMockJobDAO())
	job, err := s.Get(nil, 1)
	if assert.Nil(t, err) && assert.NotNil(t, job) {
		assert.Equal(t, "aaa", job.Name)
	}

	job, err = s.Get(nil, 100)
	assert.NotNil(t, err)
}

func TestJobService_Create(t *testing.T) {
	s := NewJobService(newMockJobDAO())
	job, err := s.Create(nil, &models.Job{
		Name: "ddd",
	})
	if assert.Nil(t, err) && assert.NotNil(t, job) {
		assert.Equal(t, 4, job.ID)
		assert.Equal(t, "ddd", job.Name)
	}

	// dao error
	_, err = s.Create(nil, &models.Job{
		Id:   100,
		Name: "ddd",
	})
	assert.NotNil(t, err)

	// validation error
	_, err = s.Create(nil, &models.Job{
		Name: "",
	})
	assert.NotNil(t, err)
}

func TestJobService_Update(t *testing.T) {
	s := NewJobService(newMockJobDAO())
	job, err := s.Update(nil, 2, &models.Job{
		Name: "ddd",
	})
	if assert.Nil(t, err) && assert.NotNil(t, job) {
		assert.Equal(t, 2, job.ID)
		assert.Equal(t, "ddd", job.Name)
	}

	// dao error
	_, err = s.Update(nil, 100, &models.Job{
		Name: "ddd",
	})
	assert.NotNil(t, err)

	// validation error
	_, err = s.Update(nil, 2, &models.Job{
		Name: "",
	})
	assert.NotNil(t, err)
}

func TestJobService_Delete(t *testing.T) {
	s := NewJobService(newMockJobDAO())
	job, err := s.Delete(nil, 2)
	if assert.Nil(t, err) && assert.NotNil(t, job) {
		assert.Equal(t, 2, job.ID)
		assert.Equal(t, "bbb", job.Name)
	}

	_, err = s.Delete(nil, 2)
	assert.NotNil(t, err)
}

func TestJobService_Query(t *testing.T) {
	s := NewJobService(newMockJobDAO())
	result, err := s.Query(nil, 1, 2)
	if assert.Nil(t, err) {
		assert.Equal(t, 2, len(result))
	}
}

func newMockJobDAO() jobDAO {
	return &mockJobDAO{
		records: []models.Job{
			{ID: 1, Name: "aaa"},
			{ID: 2, Name: "bbb"},
			{ID: 3, Name: "ccc"},
		},
	}
}

type mockJobDAO struct {
	records []models.Job
}

func (m *mockJobDAO) Get(rs app.RequestScope, id int) (*models.Job, error) {
	for _, record := range m.records {
		if record.ID == id {
			return &record, nil
		}
	}
	return nil, errors.New("not found")
}

func (m *mockJobDAO) Query(rs app.RequestScope, offset, limit int) ([]models.Job, error) {
	return m.records[offset : offset+limit], nil
}

func (m *mockJobDAO) Count(rs app.RequestScope) (int, error) {
	return len(m.records), nil
}

func (m *mockJobDAO) Create(rs app.RequestScope, job *models.Job) error {
	if job.ID != 0 {
		return errors.New("Id cannot be set")
	}
	job.ID = len(m.records) + 1
	m.records = append(m.records, *job)
	return nil
}

func (m *mockJobDAO) Update(rs app.RequestScope, id int, job *models.Job) error {
	job.ID = id
	for i, record := range m.records {
		if record.ID == id {
			m.records[i] = *job
			return nil
		}
	}
	return errors.New("not found")
}

func (m *mockJobDAO) Delete(rs app.RequestScope, id int) error {
	for i, record := range m.records {
		if record.ID == id {
			m.records = append(m.records[:i], m.records[i+1:]...)
			return nil
		}
	}
	return errors.New("not found")
}
