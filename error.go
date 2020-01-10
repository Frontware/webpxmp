package xmp

import "errors"

var (
	// ErrInvalidRIFF means that the FourCC "RIFF" is not found from byte 0 to 3.
	ErrInvalidRIFF = errors.New("invalid RIFF container")

	// ErrInvalidWEBP means that the FourCC "WEBP" is not found from byte 8 to 11.
	ErrInvalidWEBP = errors.New("invalid WEBP format")

	// ErrVP8XNotFound means that the FourCC "VP8X" is not found from byte 12 to 15.
	ErrVP8XNotFound = errors.New("VP8 extensible not found")

	// ErrXMPNotFound means either the flag "XMP" is false in VP8X metadata or XMP chunk is
	// nowhere to be found.
	ErrXMPNotFound = errors.New("XMP not found")
)
