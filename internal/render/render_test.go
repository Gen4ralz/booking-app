package render

import (
	"net/http"
	"testing"

	"github.com/gen4ralz/booking-app/internal/models"
)

func TestAddDefaultData(t *testing.T){
	var td models.TemplateData

	r, err := getSession()
	if err != nil {
		t.Error(err)
	}

	session.Put(r.Context(),"flash","123")

	result := AddDefaultData(&td, r)

	if result.Flash != "123" {
		t.Error("flash value of 123 not found in session")
	}
}

func TestRenderTemplate(t *testing.T){
	pathToTemplate = "./../../templates"
	tc, err := CreateTemplateCache()
	if err != nil {
		t.Error(err)
	}
	app.TemplateCache = tc

	r,err := getSession()
	if err !=nil {
		t.Error(err)
	}

	var ww myWriter

	err = Template(&ww, r, "home.gohtml", &models.TemplateData{})
	if err != nil {
		t.Error("error writing template to user")
	}

	err = Template(&ww, r, "non-existent.gohtml", &models.TemplateData{})
	if err == nil {
		t.Error("rendered template that does not exist")
	}
}

func getSession()(*http.Request,error){
	r, err := http.NewRequest("GET", "some", nil)
	if err != nil {
		return nil,err
	}

	ctx := r.Context()
	ctx,_ = session.Load(ctx, r.Header.Get("X-Session"))
	r = r.WithContext(ctx)
	return r,nil
}

func TestNewTemplates(t *testing.T) {
	NewRenderer(app)
}

func TestCreateTemplateCache(t *testing.T){
	pathToTemplate = "./../../templates"
	_, err := CreateTemplateCache()
	if err != nil {
		t.Error(err)
	}
}