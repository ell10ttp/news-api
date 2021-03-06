package application

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"news-api/pkg/models"
	"news-api/pkg/sourceapi"
	"news-api/pkg/testserver"

	"github.com/stretchr/testify/assert"
)

func TestPing(t *testing.T) {
	t.Run("Returns 200 OK", func(t *testing.T) {
		app := newTestApplication()
		ts := testserver.NewTestServer(t, app.Routes())
		defer ts.Close()

		code, _, _ := ts.Get(t, "/ping")

		assertStatus(t, code, http.StatusOK)
	})
	t.Run("Returns correct response", func(t *testing.T) {
		app := newTestApplication()
		ts := testserver.NewTestServer(t, app.Routes())
		defer ts.Close()

		_, _, body := ts.Get(t, "/ping")

		got := strings.Trim(string(body), "\n\"")
		want := "ping successful"

		assert.Equal(t, want, got)
	})

}

func TestGetSourceList(t *testing.T) {
	t.Run("Returns 200 OK", func(t *testing.T) {
		app := newTestApplication()
		ts := testserver.NewTestServer(t, app.Routes())
		defer ts.Close()

		code, _, _ := ts.Get(t, "/source")

		assertStatus(t, http.StatusOK, code)
	})
	t.Run("Returns correct response", func(t *testing.T) {
		app := newTestApplication()
		ts := testserver.NewTestServer(t, app.Routes())
		defer ts.Close()

		_, _, body := ts.Get(t, "/source")

		var got models.SourceList
		want := app.SourceAPI.GetSourceList()
		if err := json.Unmarshal(body, &got); err != nil {
			t.Fail()
		}
		assert.Equal(t, want, got)
	})
}

func TestGetSource(t *testing.T) {
	t.Run("Returns 200 OK", func(t *testing.T) {
		app := newTestApplication()
		ts := testserver.NewTestServer(t, app.Routes())
		defer ts.Close()

		code, _, _ := ts.Get(t, "/source/1")

		assertStatus(t, http.StatusOK, code)
	})
	t.Run("Returns correct response", func(t *testing.T) {
		app := newTestApplication()
		ts := testserver.NewTestServer(t, app.Routes())
		defer ts.Close()

		_, _, body := ts.Get(t, "/source/1")

		var got models.Source
		fmt.Println(string(body))
		want := app.SourceAPI.GetSourceList().Sources[0]
		if err := json.Unmarshal(body, &got); err != nil {
			t.Fail()
		}
		assert.Equal(t, want, got)
	})
}

func TestPostSource(t *testing.T) {
	t.Run("Returns 200 OK", func(t *testing.T) {
		app := newTestApplication()
		ts := testserver.NewTestServer(t, app.Routes())
		defer ts.Close()

		values := map[string]string{"Name": "test", "Description": "tester", "Url": "http://test.com", "Language": "en-test", "Country": "testland"}
		jsonValue, _ := json.Marshal(values)
		code, _, _ := ts.Post(t, "/source", jsonValue)
		assertStatus(t, http.StatusOK, code)
	})
	t.Run("Source ID returns as a valid and is able to be retrieved after the fact", func(t *testing.T) {
		app := newTestApplication()
		ts := testserver.NewTestServer(t, app.Routes())
		defer ts.Close()

		values := map[string]string{"Name": "test", "Description": "tester", "Url": "http://test.com", "Language": "en-test", "Country": "testland"}
		jsonValue, _ := json.Marshal(values)

		_, _, body := ts.Post(t, "/source", jsonValue)

		type Response struct {
			Action     string `json:"action"`
			Successful bool   `json:"Successful"`
			SourceID   int    `json:"sourceId"`
		}
		var resp Response
		if err := json.Unmarshal(body, &resp); err != nil {
			t.Fail()
		}
		gotSource, _ := app.SourceAPI.GetSource(resp.SourceID)

		assert.Equal(t, gotSource.ID, resp.SourceID)
	})
}

func TestGetSourceCategories(t *testing.T) {
	t.Run("Returns 200 OK", func(t *testing.T) {
		app := newTestApplication()
		ts := testserver.NewTestServer(t, app.Routes())
		defer ts.Close()

		code, _, _ := ts.Get(t, "/source/1/category")

		assertStatus(t, http.StatusOK, code)
	})
	t.Run("Returns correct response", func(t *testing.T) {
		app := newTestApplication()
		ts := testserver.NewTestServer(t, app.Routes())
		defer ts.Close()

		_, _, body := ts.Get(t, "/source/1/category")

		type Response struct {
			Action             string   `json:"action"`
			Successful         bool     `json:"successful"`
			NumberOfCategories int      `json:"numberOfCategories"`
			Categories         []string `json:"categories"`
		}
		var got Response
		fmt.Println(string(body))
		want := []string{"business", "entertainment", "politics", "technology", "uk", "world"}
		if err := json.Unmarshal(body, &got); err != nil {
			t.Fail()
		}
		assert.Equal(t, 6, len(got.Categories))
		assert.Equal(t, want, got.Categories)
	})
}

// NewTestApplication create instance of application for testing
func newTestApplication() *Model {
	app := Init()
	app.SetSourceAPI(sourceapi.NewSourceAPI())
	return app
}

func assertStatus(t *testing.T, want int, got int) {
	t.Helper()
	if got != want {
		t.Errorf("got %d want %d", got, want)
	}
}
