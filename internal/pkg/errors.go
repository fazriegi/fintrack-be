package pkg

import "errors"

var (
	ErrServer          = errors.New("unexpected error occured")
	ErrParseReqBody    = errors.New("error parsing request body")
	ErrParseQueryParam = errors.New("error parsing query param")
	ErrValidation      = errors.New("validation error")
	ErrNotAuthorized   = errors.New("you're not authorized")
)
