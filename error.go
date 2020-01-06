package xmp

import "errors"

var (
	InvalidRIFF  = errors.New("invalid RIFF container")
	InvalidWEBP  = errors.New("invalid WEBP format")
	VP8XNotFound = errors.New("VP8 extensible not found")
	XMPNotFound  = errors.New("XMP not found")
)
