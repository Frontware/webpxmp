package xmp

import (
	"encoding/binary"
	"encoding/xml"
	"io/ioutil"
	"net/http"
	"strconv"
)

func Read(path string) (*Profile, error) {
	var p *Profile
	xmp := Xmpmeta{}

	b, err := ioutil.ReadFile(path)
	if err == nil {
		err = decode(b, &xmp)
		p = &xmp.RDF.Description.Profile
		return p, err
	}

	resp, err := http.Get(path)
	if err != nil {
		return p, err
	}
	defer resp.Body.Close()

	b, err = ioutil.ReadAll(resp.Body)
	if err == nil {
		err = decode(b, &xmp)
		p = &xmp.RDF.Description.Profile
		return p, err
	}

	return p, err
}

// Refer to https://developers.google.com/speed/webp/docs/riff_container

func decode(b []byte, v *Xmpmeta) error {

	chunkOffset := 0
	chunkID := string(b[chunkOffset : chunkOffset+4])
	if chunkID != "RIFF" {
		return InvalidRIFF
	}
	if string(b[chunkOffset+8:chunkOffset+12]) != "WEBP" {
		return InvalidWEBP
	}
	containerSize := 12

	chunkOffset = chunkOffset + containerSize
	chunkID = string(b[chunkOffset : chunkOffset+4])
	if chunkID != "VP8X" {
		return VP8XNotFound
	}
	containerSize = int(binary.LittleEndian.Uint32(b[chunkOffset+4:chunkOffset+8])) + 8
	metadata := b[chunkOffset+8]
	xmpFlag := metadata & 0x4
	if xmpFlag != 4 {
		return XMPNotFound
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
		return XMPNotFound
	}

	content := b[chunkOffset+8 : chunkOffset+containerSize+8]

	err := xml.Unmarshal(content, v)
	if err != nil {
		return err
	}

	escapedName, err := strconv.Unquote(v.RDF.Description.Profile.Name)
	if err != nil {
		return err
	}
	v.RDF.Description.Profile.Name = escapedName

	escapedLocation, err := strconv.Unquote(v.RDF.Description.Profile.Location)
	if err != nil {
		return err
	}
	v.RDF.Description.Profile.Location = escapedLocation

	return err
}
