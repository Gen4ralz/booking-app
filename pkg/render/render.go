package render

import (
	"bytes"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/gen4ralz/booking-app/pkg/config"
	"github.com/gen4ralz/booking-app/pkg/models"
)

// var functions = template.FuncMap{}

var app *config.AppConfig

// NewTemplates sets the config for the template package
func NewTemplates(a *config.AppConfig) {
	app = a
}

func AddDefaultData(td *models.TemplateData) *models.TemplateData{
	return td
}

func RenderTemplate(res http.ResponseWriter, tpl string, td *models.TemplateData){
	var tc map[string]*template.Template
	// because in development mode we prefer to re-read template from disk of every request.
	// But in production mode we need to use cache template
	if app.UseCache {
		tc = app.TemplateCache
	} else {
		tc,_ = CreateTemplateCache()
	}
	// get the template cache from the app config
	//Instead of creating the template cache, let's just use it.
	// tc := app.TemplateCache


	// get requested template from cache
	t, ok := tc[tpl]
	if !ok {
		log.Fatal("could not get template from template cache")
	}

	buf := new(bytes.Buffer)

	td = AddDefaultData(td)

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