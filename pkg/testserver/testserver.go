package testserver

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

type testServer struct {
	*httptest.Server
}

func NewTestServer(t *testing.T, h http.Handler) *testServer {
	ts := httptest.NewServer(h)
	return &testServer{ts}
}

func (ts *testServer) Get(t *testing.T, urlPath string) (int, http.Header, []byte) {
	var body = []byte{}
	return ts.requestConfig(t, "GET", urlPath, body)
}

func (ts *testServer) Post(t *testing.T, urlPath string, body []byte) (int, http.Header, []byte) {
	return ts.requestConfig(t, "POST", urlPath, body)
}

func (ts *testServer) Put(t *testing.T, urlPath string, body []byte) (int, http.Header, []byte) {
	return ts.requestConfig(t, "PUT", urlPath, body)
}

func (ts *testServer) requestConfig(t *testing.T, action, urlPath string, body []byte) (int, http.Header, []byte) {
	client := &http.Client{}

	req, err := http.NewRequest(action, ts.URL+urlPath, bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	defer resp.Body.Close()

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("%+v", resp)
	return resp.StatusCode, resp.Header, body
}
