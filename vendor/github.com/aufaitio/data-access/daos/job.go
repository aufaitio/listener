package daos

import (
	"context"
	"github.com/aufaitio/data-access"
	"github.com/aufaitio/data-access/models"
	"github.com/mongodb/mongo-go-driver/bson"
)

// JobDAO persists job data in database
type JobDAO struct{}

// NewJobDAO creates a new JobDAO
func NewJobDAO() *JobDAO {
	return &JobDAO{}
}

// Get reads the job with the specified ID from the database.
func (dao *JobDAO) Get(rs access.Scope, id int64) (*models.Job, error) {
	return dao.get(
		rs,
		bson.NewDocument(
			bson.EC.Int64("id", id),
		),
	)
}

// GetByName reads the job with the specified name from the database.
func (dao *JobDAO) GetByName(rs access.Scope, name string) (*models.Job, error) {
	return dao.get(
		rs,
		bson.NewDocument(
			bson.EC.String("name", name),
		),
	)
}

func (dao *JobDAO) get(rs access.Scope, doc *bson.Document) (*models.Job, error) {
	var job *models.Job
	col := rs.DB().Collection("job")
	result := bson.NewDocument()

	err := col.FindOne(
		context.Background(),
		doc,
	).Decode(result)

	if err != nil {
		return job, err
	}

	job, err = models.NewJobFromDoc(result)

	return job, err
}

// Create saves a new job record in the database.
// The Job.ID field will be populated with an automatically generated ID upon successful saving.
func (dao *JobDAO) Create(rs access.Scope, job *models.Job) error {
	col := rs.DB().Collection("job")

	jobBson := models.NewDocFromJob(job)
	_, err := col.InsertOne(
		context.Background(),
		jobBson,
	)

	return err
}

// Update saves the changes to an job in the database.
func (dao *JobDAO) Update(rs access.Scope, id int64, job *models.Job) error {
	if _, err := dao.Get(rs, id); err != nil {
		return err
	}

	jobBson := models.NewDocFromJob(job)
	col := rs.DB().Collection("job")
	_, err := col.UpdateOne(
		context.Background(),
		bson.NewDocument(
			bson.EC.Int64("_id", job.ID),
		),
		jobBson,
	)

	return err
}

// Delete deletes an job with the specified ID from the database.
func (dao *JobDAO) Delete(rs access.Scope, id int64) error {
	_, err := dao.Get(rs, id)
	if err != nil {
		return err
	}

	col := rs.DB().Collection("job")
	_, err = col.DeleteOne(
		context.Background(),
		bson.NewDocument(
			bson.EC.Int64("id", id),
		),
	)

	return err
}

// Count returns the number of the job records in the database.
func (dao *JobDAO) Count(rs access.Scope) (int64, error) {
	col := rs.DB().Collection("job")

	return col.Count(
		context.Background(),
		bson.NewDocument(),
	)
}

// Query retrieves the job records with the specified offset and limit from the database.
func (dao *JobDAO) Query(rs access.Scope, offset, limit int) ([]*models.Job, error) {
	jobList := []*models.Job{}
	col := rs.DB().Collection("job")
	ctx := context.Background()

	cursor, err := col.Find(
		ctx,
		bson.NewDocument(),
	)
	defer cursor.Close(ctx)
	elm := bson.NewDocument()

	for cursor.Next(ctx) {
		elm.Reset()

		if err := cursor.Decode(elm); err != nil {
			return jobList, err
		}
		job, err := models.NewJobFromDoc(elm)

		if err != nil {
			return jobList, err
		}

		jobList = append(jobList, job)
	}

	return jobList, err
}
