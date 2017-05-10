package model

import "database/sql"

// redis keys
const (
	RateLimitKey = "rl:"
)

// UserSettings - model
type UserSettings struct {
	UserID        string `json:"userID"`
	Notifications bool   `json:"notifications"`
	FcmToken      string `json:"-"`
}

// User - user model
type User struct {
	ID              int64          `json:"id,omitempty"`
	Email           string         `json:"email,omitempty"`
	FirstName       string         `json:"firstName,omitempty"`
	LastName        string         `json:"lastName,omitempty"`
	Password        string         `json:"-"`
	FcmToken        sql.NullString `json:"-"`
	FacebookID      sql.NullString `json:"-"`
	FacebookUser    bool           `json:"facebookUser"`
	Jwt             string         `json:"jwt,omitempty"`
	LastRefreshTime int64          `json:"lastRefreshTime,omitempty"`
}

// Group - group model
type Group struct {
	ID         int64          `json:"id,omitempty"`
	Name       string         `json:"name,omitempty"`
	OwnerID    int64          `json:"ownerID,omitempty"`
	OwnerEmail string         `json:"ownerEmail,omitempty"`
	OwnerName  string         `json:"ownerName,omitempty"`
	UserCount  int            `json:"userCount,omitempty"`
	Password   sql.NullString `json:"-"`
	Locked     bool           `json:"locked,omitempty"` // if the group is password protected
	Public     bool           `json:"public,omitempty"`
}

// Message - message model
type Message struct {
	ID        int64  `json:"id,omitempty"`
	UserID    int64  `json:"userID,omitempty"`
	GroupID   int64  `json:"groupID,omitempty"`
	FirstName string `json:"firstName,omitempty"`
	LastName  string `json:"lastName,omitempty"`
	Content   string `json:"content,omitempty"`
	Timestamp int64  `json:"timestamp,omitempty"`
}

// GeoLocation - geo location model
type GeoLocation struct {
	ID        int64  `json:"id,omitempty"`
	UserID    int64  `json:"userID,omitempty"`
	GroupID   int64  `json:"groupID"`
	FirstName string `json:"firstName,omitempty"`
	LastName  string `json:"lastName,omitempty"`
	Email     string `json:"email,omitempty"`
	Latitude  string `json:"latitude,omitempty"`
	Longitude string `json:"longitude,omitempty"`
	Timestamp int64  `json:"timestamp,omitempty"`
	Waypoint  bool   `json:"waypoint,omitempty"`
}
