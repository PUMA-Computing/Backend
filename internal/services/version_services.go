package services

import (
	"Backend/internal/database/app"
	"Backend/internal/models"
)

type VersionService struct {
}

func NewVersionService() *VersionService {
	return &VersionService{}
}

func (vs *VersionService) GetVersion() (*models.Version, error) {
	version, err := app.LoadVersion()
	if err != nil {
		return nil, err
	}
	return version, nil
}

//func (vs *VersionService) UpdateVersion(version *models.Version) error {
//	if err := app.UpdateVersion(version); err != nil {
//		return err
//	}
//	return nil
//}
