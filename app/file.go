package main

import (
	"bytes"
	"fmt"
	"strings"
)

type File struct {
	Header   map[string]string
	MetaData map[string]string
	Data     []byte
}

func NewFile() *File {
	return &File{
		Header:   map[string]string{},
		MetaData: map[string]string{},
		Data:     []byte{},
	}
}

func (f *File) Parse(data []byte) {
	done := false
	step := 0
	sequence := []byte("\r\n")
	var results []string
	var line string
	var index int
	for !done {
		switch step {
		case 0:
			index = bytes.Index(data, sequence) + 2
			line = string(data[:index])
			// fmt.Printf("line: %s\n", line)
			data = data[index:]

			if line == "\r\n" {
				results = strings.Split(f.Header["Content-Disposition"], ";")
				for _, result := range results {
					meta := strings.Split(strings.TrimSpace(result), "=")
					if len(meta) == 2 {
						key := strings.TrimSpace(meta[0])
						value := strings.TrimSpace(meta[1])
						f.MetaData[key] = value
					}
				}
				step = 1
			} else {
				results = strings.Split(line, ":")
				key := strings.TrimSpace(results[0])
				value := strings.TrimSpace(results[1])
				fmt.Printf("key: %s, value: %s\n", key, value)
				f.Header[key] = value
			}
		case 1:
			f.Data = []byte{}
			f.Data = append(f.Data, data...)
			done = true
			fmt.Printf("Data length: %d\n", len(f.Data))
		}
	}
}
