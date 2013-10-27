package webloop

import (
	"net/http"
	"net/http/httptest"
)

var (
	// mux is the HTTP request multiplexer used with the test server.
	mux *http.ServeMux

	// server is a test HTTP server used to provide mock API responses.
	server *httptest.Server
)

// setup sets up a test HTTP server. Tests should register handlers on mux which
// provide mock responses for the method being tested.
func setup() {
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)
}

// teardown closes the test HTTP server.
func teardown() {
	server.Close()
}
