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
	s := NewJobService(dao, newMockRepositoryDAO())
	assert.Equal(t, dao, s.dao)
}

func TestJobService_Get(t *testing.T) {
	s := NewJobService(newMockJobDAO(), newMockRepositoryDAO())
	job, err := s.Get(nil, 1)
	if assert.Nil(t, err) && assert.NotNil(t, job) {
		assert.Equal(t, "aaa", job.Name)
	}

	job, err = s.Get(nil, 100)
	assert.NotNil(t, err)
}

func TestJobService_Create(t *testing.T) {
	s := NewJobService(newMockJobDAO(), newMockRepositoryDAO())
	job, err := s.Create(nil, createJob("ddd", "testing", "1.1.1"))
	if assert.Nil(t, err) && assert.NotNil(t, job) {
		assert.Equal(t, int64(4), job.ID)
		assert.Equal(t, "ddd", job.Name)
	}

	// dao error
	_, err = s.Create(nil, &models.Job{
		ID:   100,
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
	s := NewJobService(newMockJobDAO(), newMockRepositoryDAO())
	job, err := s.Update(nil, 2, createJob("ddd", "a", "1.2.4"))
	if assert.Nil(t, err) && assert.NotNil(t, job) {
		assert.Equal(t, int64(2), job.ID)
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
	s := NewJobService(newMockJobDAO(), newMockRepositoryDAO())
	job, err := s.Delete(nil, 2)
	if assert.Nil(t, err) && assert.NotNil(t, job) {
		assert.Equal(t, int64(2), job.ID)
		assert.Equal(t, "bbb", job.Name)
	}

	_, err = s.Delete(nil, 2)
	assert.NotNil(t, err)
}

func TestJobService_Query(t *testing.T) {
	s := NewJobService(newMockJobDAO(), newMockRepositoryDAO())
	result, err := s.Query(nil, 1, 2)
	if assert.Nil(t, err) {
		assert.Equal(t, 2, len(result))
	}
}

func createJob(name string, depName string, depVersion string) *models.Job {
	return &models.Job{
		Name:  name,
		State: models.Idle,
		Dependencies: []*models.PublishedDependency{
			&models.PublishedDependency{Name: depName, Version: depVersion},
		},
	}
}

func newMockJobDAO() jobDAO {
	jobList := []*models.Job{
		createJob("aaa", "test", "1.2.3"),
		createJob("bbb", "test", "2.2.3"),
		createJob("ccc", "test", "3.2.3"),
	}

	for i, job := range jobList {
		job.ID = int64(i + 1)
	}

	return &mockJobDAO{records: jobList}
}

type mockJobDAO struct {
	records []*models.Job
}

func (m *mockJobDAO) Get(rs app.RequestScope, id int64) (*models.Job, error) {
	for _, record := range m.records {
		if record.ID == id {
			return record, nil
		}
	}
	return nil, errors.New("not found")
}

func (m *mockJobDAO) GetByName(rs app.RequestScope, name string) (*models.Job, error) {
	for _, record := range m.records {
		if record.Name == name {
			return record, nil
		}
	}
	return nil, errors.New("not found")
}

func (m *mockJobDAO) Query(rs app.RequestScope, offset, limit int) ([]*models.Job, error) {
	return m.records[offset : offset+limit], nil
}

func (m *mockJobDAO) Count(rs app.RequestScope) (int64, error) {
	return int64(len(m.records)), nil
}

func (m *mockJobDAO) Create(rs app.RequestScope, job *models.Job) error {
	if job.ID != 0 {
		return errors.New("Id cannot be set")
	}
	job.ID = int64(len(m.records) + 1)

	m.records = append(m.records, job)
	return nil
}

func (m *mockJobDAO) Update(rs app.RequestScope, id int64, job *models.Job) error {
	job.ID = id
	for i, record := range m.records {
		if record.ID == id {
			m.records[i] = job
			return nil
		}
	}
	return errors.New("not found")
}

func (m *mockJobDAO) Delete(rs app.RequestScope, id int64) error {
	for i, record := range m.records {
		if record.ID == id {
			m.records = append(m.records[:i], m.records[i+1:]...)
			return nil
		}
	}
	return errors.New("not found")
}
