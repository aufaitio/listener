package services

import (
	"github.com/aufaitio/data-access"
	"github.com/aufaitio/data-access/models"
	"github.com/aufaitio/listener/app"
)

// JobService provides services related with repositories.
type JobService struct {
	dao    access.JobDAO
	repDao access.RepositoryDAO
}

// NewJobService creates a new JobService with the given job DAO.
func NewJobService(dao access.JobDAO, repDao access.RepositoryDAO) *JobService {
	return &JobService{dao, repDao}
}

// Get returns the job with the specified the job ID.
func (s *JobService) Get(rs app.RequestScope, id int64) (*models.Job, error) {
	return s.dao.Get(rs.DB(), id)
}

// CreateJobsFromHook creates a list of jobs from a NPM Hook dependency
func (s *JobService) CreateJobsFromHook(rs app.RequestScope, hook *models.NpmHook) ([]*models.Job, error) {
	var jobList []*models.Job

	repList, err := s.repDao.QueryByDependency(rs.DB(), hook.Name)

	if err != nil {
		return jobList, err
	}

	filterRepList := FilterByVersion(repList, hook)

	for _, rep := range filterRepList {
		existingJob, err := s.dao.GetByName(rs.DB(), rep.Name)
		job := existingJob

		if err != nil {
			return jobList, err
		}

		publishedDep := models.PublishedDependency{Name: hook.Name, Version: hook.Version}

		if existingJob.Name != rep.Name || existingJob.State == models.InProgress {
			publishedDepList := []*models.PublishedDependency{&publishedDep}
			job = models.NewJobFromRepository(rep, publishedDepList)

			// Jobs in progress that get new dependencies, get a new job that is locked until it is complete.
			if existingJob.State == models.InProgress {
				existingJob.State = models.Locked
			}
		} else {
			job.Dependencies = append(job.Dependencies, &publishedDep)
		}

		jobList = append(jobList, job)
	}

	return jobList, nil
}

// Create creates a new job.
func (s *JobService) Create(rs app.RequestScope, model *models.Job) (*models.Job, error) {
	if err := model.Validate(); err != nil {
		return nil, err
	}
	if err := s.dao.Create(rs.DB(), model); err != nil {
		return nil, err
	}
	return s.dao.Get(rs.DB(), model.ID)
}

// Update updates the job with the specified ID.
func (s *JobService) Update(rs app.RequestScope, id int64, model *models.Job) (*models.Job, error) {
	if err := model.Validate(); err != nil {
		return nil, err
	}
	if err := s.dao.Update(rs.DB(), id, model); err != nil {
		return nil, err
	}
	return s.dao.Get(rs.DB(), id)
}

// Delete deletes the job with the specified ID.
func (s *JobService) Delete(rs app.RequestScope, id int64) (*models.Job, error) {
	job, err := s.dao.Get(rs.DB(), id)
	if err != nil {
		return nil, err
	}
	err = s.dao.Delete(rs.DB(), id)
	return job, err
}

// Count returns the number of repositories.
func (s *JobService) Count(rs app.RequestScope) (int64, error) {
	return s.dao.Count(rs.DB())
}

// Query returns the repositories with the specified offset and limit.
func (s *JobService) Query(rs app.RequestScope, offset, limit int) ([]*models.Job, error) {
	return s.dao.Query(rs.DB(), offset, limit)
}
