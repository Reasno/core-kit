package kits3

import "io"

// Request models a go kit request in UploadEndpoint
type Request struct {
	name string
	data io.Reader
}

// Response models a go kit Response in UploadEndpoint
type Response struct {
	Data struct {
		Url string `json:"url"`
	} `json:"data"`
	Code    int    `json:"code"`
	Message string `json:"message,omitempty"`
}
