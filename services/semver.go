package services

import (
	"github.com/aufaitio/listener/models"
	"github.com/blang/semver"
)

// FilterByVersion filters out repositories that are out of range of dependency update.
func FilterByVersion(repList []*models.Repository, hook *models.NpmHook) []*models.Repository {
	var filteredList []*models.Repository

	for _, rep := range repList {
		var desiredVersion semver.Version
		pub, err := semver.Make(hook.Version)

		for _, dep := range rep.Dependencies {
			if dep.Name == hook.Name {
				// Simply grab the first, there should never be two but if there is...
				desiredVersion, err = semver.Make(dep.Semver)

				// Ensure version is within the range of package version. Perhaps we will have a config to override
				// that and always update or at least try...
				if err == nil && pub.LTE(desiredVersion) {
					filteredList = append(filteredList, rep)
				}
				break
			}
		}
	}

	return filteredList
}
