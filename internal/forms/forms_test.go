package forms

import (
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestForm_Valid(t *testing.T) {
	r := httptest.NewRequest("GET", "/some-url", nil)
	form := New(r.PostForm)

	isValid := form.Valid()
	if !isValid {
		t.Error("got invalid when should have been valid")
	}
}

func TestForm_Required(t *testing.T) {
	r := httptest.NewRequest("POST", "/some-url", nil)
	form := New(r.PostForm)

	form.Required("a", "b", "c")
	if form.Valid() {
		t.Error("form shows valid when there are required fields")
	}

	postedData := url.Values{}
	postedData.Add("a", "a")
	postedData.Add("b", "a")
	postedData.Add("c", "a")

	form = New(postedData)
	form.Required("a", "b", "c")
	if !form.Valid() {
		t.Error("shows does not have required fields when it does")
	}
}

func TestForm_MinLength(t *testing.T) {
	postData := url.Values{}
	postData.Add("a", "a")
	postData.Add("b", "bbb")
	postData.Add("c", "ccccc")

	form := New(postData)
	form.MinLength("a", 2)

	if form.Errors.Get("a") == "" {
		t.Error(`expected errors for field "a"`)
	}

	if form.Valid() {
		t.Error("form MinLength - expected failure but got success")
	}

	form = New(postData)
	form.MinLength("c", 5)

	if form.Errors.Get("c") != "" {
		t.Error(`did not expect errors for field "c"`)
	}

	if !form.Valid() {
		t.Error("form MinLength - expected success but failed")
	}
}

var testEmails = []struct {
	description    string
	value          string
	expectedResult bool
}{
	{"valid email", "test@me.com", true},
	{"blank email", "", false},
	{"invalid domain", "test@me.", false},
	{"no @", "test.com", false},
}

func TestForm_IsEmail(t *testing.T) {
	for _, e := range testEmails {
		postData := url.Values{}
		postData.Add("email", e.value)

		form := New(postData)
		result := form.IsEmail("email")
		if result != e.expectedResult {
			t.Errorf("%s - expected %s to be valid %v, got %v", e.description, e.value, e.expectedResult, result)
		}
	}
}

func TestForm_Has(t *testing.T) {
	postData := url.Values{}
	postData.Add("a", "a")
	postData.Add("b", "bbb")
	postData.Add("c", "ccccc")

	// r, _ := http.NewRequest("POST", "/some-url", nil)
	// r.PostForm = postData

	form := New(postData)
	form.Has("d")
	if form.Valid() {
		t.Error(`should not have found field "e"`)
	}

	form = New(postData)
	form.Has("a")
	if !form.Valid() {
		t.Error(`should have found field "a"`)
	}
}
