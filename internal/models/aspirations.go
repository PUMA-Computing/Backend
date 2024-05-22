package models

import (
	"github.com/google/uuid"
	"time"
)

type Aspiration struct {
	ID             int                    `json:"id"`
	UserID         uuid.UUID              `json:"user_id"`
	Subject        string                 `json:"subject"`
	Message        string                 `json:"message"`
	Anonymous      bool                   `json:"anonymous"`
	OrganizationID int                    `json:"organization_id"`
	Closed         bool                   `json:"closed"`
	CreatedAt      time.Time              `json:"created_at"`
	UpdatedAt      time.Time              `json:"updated_at"`
	Upvote         int                    `json:"upvote"`
	AdminReply     *string                `json:"admin_reply"`
	Organization   AspirationOrganization `json:"organization"`
	Author         AspirationAuthor       `json:"author"`
}

type AspirationAuthor struct {
	Name           string `json:"name"`
	ProfilePicture string `json:"profile_picture"`
}

type AspirationOrganization struct {
	Name string `json:"name"`
}
