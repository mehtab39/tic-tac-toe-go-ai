// server/server_test.go

package server_test

import (
	"github.com/stretchr/testify/assert"
	"go-be-ai/server"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPingHandler(t *testing.T) {
	// Create a request to the /ping endpoint
	req, err := http.NewRequest("GET", "/ping", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Call the PingHandler function directly, passing in the ResponseRecorder and Request
	server.PingHandler(rr, req)

	// Check the status code and response body
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "pong", rr.Body.String())
}
