/*
 * Copyright 2022-present Kuei-chun Chen. All rights reserved.
 * http_reader_test.go
 */

package hatchet

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetHTTPContent(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "hello, client")
	}))
	defer server.Close()

	reader, err := GetHTTPContent(server.URL, "", "")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	content, err := io.ReadAll(reader)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	expected := "hello, client\n"
	if string(content) != expected {
		t.Errorf("unexpected content: got %q, want %q", string(content), expected)
	}
}
