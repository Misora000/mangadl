package types

import "errors"

// Err definds the common errors.
var (
	ErrInvalidURL             = errors.New("invalid URL")
	ErrInsuficientChapterInfo = errors.New("chapter info is insuficient")
)
