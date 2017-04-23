package session

import (
	"time"
)

const (
	tableName  = "sessions"
	timeFormat = "Mon Jan 2 15:04:05 -0700 MST 2006"
)

// Session defines a struct which holds the the details of a giving user session.
type Session struct {
	UserID   string    `json:"user_id"`
	PublicID string    `json:"public_id"`
	Token    string    `json:"token"`
	Expires  time.Time `json:"expires"`
}

// Table returns the given table which the given struct corresponds to.
func (u Session) Table() string {
	return tableName
}

// WithFields attempts to syncing the giving data within the provided
// map into it's own fields.
func (u *Session) WithFields(fields map[string]interface{}) error {
	if user, ok := fields["user_id"].(string); ok {
		u.UserID = user
	}

	if public, ok := fields["public_id"].(string); ok {
		u.PublicID = public
	}

	if token, ok := fields["token"].(string); ok {
		u.Token = token
	}

	if expires, ok := fields["expires"]; ok {
		switch co := expires.(type) {
		case string:
			t, err := time.Parse(timeFormat, co)
			if err != nil {
				return err
			}

			u.Expires = t.UTC()
		case time.Time:
			u.Expires = co.UTC()
		}
	}

	return nil
}

// Fields returns a map representing the data of the session.
func (u *Session) Fields() map[string]interface{} {
	return map[string]interface{}{
		"user_id":   u.UserID,
		"expires":   u.Expires,
		"token":     u.Token,
		"public_id": u.PublicID,
	}
}
