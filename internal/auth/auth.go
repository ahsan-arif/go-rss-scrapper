package auth

import (
	"errors"
	"net/http"
	"strings"
)

// GetAPIKey extracts API Key from the Header of an HTTP request
// Example:
// Authorization: APIKey {insert apikey here}
func GetAPIKey(header http.Header) (string, error) {
	val := header.Get("Authorization")

	if val == "" {
		return "", errors.New("no authentication info found")
	}
	vals := strings.Split(val, " ")

	if len(vals) != 2 {
		return "", errors.New("malformed auth header")
	}
	if vals[0] != "APIKey" {
		return "", errors.New("malformed first part of auth header")
	}
	return vals[1], nil
}
