package models

type Version struct {
	LatestVersion string           `json:"latest_version"`
	Changelog     []ChangelogEntry `json:"changelog"`
}

type ChangelogEntry map[string][]string

type GithubVersion struct {
	TagName     string `json:"tag_name"`
	Body        string `json:"body"`
	PublishedAt string `json:"published_at"`
}
