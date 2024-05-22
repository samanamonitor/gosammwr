package transport

import (
	"io"
	"bytes"
	"errors"
	"fmt"
)

type MultiPart struct {
	boundary string
	headers map[string]string

	body []byte
	buffer *bytes.Buffer
	encoded bool
}

func (mp *MultiPart) Body() (*bytes.Buffer) {
	if ! mp.encoded {
		mp.Encode()
		mp.encoded = true
	}
	return mp.buffer
}

func (mp *MultiPart) AddHeader(key string, value string) {
	mp.headers[key] = value
}

func (mp *MultiPart) WriteBoundary(end bool) {
	mp.buffer.WriteString("--")
	mp.buffer.WriteString(mp.boundary)
	if end {
		mp.buffer.WriteString("--")
	}
	mp.buffer.WriteString("\r\n")
}

func (mp *MultiPart) GetHeader(key string) (value string) {
	return fmt.Sprintf("%s: %s\r\n", key, mp.headers[key])
}

func (mp *MultiPart) Encode() () {
	/*
	NOTE IMPORTANT!!! Windows will reject requests with the headers in a different order
	https://learn.microsoft.com/en-us/openspecs/windows_protocols/ms-wsmv/080c0681-bc50-412d-b9ca-42047a8011d5
	*/
	mp.WriteBoundary(false)
	mp.buffer.WriteString(mp.GetHeader("Content-Type"))
	mp.buffer.WriteString(mp.GetHeader("OriginalContent"))
	mp.WriteBoundary(false)
	mp.buffer.WriteString("Content-Type: application/octet-stream\r\n")
	mp.buffer.Write(mp.body)
	mp.WriteBoundary(true)
}

func (mp *MultiPart) Decode() (error) {
	return nil
}

func NewMultiPart(Boundary string, message []byte) (*MultiPart) {
	return &MultiPart{
		buffer: bytes.NewBuffer(nil),
		headers: make(map[string]string),
		boundary: Boundary,
		body: message,
		encoded: false,
	}
}


func removeBlanks(data []byte) ([]byte) {
	for ;len(data) > 1; {
		if data[0] == ' ' || data[0] == '\t' {
			data = data[1:]
		} else {
			return data
		}
	}
	return data
}

func extractPart(data []byte, boundary []byte) ([]byte, []byte, error) {
	bl := len(boundary)
	if ! bytes.HasPrefix(data, boundary) {
		return data, []byte{}, errors.New("Missing starting boundary")
	}
	if data[bl] == '-' && data[bl+1] == '-' {
		return data[bl+2:], []byte{}, io.EOF
	}
	sectionData := data[bl+2:]
	i := bytes.Index(sectionData, boundary)
	if i == -1 {
		return data, []byte{}, errors.New("Boundary not found")
	}
	if i + bl + 2 > len(sectionData) {
		return data, []byte{}, errors.New("Invalid data after boundary")
	}
	return sectionData[i:], sectionData[:i], nil
}

func extractHeader(data []byte) ([]byte, string, string, error) {
	data = removeBlanks(data)
	nl := bytes.IndexByte(data, '\n')
	if nl == -1 {
		return data, "", "", errors.New("Invalid Header. No new line found.")
	}
	header := data[:nl]
	data = data[nl+1:]

	sep := bytes.IndexByte(header, ':')
	if sep == -1 {
		return data, "", "", errors.New("Invalid Header. No ':' found")
	}
	key := string(header[:sep])
	value := string(bytes.Trim(removeBlanks(header[sep+1:]), "\r"))
	return data, key, value, nil
}

func MultipartDecode(data []byte, boundary_string string) (map[string]string, []byte, error) {
	headers := make(map[string]string)
	body := []byte{}
	boundary := []byte(fmt.Sprint("--", boundary_string))

	for {
		var section []byte
		var err error
		data, section, err = extractPart(data, boundary)
		if err == io.EOF {
			break
		} else if err != nil {
			return headers, body, err
		}
		for ; len(section) > 0;  {
			var key, value string
			section, key, value, err = extractHeader(section)
			if err != nil {
				return headers, body, err
			}
			if key == "Content-Type" && value == "application/octet-stream" {
				body = section
				break
			}
			headers[key] = value
		}
	}
	return headers, body, nil
}
