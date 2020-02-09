package cmd

import "net/http"

// NewHTTPServer returns a new HTTP server
func NewHTTPServer() *http.Server {
	return &http.Server{}
}
