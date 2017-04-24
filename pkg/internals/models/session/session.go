package session

import (
	"encoding/base64"
	"errors"
	"time"

	uuid "github.com/satori/go.uuid"
)

const (
	tableName = "sessions"

	// UniqueIndex defines the unique index name used by the models db for model query optimization.
	UniqueIndex = "user_id"

	// UniqueIndexField defines the unique index field used by the model in it's field.
	UniqueIndexField = "user_public_id"
)

// NewSession defines the set of data received to create a new user.
type NewSession struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Session defines a struct which holds the the details of a giving user session.
type Session struct {
	UserID   string    `json:"user_public_id"`
	PublicID string    `json:"public_id"`
	Token    string    `json:"token"`
	Expires  time.Time `json:"expires"`
}

// New returns a new instance of a session.
func New(userID string, expiration time.Time) *Session {
	return &Session{
		UserID:   userID,
		PublicID: uuid.NewV4().String(),
		Token:    uuid.NewV4().String(),
		Expires:  expiration,
	}
}

// ValidateToken validates the provide base64 encoded token, that it matches the
// expected token value with that of the session.
func (u Session) ValidateToken(token string) bool {
	decoded, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		return false
	}

	if string(decoded) != u.Token {
		return false
	}

	return true
}

// SessionToken returns the Session.Token has a base64 encoded string.
func (u Session) SessionToken() string {
	return base64.StdEncoding.EncodeToString([]byte(u.Token))
}

// Table returns the given table which the given struct corresponds to.
func (u Session) Table() string {
	return tableName
}

// Fields returns a map representing the data of the session.
func (u Session) Fields() map[string]interface{} {
	return map[string]interface{}{
		"user_public_id": u.UserID,
		"token":          u.Token,
		"public_id":      u.PublicID,
		"expires":        u.Expires.Format(time.RFC3339),
	}
}

// WithFields attempts to syncing the giving data within the provided
// map into it's own fields.
func (u *Session) WithFields(fields map[string]interface{}) error {
	if user, ok := fields["user_id"].(string); ok {
		u.UserID = user
	} else {
		return errors.New("Expected 'user_id' key")
	}

	if public, ok := fields["public_id"].(string); ok {
		u.PublicID = public
	} else {
		return errors.New("Expected 'public_id' key")
	}

	if token, ok := fields["token"].(string); ok {
		u.Token = token
	} else {
		return errors.New("Expected 'token' key")
	}

	if expires, ok := fields["expires"]; ok && expires != "" {
		switch co := expires.(type) {
		case string:
			t, err := time.Parse(time.RFC3339, co)
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
