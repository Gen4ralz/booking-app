package render

import (
	"bytes"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/gen4ralz/booking-app/pkg/config"
	"github.com/gen4ralz/booking-app/pkg/models"
	"github.com/justinas/nosurf"
)

var app *config.AppConfig

func NewTemplates(a *config.AppConfig){
	app = a
}

func AddDefaultData(td *models.TemplateData, req *http.Request) *models.TemplateData{
	td.CSRFToken = nosurf.Token(req)
	return td
}

func RenderTemplate(res http.ResponseWriter, req *http.Request,tpl string, td *models.TemplateData) {

	var tc map[string]*template.Template
	if app.UseCache{
		tc = app.TemplateCache
	} else {
		tc,_ = CreateTemplateCache()
	}

	// get requested template from cache
	t, ok := tc[tpl]
	if !ok {
		log.Fatal("Could not get template from template cache")
	}

	buf := new(bytes.Buffer)

	td = AddDefaultData(td, req)

	_ = t.Execute(buf,td)

	// render the template
	_, err := buf.WriteTo(res)
	if err != nil {
		log.Println(err)
	}

}

func CreateTemplateCache()(map[string]*template.Template, error){
	// myCache := make(map[string]*template.Template)
	myCache := map[string]*template.Template{}

	// get all of the files named *.gohtml from ./templates
	pages, err := filepath.Glob("./templates/*.gohtml")
	if err != nil {
		return myCache,err
	}

	// range through all files ending with *.gohtml
	for _,page := range pages {
		name := filepath.Base(page)

		//ts = template set
		ts, err := template.New(name).ParseFiles(page)
		if err != nil {
			return myCache, err
		}

		matches, err := filepath.Glob("./templates/*.layout.gohtml")
		if err != nil {
			return myCache, err
		}

		if len(matches) > 0 {
			ts, err = ts.ParseGlob("./templates/*.layout.gohtml")
			if err != nil {
				return myCache, err
			}
		}
		myCache[name] = ts
	}
	return myCache, nil
}