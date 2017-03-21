package model

import "database/sql"

// User - user model
type User struct {
	ID              string         `json:"id,omitempty"`
	Email           string         `json:"email,omitempty"`
	FirstName       string         `json:"firstName,omitempty"`
	LastName        string         `json:"lastName,omitempty"`
	Password        string         `json:"-"`
	FcmToken        sql.NullString `json:"-"`
	FacebookID      sql.NullString `json:"-"`
	Jwt             string         `json:"jwt,omitempty"`
	LastRefreshTime int64          `json:"lastRefreshTime,omitempty"`
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
	Locked     bool           `json:"locked,omitempty"` // if the group is password protected
	Public     bool           `json:"public,omitempty"`
}

// Message - message model
type Message struct {
	ID        string `json:"id,omitempty"`
	UserID    string `json:"userID,omitempty"`
	GroupID   string `json:"groupID"`
	FirstName string `json:"firstName,omitempty"`
	LastName  string `json:"lastName,omitempty"`
	Content   string `json:"content,omitempty"`
	Timestamp string `json:"timestamp,omitempty"`
}

// GeoLocation - geo location model
type GeoLocation struct {
	ID        string `json:"id,omitempty"`
	UserID    string `json:"userID,omitempty"`
	GroupID   string `json:"groupID"`
	FirstName string `json:"firstName,omitempty"`
	LastName  string `json:"lastName,omitempty"`
	Email     string `json:"email,omitempty"`
	Latitude  string `json:"latitude,omitempty"`
	Longitude string `json:"longitude,omitempty"`
	Timestamp string `json:"timestamp,omitempty"`
}
