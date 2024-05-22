package services

import (
	"Backend/internal/database/app"
	"Backend/internal/models"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type VersionService struct {
	client *http.Client
	token  string
}

func NewVersionService(token string) *VersionService {
	return &VersionService{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
		token: token,
	}
}

func (vs *VersionService) GetVersion() (string, error) {
	version, err := app.GetVersion()
	if err != nil {
		return "", err
	}

	return version, nil
}

func (vs *VersionService) GetChangelog() ([]models.ChangelogEntry, error) {
	changelog, err := app.GetChangeLog()
	if err != nil {
		return nil, err
	}

	return changelog, nil
}

func (vs *VersionService) FetchVersion() (*models.GithubVersion, error) {
	url := "https://api.github.com/repos/PUFA-Computing/Frontend/releases/latest"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("Authorization", "Bearer "+vs.token)
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")

	resp, err := vs.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch latest version: %w", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("Failed to close response body")
		}
	}(resp.Body)

	var release struct {
		TagName     string `json:"tag_name"`
		Body        string `json:"body"`
		PublishedAt string `json:"published_at"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	version := &models.GithubVersion{
		TagName:     release.TagName,
		Body:        release.Body,
		PublishedAt: release.PublishedAt,
	}

	return version, nil
}

func (vs *VersionService) CheckVersion(version string) (bool, error) {
	latestVersion, err := app.CheckVersion(version)
	if err != nil {
		return false, err
	}

	return latestVersion, nil
}

func (vs *VersionService) UpdateLatestVersion(version string) error {
	err := app.UpdateLatestVersion(version)
	if err != nil {
		return err
	}

	return nil
}

func (vs *VersionService) UpdateChangelog(version string, changelog string) error {
	err := app.UpdateChangelog(version, changelog)
	if err != nil {
		return err
	}

	return nil
}
