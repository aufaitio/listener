package models

import (
	"github.com/go-ozzo/ozzo-validation"
	"github.com/mongodb/mongo-go-driver/bson"
)

// PublishedDependency represents a packages dependencies
type PublishedDependency struct {
	Name    string `json:"name" bson:"name"`
	Version string `json:"version" bson:"version"`
}

// Job represents Builder Job.
type Job struct {
	Dependencies []PublishedDependency `json:"dependencies" bson:"dependencies"`
	Expiration   string                `json:"expiration" bson:"expiration"`
	ID           int64                 `json:"id" bson:"_id"`
	Name         string                `json:"name" bson:"name"`
	State        string                `json:"state" bson:"state"`
}

// Validate validates the repository fields.
func (job Job) Validate() error {
	return validation.ValidateStruct(&job,
		validation.Field(&job.Name, validation.Required),
		validation.Field(&job.State, validation.Required),
		validation.Field(&job.Dependencies, validation.Required),
	)
}

// NewJobFromDoc creates a new Job instance from DB doc.
// I need to look into better ways to do the decoding and encoding. Probably could use the Decoder interface.
func NewJobFromDoc(doc *bson.Document) (*Job, error) {
	var job Job
	keys, err := doc.Keys(false)

	if err != nil {
		return &job, err
	}

	job = Job{}

	for _, key := range keys {
		keyString := key.String()
		elm := doc.Lookup(keyString)
		if err != nil {
			return &job, err
		}

		// I need to find a better way to marshal these.
		switch keyString {
		case "dependencies":
			rawDepList := elm.MutableArray()
			depList := make([]PublishedDependency, rawDepList.Len())

			for i := uint(0); i < uint(rawDepList.Len()); i++ {
				elm, err := rawDepList.Lookup(i)

				if err != nil {
					return &job, err
				}

				doc := elm.MutableDocument()
				name := doc.Lookup("name")
				version := doc.Lookup("version")

				depList = append(
					depList,
					PublishedDependency{Name: name.StringValue(), Version: version.StringValue()},
				)
			}

			job.Dependencies = depList
		case "expiration":
			job.Expiration = elm.StringValue()
		case "name":
			job.Name = elm.StringValue()
		case "state":
			job.State = elm.StringValue()
		case "_id":
			job.ID = elm.Int64()
		default:
		}
	}

	return &job, err
}

// NewDocFromJob create bson Document from Job
func NewDocFromJob(job *Job) *bson.Document {
	var depList *bson.Array

	for _, dep := range job.Dependencies {
		depBson := bson.VC.Document(bson.NewDocument(
			bson.EC.String("name", dep.Name),
			bson.EC.String("version", dep.Version),
		))
		depList = depList.Append(depBson)
	}

	return bson.NewDocument(
		bson.EC.String("name", job.Name),
		bson.EC.Array("dependencies", depList),
		bson.EC.String("state", job.State),
	)
}
