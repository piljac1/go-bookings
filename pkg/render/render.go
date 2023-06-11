package render

import (
	"bytes"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/piljac1/bookings/pkg/config"
	"github.com/piljac1/bookings/pkg/models"
)

var app *config.AppConfig

func NewTemplates(a *config.AppConfig) {
	app = a
}

func AddDefaultData(templateData *models.TemplateData) *models.TemplateData {
	return templateData
}

func RenderTemplate(writer http.ResponseWriter, view string, templateData *models.TemplateData) {
	viewPath := strings.Replace(view, ".", "/", -1) + ".page.tmpl"

	var templateCache map[string]*template.Template

	if app.UseCache {
		templateCache = app.TemplateCache
	} else {
		var err error

		templateCache, err = CreateTemplateCache()

		if err != nil {
			log.Fatal(err)
		}
	}

	templateSet, ok := templateCache[viewPath]

	if !ok {
		log.Fatal("Cannot retrieve view from template cache")
	}

	buffer := new(bytes.Buffer)

	templateData = AddDefaultData(templateData)

	err := templateSet.Execute(buffer, templateData)

	if err != nil {
		log.Println(err)
	}

	_, err = buffer.WriteTo(writer)

	if err != nil {
		log.Println(err)
	}
}

func CreateTemplateCache() (map[string]*template.Template, error) {
	templateCache := map[string]*template.Template{}

	pages, err := filepath.Glob("./templates/*.page.tmpl")

	if err != nil {
		return templateCache, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		templateSet, err := template.New(name).ParseFiles(page)

		if err != nil {
			return templateCache, err
		}

		matches, err := filepath.Glob("./templates/*.layout.tmpl")

		if err != nil {
			return templateCache, err
		}

		if len(matches) > 0 {
			templateSet, err = templateSet.ParseGlob("./templates/*.layout.tmpl")

			if err != nil {
				return templateCache, err
			}
		}

		templateCache[name] = templateSet
	}

	return templateCache, nil
}
