package application

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPing(t *testing.T) {
	t.Run("Returns 200 OK", func(t *testing.T) {
		app := newTestApplication()
		req := httptest.NewRequest(http.MethodGet, "/ping", nil)
		w := httptest.NewRecorder()
		app.ping(w, req)
		res := w.Result()
		defer res.Body.Close()
		data, err := ioutil.ReadAll(res.Body)
		if err != nil {
			t.Errorf("expected err to be nil got %v", err)
		}
		assert.Equal(t, "\"ping successful\"", string(data))
	})
}

// NewTestApplication create instance of application for testing
func newTestApplication() *Model {
	app := Init()
	return app
}

func assertStatus(t *testing.T, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("got %d want %d", got, want)
	}
}
