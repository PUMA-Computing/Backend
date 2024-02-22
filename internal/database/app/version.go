package app

import (
	"Backend/internal/database"
	"Backend/internal/models"
	"context"
)

type Version struct {
}

func LoadVersion() (*models.Version, error) {
	var versionInfo models.Version
	err := database.DB.QueryRow(context.Background(), "SELECT latest_version, changelog FROM version_info ORDER BY id DESC LIMIT 1").
		Scan(&versionInfo.LatestVersion, &versionInfo.Changelog)
	if err != nil {
		return nil, err
	}

	return &versionInfo, nil
}

func UpdateVersion(version *models.Version) error {
	_, err := database.DB.Exec(context.Background(), "INSERT INTO version_info (latest_version, changelog) VALUES ($1, $2)", version.LatestVersion, version.Changelog)
	if err != nil {
		return err
	}

	return nil
}
