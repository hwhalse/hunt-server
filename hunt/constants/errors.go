package constants

import "errors"

var (
	ErrNoDocsDeleted   = errors.New("no documents deleted")
	ErrInvalidCallsign = errors.New("invalid callsign")
	ErrUnknownEvent    = errors.New("unknown event")
	ErrNoDocsInserted  = errors.New("no documents inserted")
)
