package handlers

import (
	"context"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
)

// getRequest against the given handler and target URL path.
// Returns the status code, headers, and response body.
func getRequest(handler http.Handler, target string) (int, http.Header, string) {
	return request(handler, http.MethodGet, target, nil, "")
}

// postRequest against the given handler and target URL path.
// Returns the status code, headers, and response body.
func postRequest(handler http.Handler, target, body string) (int, http.Header, string) {
	return request(handler, http.MethodPost, target, nil, body)
}

// putRequest against the given handler and target URL path.
// Returns the status code, headers, and response body.
func putRequest(handler http.Handler, target, body string) (int, http.Header, string) {
	return request(handler, http.MethodPut, target, nil, body)
}

// deleteRequest against the given handler and target URL path.
// Returns the status code, headers, and response body.
func deleteRequest(handler http.Handler, target string) (int, http.Header, string) {
	return request(handler, http.MethodDelete, target, nil, "")
}

// request against the given handler and target URL path.
// Returns the status code, headers, and response body.
func request(handler http.Handler, method, target string, ctx context.Context, body string) (int, http.Header, string) {
	var bodyReader io.Reader
	if body != "" {
		bodyReader = strings.NewReader(body)
	}
	req := httptest.NewRequest(method, target, bodyReader)
	if method == http.MethodPost || method == http.MethodPut {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}
	if ctx != nil {
		req = req.WithContext(ctx)
	}
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, req)
	result := recorder.Result()
	resultBody, err := ioutil.ReadAll(result.Body)
	if err != nil {
		panic(err)
	}
	return result.StatusCode, result.Header, strings.TrimSpace(string(resultBody))
}

func noopHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {}
}
