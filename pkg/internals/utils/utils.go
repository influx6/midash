package utils

import "fmt"

// ErrorMessage returns a string which contains a json value of a
// given error message to be delivered.
func ErrorMessage(status int, header string, err error) string {
	return fmt.Sprintf(`{
		"status": %d,
		"title": %+q,
		"message": %+q,
	}`, status, header, err)
}
