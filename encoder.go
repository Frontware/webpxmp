package xmp

import (
	"encoding/binary"
	"encoding/xml"
	"io"
	"io/ioutil"
	"os"
)

// Write writes XML-encoded data format to WEBP RIFF container.
func Write(fileName string, v interface{}) error {
	b, err := ioutil.ReadFile(fileName)
	if err != nil {
		return err
	}

	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	return encode(file, b, v)
}

// WriteXMP writes profile information in XMP format to WEBP RIFF container.
func WriteXMP(fileName, employeeName string, timeStamp int64, lat, long float64, location string) error {
	xmp := xmpMeta{}
	profile := Profile{
		Name:      employeeName,
		Timestamp: timeStamp,
		Location:  location,
		Lat:       lat,
		Long:      long,
	}
	xmp.RDF.Description.Profile = profile

	return Write(fileName, xmp)
}

// Refer to https://developers.google.com/speed/webp/docs/riff_container
// encode writes xml to WEBP RIFF container.
func encode(w io.Writer, b []byte, v interface{}) error {

	var buffer []byte

	content, err := xml.Marshal(v)
	if err != nil {
		return err
	}

	chunkOffset := 0
	chunkID := string(b[chunkOffset : chunkOffset+4])
	if chunkID != "RIFF" {
		return ErrInvalidRIFF
	}
	if string(b[chunkOffset+8:chunkOffset+12]) != "WEBP" {
		return ErrInvalidWEBP
	}

	size := 0
	chunkOffset = 8
	containerSize := 4

	for chunkOffset < len(b) && chunkID != "XMP " {
		size += containerSize
		chunkOffset += containerSize
		chunkID = string(b[chunkOffset : chunkOffset+4])
		containerSize = int(binary.LittleEndian.Uint32(b[chunkOffset+4:chunkOffset+8])) + 8
	}

	chunkID = string(b[12:16])

	size += 8 + len(content)

	if chunkID != "VP8X" {
		size += 18
		if size%2 == 0 {
			buffer = make([]byte, chunkOffset+len(content)+26)
		} else {
			buffer = make([]byte, chunkOffset+len(content)+27)
		}

		copy(buffer, b)
		copy(buffer[30:], buffer[12:])

		width := 0
		height := 0
		alphaF := false

		if chunkID == "VP8 " {
			width = (((int(buffer[45]) << 8) | int(buffer[44])) & 0x3fff) - 1
			height = (((int(buffer[47]) << 8) | int(buffer[46])) & 0x3fff) - 1
		} else {
			width = (((int(buffer[40]) << 8) | int(buffer[39])) & 0x3fff) - 1
			height = ((int(buffer[42]) << 10) | (int(buffer[41]) << 2) | (int(buffer[40])>>6)&0x3fff) - 1
			// alpha flag
			if ((buffer[42] >> 4) & 0x1) == 0x1 {
				alphaF = true
			}
		}

		// VP8X
		buffer[12] = 0x56
		buffer[13] = 0x50
		buffer[14] = 0x38
		buffer[15] = 0x58

		// VP8X size
		buffer[16] = 0xa
		buffer[17] = 0x0
		buffer[18] = 0x0
		buffer[19] = 0x0

		// Flags
		buffer[20] = 0x4
		if alphaF {
			buffer[20] = buffer[20] | 0x10
		}

		// Reserved
		buffer[21] = 0x0
		buffer[22] = 0x0
		buffer[23] = 0x0

		// Canvas dimension
		buffer[24] = byte(width)
		buffer[25] = byte(width >> 8)
		buffer[26] = byte(width >> 16)
		buffer[27] = byte(height)
		buffer[28] = byte(height >> 8)
		buffer[29] = byte(height >> 16)

		chunkOffset += 18

	} else {
		if size%2 == 0 {
			buffer = make([]byte, chunkOffset+len(content)+8)
		} else {
			buffer = make([]byte, chunkOffset+len(content)+9)
		}
		copy(buffer, b)
		buffer[20] = buffer[20] | 0x4
	}

	// XMP ID
	buffer[chunkOffset] = 0x58
	buffer[chunkOffset+1] = 0x4d
	buffer[chunkOffset+2] = 0x50
	buffer[chunkOffset+3] = 0x20
	chunkOffset += 4

	// XMP Size
	buffer[chunkOffset] = byte(len(content))
	buffer[chunkOffset+1] = byte(len(content) >> 8)
	buffer[chunkOffset+2] = byte(len(content) >> 16)
	buffer[chunkOffset+3] = byte(len(content) >> 24)
	chunkOffset += 4

	// XMP Content
	copy(buffer[chunkOffset:], content)
	if size%2 != 0 {
		buffer[len(buffer)-1] = 0x0
		size++
	}

	// Update container size
	buffer[4] = byte(size)
	buffer[5] = byte(size >> 8)
	buffer[6] = byte(size >> 16)
	buffer[7] = byte(size >> 24)

	_, err = w.Write(buffer)

	return err
}
