package auth

import "errors"

var (
	ErrNoAuthHeaderIncluded = errors.New("no authorization header included")
	ErrMalformedAuthHeader  = errors.New("malformed authorization header")
)
