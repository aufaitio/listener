package models

import (
	"github.com/go-ozzo/ozzo-validation"
)

type (
	// NpmHook structure for open source NPM version:publish hook
	NpmHook struct {
		Name    string `json:"name"`
		Source  string `json:source`
		Version string `json:"version"`
	}
)

// Validate validates NpmHook structure
func (m NpmHook) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.Name, validation.Required),
		validation.Field(&m.Version, validation.Required),
	)
}
