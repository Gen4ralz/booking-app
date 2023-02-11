package render

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/gen4ralz/booking-app/internal/config"
	"github.com/gen4ralz/booking-app/internal/models"
	"github.com/justinas/nosurf"
)

var functions = template.FuncMap{}

var app *config.AppConfig
var pathToTemplate = "./templates"

func NewTemplates(a *config.AppConfig){
	app = a
}

func AddDefaultData(td *models.TemplateData, req *http.Request) *models.TemplateData{
	td.Flash = app.Session.PopString(req.Context(), "flash")
	td.Error = app.Session.PopString(req.Context(), "error")
	td.Warning = app.Session.PopString(req.Context(), "warning")
	td.CSRFToken = nosurf.Token(req)
	return td
}

func RenderTemplate(res http.ResponseWriter, req *http.Request,tpl string, td *models.TemplateData) error {

	var tc map[string]*template.Template
	if app.UseCache{
		tc = app.TemplateCache
	} else {
		tc,_ = CreateTemplateCache()
	}

	// get requested template from cache
	t, ok := tc[tpl]
	if !ok {
		// log.Println("Could not get template from template cache")
		return errors.New("can't get template from cache")
	}

	buf := new(bytes.Buffer)

	td = AddDefaultData(td, req)

	_ = t.Execute(buf,td)

	// render the template
	_, err := buf.WriteTo(res)
	if err != nil {
		log.Println("Error writing template to user", err)
		return err

	}
	return nil

}

func CreateTemplateCache()(map[string]*template.Template, error){
	// myCache := make(map[string]*template.Template)
	myCache := map[string]*template.Template{}

	// get all of the files named *.gohtml from ./templates
	pages, err := filepath.Glob(fmt.Sprintf("%s/*.gohtml", pathToTemplate))
	if err != nil {
		return myCache,err
	}

	// range through all files ending with *.gohtml
	for _,page := range pages {
		name := filepath.Base(page)

		//ts = template set
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return myCache, err
		}

		matches, err := filepath.Glob(fmt.Sprintf("%s/*.layout.gohtml", pathToTemplate))
		if err != nil {
			return myCache, err
		}

		if len(matches) > 0 {
			ts, err = ts.ParseGlob(fmt.Sprintf("%s/*.layout.gohtml", pathToTemplate))
			if err != nil {
				return myCache, err
			}
		}
		myCache[name] = ts
	}
	return myCache, nil
}