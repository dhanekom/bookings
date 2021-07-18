package render

import (
	"net/http"
	"testing"

	"github.com/dhanekom/bookings/internal/models"
)

func TestAddDefaultData(t *testing.T) {
	var td models.TemplateData

	r, err := getSession()
	if err != nil {
		t.Error(err)
	}

	session.Put(r.Context(), "flash", "123")
	result := AddDefaultData(&td, r)
	if result.Flash != "123" {
		t.Errorf("Flash value - expected %s, got %s", "1234", result.Flash)
	}
}

func TestRenderTemplate(t *testing.T) {
	app.TemplatePath = "../../templates"

	tc, err := CreateTemplateCache(app.TemplatePath)
	if err != nil {
		t.Errorf("unable to create template cache - %s", err)
	}

	app.TemplateCache = tc

	r, err := getSession()
	if err != nil {
		t.Error(err)
	}

	var ww myResponse

	err = Template(&ww, r, "home.page.tmpl", &models.TemplateData{})
	if err != nil {
		t.Error("error writing template to browser")
	}

	err = Template(&ww, r, "non-existent.page.tmpl", &models.TemplateData{})
	if err == nil {
		t.Error("rendered templated that does not exist")
	}
}

type myResponse struct{}

func (rw *myResponse) Header() http.Header {
	var h http.Header
	return h
}

func (rw *myResponse) Write(b []byte) (int, error) {
	length := len(b)
	return length, nil
}

func (rw *myResponse) WriteHeader(statusCode int) {

}

func getSession() (*http.Request, error) {
	r, err := http.NewRequest("GET", "some-url", nil)
	if err != nil {
		return nil, err
	}

	ctx := r.Context()
	ctx, _ = session.Load(ctx, r.Header.Get("X-Session"))
	r = r.WithContext(ctx)

	return r, nil
}

func TestNewTemplate(t *testing.T) {
	NewRendered(app)
}

func TestCreateTemplateCache(t *testing.T) {
	app.TemplatePath = "../../templates"
	_, err := CreateTemplateCache(app.TemplatePath)
	if err != nil {
		t.Error(err)
	}

}
