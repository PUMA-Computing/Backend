package services

import (
	"Backend/internal/database/app"
	"log"
	"time"
)

type VersionUpdater struct {
	VersionService *VersionService
}

func NewVersionUpdater(versionService *VersionService) *VersionUpdater {
	return &VersionUpdater{VersionService: versionService}
}

func (v *VersionUpdater) Run() {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			log.Println("VersionUpdater: Version update loop started")
			// Get latest version
			latestVersion, err := v.VersionService.FetchVersion()
			if err != nil {
				// Log the error
				log.Println("Error fetching latest version:", err)
				continue
			}

			// Check current version
			isLatest, err := app.CheckVersion(latestVersion.TagName)
			if err != nil {
				// Log the error
				log.Println("Error checking version:", err)
				continue
			}

			// Update version if not latest
			if !isLatest {
				err := app.UpdateLatestVersion(latestVersion.TagName)
				if err != nil {
					// Log the error
					log.Println("Error updating version:", err)
					continue
				}

				// Update changelog
				err = app.UpdateChangelog(latestVersion.TagName, latestVersion.Body)
				if err != nil {
					// Log the error
					log.Println("Error updating changelog:", err)
					continue
				}
			}

			log.Println("VersionUpdater: Version update loop completed")
		}
	}
}
