package comms

import (
	"context"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEmailer_SendWithTemplate(t *testing.T) {
	t.Run("sends email with template", func(t *testing.T) {
		e, cleanup := newEmailer(func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, http.MethodPost, r.Method)
			require.Equal(t, "/email/withTemplate", r.URL.Path)
			require.Equal(t, "transactional", r.Header.Get("X-Postmark-Server-Token"))

			body := readBody(r.Body)
			require.Contains(t, body, `"From":"Me \u003cme@example.com\u003e"`)
			require.Contains(t, body, `"To":"you@example.com"`)
			require.Contains(t, body, `"TemplateAlias":"template1"`)
			require.Contains(t, body, `"TemplateModel":{"key1":"value1"}`)
		})
		defer cleanup()

		err := e.SendWithTemplate(context.Background(), CreateNameAndEmail("Me", "me@example.com"), "you@example.com", "template1", map[string]string{
			"key1": "value1",
		})
		require.NoError(t, err)
	})

	t.Run("errors on cannot connect", func(t *testing.T) {
		e := NewEmailer("")
		e.endpoint = "http://localhost:12345"

		err := e.SendWithTemplate(context.Background(), "me@example.com", "you@example.com", "template1", map[string]string{})
		require.Error(t, err)
	})

	t.Run("errors on http error response", func(t *testing.T) {
		e, cleanup := newEmailer(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "boo", http.StatusBadRequest)
		})
		defer cleanup()

		err := e.SendWithTemplate(context.Background(), "me@example.com", "you@example.com", "template1", map[string]string{})
		require.Error(t, err)
	})
}

func TestEmailer_Send(t *testing.T) {
	t.Run("sends email", func(t *testing.T) {
		e, cleanup := newEmailer(func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, http.MethodPost, r.Method)
			require.Equal(t, "/email", r.URL.Path)
			require.Equal(t, "transactional", r.Header.Get("X-Postmark-Server-Token"))

			body := readBody(r.Body)
			require.Contains(t, body, `"From":"Me \u003cme@example.com\u003e"`)
			require.Contains(t, body, `"To":"you@example.com"`)
			require.Contains(t, body, `"Subject":"Subjective"`)
			require.Contains(t, body, `"HtmlBody":"\u003chtml\u003e\u003cbody\u003eHi!\u003c/body\u003e\u003c/html\u003e"`)
		})
		defer cleanup()

		err := e.Send(context.Background(), CreateNameAndEmail("Me", "me@example.com"), "you@example.com", "Subjective", "<html><body>Hi!</body></html>")
		require.NoError(t, err)
	})
}

func newEmailer(h http.HandlerFunc) (*Emailer, func()) {
	s := httptest.NewServer(h)
	e := NewEmailer("transactional")
	e.endpoint = s.URL
	cleanup := func() {
		s.Close()
	}
	return e, cleanup
}

func readBody(r io.ReadCloser) string {
	body, err := ioutil.ReadAll(r)
	if err != nil {
		panic(err)
	}
	return string(body)
}
