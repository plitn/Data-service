package model

import "io"

type File struct {
	ID      string `json:"id" :"id"`
	Name    string `json:"name" :"name"`
	Size    int64  `json:"size" :"size"`
	Bytes   []byte `json:"file" :"bytes"`
	FileUrl string `json:"file_url"`
}

type CreateFileDTO struct {
	Name   string `json:"name"`
	Size   int64  `json:"size"`
	Reader io.Reader
}
