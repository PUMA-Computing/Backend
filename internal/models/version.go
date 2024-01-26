package models

type Version struct {
	LatestVersion string      `json:"latest_version"`
	Changelog     []Changelog `json:"changelog"`
}

type Changelog struct {
	Version string   `json:"version"`
	Changes []string `json:"changes"`
}
