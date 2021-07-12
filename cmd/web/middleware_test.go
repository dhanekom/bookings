package main

import (
	"net/http"
	"testing"
)

func TestNoSurf(t *testing.T) {
	var myH myHandler
	h := NoSurf(&myH)

	switch v := h.(type) {
	case http.Handler:
		// Test passed
	default:
		t.Errorf("got type %T, expected type http.Handler", v)
	}
}

func TestSessionLoad(t *testing.T) {
	var myH myHandler
	h := SessionLoad(&myH)

	switch v := h.(type) {
	case http.Handler:
		// Test passed
	default:
		t.Errorf("got type %T, expected type http.Handler", v)
	}
}
