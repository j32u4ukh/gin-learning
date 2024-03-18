package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"image/png"
	"os"
	"strings"
)

type Request struct {
	Method   string
	Query    string
	Proto    string
	Header   map[string]string
	Data     []byte
	Files    []*File
	JsonData map[string]string
}

func NewRequest() *Request {
	return &Request{
		Header: map[string]string{},
		// Data: []byte{},
	}
}

func (r *Request) Parse(data []byte) error {
	done := false
	step := 0
	sequence := []byte("\r\n")
	var results []string
	var line string
	var index int
	for !done {
		if step != 2 {
			index = bytes.Index(data, sequence) + 2
			line = string(data[:index])
			data = data[index:]
		}
		switch step {
		case 0:
			results = strings.Split(line, " ")
			r.Method = strings.TrimSpace(results[0])
			r.Query = strings.TrimSpace(results[1])
			r.Proto = strings.TrimSpace(results[2])
			fmt.Printf("Method: %s, Query: %s, Proto: %s\n", r.Method, r.Query, r.Proto)
			step = 1
		case 1:
			if line == "\r\n" {
				if length, ok := r.Header["Content-Length"]; ok {
					if length == "0" {
						done = true
						break
					}
				}
				step = 2
			} else {
				results = strings.Split(line, ":")
				key := strings.TrimSpace(results[0])
				value := strings.TrimSpace(results[1])
				fmt.Printf("key: %s, value: %s\n", key, value)
				r.Header[key] = value
			}
		case 2:
			r.Data = []byte{}
			r.Data = append(r.Data, data...)
			done = true
		}
	}
	r.ParseContent()
	return nil
}

func (r *Request) ParseContent() error {
	cts := strings.Split(r.Header["Content-Type"], ";")
	ctype := strings.TrimSpace(cts[0])
	var err error
	switch ctype {
	case "application/json":
		r.JsonData = map[string]string{}
		err = json.Unmarshal(r.Data, &r.JsonData)
		if err != nil {
			return err
		}
	case "multipart/form-data":
		key, boundary, ok := strings.Cut(strings.TrimSpace(cts[1]), "=")
		if !ok || key != "boundary" {
			return errors.New("not found boundary")
		}
		var index int

		// 移除結尾 boundary
		endSequence := []byte(fmt.Sprintf("%s--", boundary))
		index = bytes.Index(r.Data, endSequence)
		if index == -1 {
			return errors.New("not found end boundary")
		}

		r.Data = r.Data[:index]
		r.Files = []*File{}
		sequence := []byte(fmt.Sprintf("%s\r\n", boundary))
		nSequence := len(sequence)

		// LastIndex
		index = bytes.LastIndex(r.Data, sequence)

		for index != -1 {
			data := []byte(r.Data[index+nSequence:])
			fmt.Printf("data: %+v\n", data)

			file := NewFile()
			file.Parse(data)
			fmt.Printf("Header: %+v\n", file.Header)
			fmt.Printf("Meta: %+v\n", file.MetaData)
			isPng := file.Header["Content-Type"] == "image/png"
			if(isPng){
				savePng(file.Data)
			}
			r.Files = append(r.Files, file)
			r.Data = r.Data[:index]
			index = bytes.LastIndex(r.Data, sequence)
		}

	default:
		// req += SliceToString(r.Data)
	}
	return nil
}

func savePng(pngData []byte){
	// 创建一个字节缓冲区，并将PNG二进制数据写入其中
	buf := bytes.NewBuffer(pngData)

	// 解码PNG数据
	img, err := png.Decode(buf)
	if err != nil {
		fmt.Println("Error decoding PNG data:", err)
		return
	}

	// 创建一个新的文件用于保存解码后的图像
	outFile, err := os.Create("./output.png") 
	if err != nil {
		fmt.Println("Error creating output file:", err)
		return
	}
	defer outFile.Close()

	// 将图像保存到文件中
	err = png.Encode(outFile, img)
	if err != nil {
		fmt.Println("Error encoding image:", err)
		return
	}

	fmt.Println("Image saved successfully")
}

func (r Request) String() string {
	req := fmt.Sprintf("%s %s %s\r\n", r.Method, r.Query, r.Proto)
	for k, v := range r.Header {
		req += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	req += "\r\n"
	// Content-Type: application/json
	ctype := r.Header["Content-Type"]
	switch ctype {
	case "application/json":
		req += fmt.Sprintf("JsonData: %+v", r.JsonData)
	default:
		req += SliceToString(r.Data)
	}
	return req
}
