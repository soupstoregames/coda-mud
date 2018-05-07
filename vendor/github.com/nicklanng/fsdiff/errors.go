package fsdiff

import "errors"

var (
	ErrPathNotFound          = errors.New("path not found")
	ErrFailedToReadDirectory = errors.New("failed to read directory")
	ErrFailedToReadFile      = errors.New("failed to read file")
	ErrFailedToComputeHash   = errors.New("failed to compute hash")
)
