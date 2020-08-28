package handlers

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
)

// makeGETRequest against the given handler and target URL path.
// Returns the status code and response body.
func makeGETRequest(handler http.Handler, target string) (int, string) {
	req := httptest.NewRequest(http.MethodGet, target, nil)
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, req)
	result := recorder.Result()
	resultBody, err := ioutil.ReadAll(result.Body)
	if err != nil {
		panic(err)
	}
	return result.StatusCode, strings.TrimSpace(string(resultBody))
}
