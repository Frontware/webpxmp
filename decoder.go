package xmp

import (
	"encoding/binary"
	"encoding/xml"
	"io/ioutil"
	"net/http"
)

// Read parses XML-encoded data inside WEBP RIFF container.
func Read(path string, v interface{}) error {
	b, err := ioutil.ReadFile(path)
	if err == nil {
		return decode(b, v)
	}

	resp, err := http.Get(path)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	b, err = ioutil.ReadAll(resp.Body)
	if err == nil {
		return decode(b, v)
	}

	return err
}

// ReadXMP returns xml information on profile.
func ReadXMP(path string) (*Profile, error) {
	xmp := xmpMeta{}
	err := Read(path, &xmp)
	if err != nil {
		return nil, err
	}

	profile := xmp.RDF.Description.Profile

	return &profile, nil
}

// Refer to https://developers.google.com/speed/webp/docs/riff_container
// decode reads XML inside WEBP RIFF container.
func decode(b []byte, v interface{}) error {

	chunkOffset := 0
	chunkID := string(b[chunkOffset : chunkOffset+4])
	if chunkID != "RIFF" {
		return ErrInvalidRIFF
	}
	if string(b[chunkOffset+8:chunkOffset+12]) != "WEBP" {
		return ErrInvalidWEBP
	}
	containerSize := 12

	chunkOffset = chunkOffset + containerSize
	chunkID = string(b[chunkOffset : chunkOffset+4])
	if chunkID != "VP8X" {
		return ErrVP8XNotFound
	}
	containerSize = int(binary.LittleEndian.Uint32(b[chunkOffset+4:chunkOffset+8])) + 8
	metadata := b[chunkOffset+8]
	xmpFlag := metadata & 0x4
	if xmpFlag != 4 {
		return ErrXMPNotFound
	}

	chunkOffset = chunkOffset + containerSize
	chunkID = string(b[chunkOffset : chunkOffset+4])
	containerSize = int(binary.LittleEndian.Uint32(b[chunkOffset+4:chunkOffset+8])) + 8

	for chunkOffset < len(b) && chunkID != "XMP " {
		chunkOffset = chunkOffset + containerSize
		chunkID = string(b[chunkOffset : chunkOffset+4])
		containerSize = int(binary.LittleEndian.Uint32(b[chunkOffset+4:chunkOffset+8])) + 8
	}

	if chunkID != "XMP " {
		return ErrXMPNotFound
	}

	content := b[chunkOffset+8 : chunkOffset+containerSize+8]

	err := xml.Unmarshal(content, v)
	if err != nil {
		return err
	}

	return err
}
