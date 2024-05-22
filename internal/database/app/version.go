package app

import (
	"Backend/internal/database"
	"Backend/internal/models"
	"context"
	"encoding/json"
)

type Version struct {
}

// GetVersion retrieves the current version from the database
func GetVersion() (string, error) {
	var version string
	err := database.DB.QueryRow(context.Background(), "SELECT latest_version FROM version_info").Scan(&version)
	if err != nil {
		return "", err
	}

	return version, nil
}

// CheckVersion checks if the current version is the latest version
func CheckVersion(version string) (bool, error) {
	var latestVersion string
	err := database.DB.QueryRow(context.Background(), "SELECT latest_version FROM version_info").Scan(&latestVersion)
	if err != nil {
		return false, err
	}

	if latestVersion == version {
		return true, nil
	}

	return false, nil
}

func GetChangeLog() ([]models.ChangelogEntry, error) {
	var changelogJSON []byte
	err := database.DB.QueryRow(context.Background(), "SELECT changelog FROM version_info").Scan(&changelogJSON)
	if err != nil {
		return nil, err
	}

	var changelog []models.ChangelogEntry
	err = json.Unmarshal(changelogJSON, &changelog)
	if err != nil {
		return nil, err
	}

	return changelog, nil
}

func UpdateLatestVersion(version string) error {
	_, err := database.DB.Exec(context.Background(), "UPDATE version_info SET latest_version = $1", version)
	if err != nil {
		return err
	}
	return nil
}

func UpdateChangelog(version string, changelog string) error {
	var existingChangelog []models.ChangelogEntry
	err := database.DB.QueryRow(context.Background(), "SELECT changelog FROM version_info").Scan(&existingChangelog)
	if err != nil {
		return err
	}

	// Append new changelog entry
	newEntry := models.ChangelogEntry{version: []string{changelog}}
	existingChangelog = append(existingChangelog, newEntry)

	// Convert back to JSON
	updatedChangelog, err := json.Marshal(existingChangelog)
	if err != nil {
		return err
	}

	// Update the database
	_, err = database.DB.Exec(context.Background(), "UPDATE version_info SET changelog = $1", updatedChangelog)
	if err != nil {
		return err
	}

	return nil
}
