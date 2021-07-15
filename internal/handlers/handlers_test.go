package handlers

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

type postData struct {
	key   string
	value string
}

var theTests = []struct {
	name               string
	url                string
	method             string
	params             []postData
	expectedStatusCode int
}{
	{"home", "/", "GET", []postData{}, http.StatusOK},
	{"about", "/about", "GET", []postData{}, http.StatusOK},
	{"contact", "/contact", "GET", []postData{}, http.StatusOK},
	{"generals quarters", "/generals-quarters", "GET", []postData{}, http.StatusOK},
	{"majors suite", "/majors-suite", "GET", []postData{}, http.StatusOK},
	{"search availability", "/search-availability", "GET", []postData{}, http.StatusOK},
	{"make reservation", "/make-reservation", "GET", []postData{}, http.StatusOK},
	{"reservation summary", "/reservation-summary", "GET", []postData{}, http.StatusOK},
	{"post search availability", "/search-availability", "POST", []postData{
		{key: "start", value: "2021-07-01"},
		{key: "end", value: "2021-07-03"},
	}, http.StatusOK},
	{"post search availability json", "/search-availability-json", "POST", []postData{
		{key: "start", value: "2021-07-01"},
		{key: "end", value: "2021-07-03"},
	}, http.StatusOK},
	{"post make reservation", "/make-reservation", "POST", []postData{
		{key: "first_name", value: "John"},
		{key: "last_name", value: "Doe"},
		{key: "email", value: "me@here.com"},
		{key: "phone", value: "555-5555"},
	}, http.StatusOK},
}

func TestHandlers(t *testing.T) {
	routes := getRoutes()
	ts := httptest.NewTLSServer(routes)
	defer ts.Close()

	for _, e := range theTests {
		if e.method == "GET" {
			resp, err := ts.Client().Get(ts.URL + e.url)
			if err != nil {
				// t.Log(err)
				t.Fatal(err)
			}

			if resp.StatusCode != e.expectedStatusCode {
				t.Errorf("%s, expected %d, get %d", e.name, e.expectedStatusCode, resp.StatusCode)
			}
		} else {
			values := url.Values{}
			for _, x := range e.params {
				values.Add(x.key, x.value)
			}

			resp, err := ts.Client().PostForm(ts.URL+e.url, values)
			if err != nil {
				t.Log(err)
				t.Fatal(err)
			}

			if resp.StatusCode != e.expectedStatusCode {
				t.Errorf("%s, expected %d, get %d", e.name, e.expectedStatusCode, resp.StatusCode)
			}
		}
	}
}
