package model

import "database/sql"

// User - user model
type User struct {
	ID             string         `json:"id,omitempty"`
	Email          string         `json:"email,omitempty"`
	FirstName      string         `json:"firstName,omitempty"`
	LastName       string         `json:"lastName,omitempty"`
	Password       string         `json:"-"`
	FcmToken       sql.NullString `json:"-"`
	FacebookID     sql.NullString `json:"-"`
	Jwt            string         `json:"jwt,omitempty"`
	LastJwtRefresh int64          `json:"lastJwtRefresh,omitempty"`
}

// Group - group model
type Group struct {
	ID         string         `json:"id,omitempty"`
	Name       string         `json:"name,omitempty"`
	OwnerID    string         `json:"ownerID,omitempty"`
	OwnerEmail string         `json:"ownerEmail,omitempty"`
	OwnerName  string         `json:"ownerName,omitempty"`
	UserCount  int            `json:"userCount,omitempty"`
	Password   sql.NullString `json:"-"`
	Locked     bool           `json:"locked"` // if the group is password protected
	Public     bool           `json:"public,omitempty"`
}

// Message - message model
type Message struct {
	ID        string `json:"id,omitempty"`
	UserID    string `json:"userID,omitempty"`
	GroupID   string `json:"groupID"`
	FirstName string `json:"firstName,omitempty"`
	LastName  string `json:"lastName,omitempty"`
	Content   string `json:"Content,omitempty"`
	Timestamp string `json:"timestamp"`
}
