package models

import (
	"github.com/go-ozzo/ozzo-validation"
	"github.com/mongodb/mongo-go-driver/bson"
)

// Config a managed repositories config
type Config struct {
	Branch string `json:"branch" bson:"branch"`
	Remote string `json:"remote" bson:"remote"`
}

// Dependency represents a packages dependencies
type Dependency struct {
	Installed string `json:"installed" bson:"installed"`
	Name      string `json:"name" bson:"name"`
	Semver    string `json:"semver" bson:"semver"`
}

// Repository represents a Repository registered with Au Fait.
type Repository struct {
	ID           int64        `json:"id" bson:"id"`
	Name         string       `json:"name" bson:"name"`
	Config       Config       `json:"config" bson:"config"`
	Dependencies []Dependency `json:"dependencies" bson:"dependencies"`
}

// Validate validates the repository fields.
func (m Repository) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.Name, validation.Required),
		validation.Field(&m.Config, validation.Required),
		validation.Field(&m.Dependencies, validation.Required),
	)
}

// NewRepositoryFromDoc creates a new Repository instance from DB doc.
// I need to look into better ways to do the decoding and encoding. Probably could use the Decoder interface.
func NewRepositoryFromDoc(doc *bson.Document) (*Repository, error) {
	var repository Repository
	keys, err := doc.Keys(false)

	if err != nil {
		return &repository, err
	}

	repository = Repository{}

	for _, key := range keys {
		keyString := key.String()
		elm := doc.Lookup(keyString)
		if err != nil {
			return &repository, err
		}

		// I need to find a better way to marshal these.
		switch keyString {
		case "dependencies":
			rawDepList := elm.MutableArray()
			depList := make([]Dependency, rawDepList.Len())

			for i := uint(0); i < uint(rawDepList.Len()); i++ {
				elm, err := rawDepList.Lookup(i)

				if err != nil {
					return &repository, err
				}

				doc := elm.MutableDocument()
				installed := doc.Lookup("installed")
				name := doc.Lookup("name")
				semver := doc.Lookup("semver")

				depList = append(
					depList,
					Dependency{
						Name:      name.StringValue(),
						Semver:    semver.StringValue(),
						Installed: installed.StringValue(),
					},
				)
			}

			repository.Dependencies = depList
		case "name":
			repository.Name = elm.StringValue()
		case "config":
			configDoc := elm.MutableDocument()
			branch := configDoc.Lookup("branch")
			remote := configDoc.Lookup("remote")

			repository.Config = Config{Branch: branch.StringValue(), Remote: remote.StringValue()}
		case "_id":
			repository.ID = elm.Int64()
		default:
		}
	}

	return &repository, err
}

// NewDocFromRepository create bson Document from Repository
func NewDocFromRepository(repository *Repository) *bson.Document {
	var depList *bson.Array

	for _, dep := range repository.Dependencies {
		depBson := bson.VC.Document(bson.NewDocument(
			bson.EC.String("name", dep.Name),
			bson.EC.String("installed", dep.Installed),
			bson.EC.String("semver", dep.Semver),
		))
		depList = depList.Append(depBson)
	}

	return bson.NewDocument(
		bson.EC.String("name", repository.Name),
		bson.EC.Array("dependencies", depList),
		bson.EC.SubDocument("config", bson.NewDocument(
			bson.EC.String("remote", repository.Config.Remote),
			bson.EC.String("branch", repository.Config.Branch),
		)),
	)
}
