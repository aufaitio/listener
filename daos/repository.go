package daos

import (
	"context"

	"github.com/aufaitio/listener/app"
	"github.com/aufaitio/listener/models"
	"github.com/mongodb/mongo-go-driver/bson"
)

// RepositoryDAO persists repository data in database
type RepositoryDAO struct{}

// NewRepositoryDAO creates a new RepositoryDAO
func NewRepositoryDAO() *RepositoryDAO {
	return &RepositoryDAO{}
}

// Get reads the repository with the specified ID from the database.
func (dao *RepositoryDAO) Get(rs app.RequestScope, id int64) (*models.Repository, error) {
	var repository *models.Repository
	col := rs.DB().Collection("repository")

	err := col.FindOne(
		context.Background(),
		bson.NewDocument(
			bson.EC.Int64("id", id),
		),
	).Decode(repository)

	if err != nil {
		return repository, err
	}

	return repository, err
}

// Create saves a new repository record in the database.
// The Repository.ID field will be populated with an automatically generated ID upon successful saving.
func (dao *RepositoryDAO) Create(rs app.RequestScope, repository *models.Repository) error {
	col := rs.DB().Collection("repository")
	repoBson := models.NewDocFromRepository(repository)

	_, err := col.InsertOne(
		context.Background(),
		repoBson,
	)

	return err
}

// Update saves the changes to an repository in the database.
func (dao *RepositoryDAO) Update(rs app.RequestScope, id int64, repository *models.Repository) error {
	if _, err := dao.Get(rs, id); err != nil {
		return err
	}

	repoBson := models.NewDocFromRepository(repository)
	col := rs.DB().Collection("repository")

	_, err := col.UpdateOne(
		context.Background(),
		bson.NewDocument(
			bson.EC.Int64("id", id),
		),
		repoBson,
	)
	return err
}

// Delete deletes an repository with the specified ID from the database.
func (dao *RepositoryDAO) Delete(rs app.RequestScope, id int64) error {
	repository, err := dao.Get(rs, id)
	if err != nil {
		return err
	}

	col := rs.DB().Collection("repository")
	_, err = col.DeleteOne(
		context.Background(),
		bson.NewDocument(
			bson.EC.String("name", repository.Name),
		),
	)

	return err
}

// Count returns the number of the repository records in the database.
func (dao *RepositoryDAO) Count(rs app.RequestScope) (int64, error) {
	return rs.DB().Collection("repository").Count(
		context.Background(),
		bson.NewDocument(),
	)
}

// Query retrieves the repository records with the specified offset and limit from the database.
func (dao *RepositoryDAO) Query(rs app.RequestScope, offset, limit int) ([]*models.Repository, error) {
	repositoryList := []*models.Repository{}
	col := rs.DB().Collection("repository")
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
			return repositoryList, err
		}
		job, err := models.NewRepositoryFromDoc(elm)

		if err != nil {
			return repositoryList, err
		}

		repositoryList = append(repositoryList, job)
	}

	return repositoryList, err
}
