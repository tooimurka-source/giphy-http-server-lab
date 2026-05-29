package server

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

type fakeSearcher struct {
	url string
	err error
}

func (f fakeSearcher) Search(query string) (string, error) {
	if f.err != nil {
		return "", f.err
	}

	return f.url, nil
}

func TestSearchGIFByQuery(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := NewRouter(fakeSearcher{url: "https://media.giphy.com/example.gif"})

	request := httptest.NewRequest(http.MethodGet, "/api/gifs/search?q=cat", nil)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, response.Code)
	}

	var body map[string]string
	if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if body["query"] != "cat" {
		t.Fatalf("expected query cat, got %q", body["query"])
	}

	if body["url"] != "https://media.giphy.com/example.gif" {
		t.Fatalf("unexpected url %q", body["url"])
	}
}

func TestSearchGIFRequiresQuery(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := NewRouter(fakeSearcher{url: "https://media.giphy.com/example.gif"})

	request := httptest.NewRequest(http.MethodGet, "/api/gifs/search", nil)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, response.Code)
	}
}

func TestSearchGIFNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := NewRouter(fakeSearcher{err: errors.New("no gifs found")})

	request := httptest.NewRequest(
		http.MethodPost,
		"/api/gifs/search",
		strings.NewReader(`{"query":"unknown"}`),
	)
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, response.Code)
	}
}
