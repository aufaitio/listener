package services

import (
	"errors"
	"github.com/aufaitio/data-access"
	"github.com/aufaitio/data-access/models"
	"github.com/aufaitio/listener/app"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

type MockRequestScope struct {
	mock.Mock
	app.RequestScope
}

func (m MockRequestScope) DB() *mongo.Database {
	return &mongo.Database{}
}

func TestNewJobService(t *testing.T) {
	dao := newMockJobDAO()
	s := NewJobService(dao, newMockRepositoryDAO())
	assert.Equal(t, dao, s.dao)
}

func TestJobService_Get(t *testing.T) {
	s := NewJobService(newMockJobDAO(), newMockRepositoryDAO())
	job, err := s.Get(new(MockRequestScope), 1)
	if assert.Nil(t, err) && assert.NotNil(t, job) {
		assert.Equal(t, "aaa", job.Name)
	}

	job, err = s.Get(new(MockRequestScope), 100)
	assert.NotNil(t, err)
}

func TestJobService_Create(t *testing.T) {
	s := NewJobService(newMockJobDAO(), newMockRepositoryDAO())
	job, err := s.Create(new(MockRequestScope), createJob("ddd", "testing", "1.1.1"))
	if assert.Nil(t, err) && assert.NotNil(t, job) {
		assert.Equal(t, int64(4), job.ID)
		assert.Equal(t, "ddd", job.Name)
	}

	// dao error
	_, err = s.Create(new(MockRequestScope), &models.Job{
		ID:   100,
		Name: "ddd",
	})
	assert.NotNil(t, err)

	// validation error
	_, err = s.Create(new(MockRequestScope), &models.Job{
		Name: "",
	})
	assert.NotNil(t, err)
}

func TestJobService_Update(t *testing.T) {
	s := NewJobService(newMockJobDAO(), newMockRepositoryDAO())
	job, err := s.Update(new(MockRequestScope), 2, createJob("ddd", "a", "1.2.4"))
	if assert.Nil(t, err) && assert.NotNil(t, job) {
		assert.Equal(t, int64(2), job.ID)
		assert.Equal(t, "ddd", job.Name)
	}

	// dao error
	_, err = s.Update(new(MockRequestScope), 100, &models.Job{
		Name: "ddd",
	})
	assert.NotNil(t, err)

	// validation error
	_, err = s.Update(new(MockRequestScope), 2, &models.Job{
		Name: "",
	})
	assert.NotNil(t, err)
}

func TestJobService_Delete(t *testing.T) {
	s := NewJobService(newMockJobDAO(), newMockRepositoryDAO())
	job, err := s.Delete(new(MockRequestScope), 2)
	if assert.Nil(t, err) && assert.NotNil(t, job) {
		assert.Equal(t, int64(2), job.ID)
		assert.Equal(t, "bbb", job.Name)
	}

	_, err = s.Delete(new(MockRequestScope), 2)
	assert.NotNil(t, err)
}

func TestJobService_Query(t *testing.T) {
	s := NewJobService(newMockJobDAO(), newMockRepositoryDAO())
	result, err := s.Query(new(MockRequestScope), 1, 2)
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

func newMockJobDAO() access.JobDAO {
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

func (m *mockJobDAO) Get(db *mongo.Database, id int64) (*models.Job, error) {
	for _, record := range m.records {
		if record.ID == id {
			return record, nil
		}
	}
	return nil, errors.New("not found")
}

func (m *mockJobDAO) GetByName(db *mongo.Database, name string) (*models.Job, error) {
	for _, record := range m.records {
		if record.Name == name {
			return record, nil
		}
	}
	return nil, errors.New("not found")
}

func (m *mockJobDAO) Query(db *mongo.Database, offset, limit int) ([]*models.Job, error) {
	return m.records[offset : offset+limit], nil
}

func (m *mockJobDAO) Count(db *mongo.Database) (int64, error) {
	return int64(len(m.records)), nil
}

func (m *mockJobDAO) Create(db *mongo.Database, job *models.Job) error {
	if job.ID != 0 {
		return errors.New("Id cannot be set")
	}
	job.ID = int64(len(m.records) + 1)

	m.records = append(m.records, job)
	return nil
}

func (m *mockJobDAO) Update(db *mongo.Database, id int64, job *models.Job) error {
	job.ID = id
	for i, record := range m.records {
		if record.ID == id {
			m.records[i] = job
			return nil
		}
	}
	return errors.New("not found")
}

func (m *mockJobDAO) Delete(db *mongo.Database, id int64) error {
	for i, record := range m.records {
		if record.ID == id {
			m.records = append(m.records[:i], m.records[i+1:]...)
			return nil
		}
	}
	return errors.New("not found")
}
