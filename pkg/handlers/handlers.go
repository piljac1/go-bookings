package handlers

import (
	"net/http"

	"github.com/piljac1/bookings/pkg/config"
	"github.com/piljac1/bookings/pkg/models"
	"github.com/piljac1/bookings/pkg/render"
)

var Repo *Repository

type Repository struct {
	App *config.AppConfig
}

func NewRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
	}
}

func NewHandlers(r *Repository) {
	Repo = r
}

func (m *Repository) Home(writer http.ResponseWriter, request *http.Request) {
	remoteIP := request.RemoteAddr

	m.App.Session.Put(request.Context(), "remote_ip", remoteIP)

	render.RenderTemplate(writer, "home", &models.TemplateData{})
}

func (m *Repository) About(writer http.ResponseWriter, request *http.Request) {
	stringMap := map[string]string{}
	stringMap["test"] = "Hello, again."

	remoteIP := m.App.Session.GetString(request.Context(), "remote_ip")

	stringMap["remote_ip"] = remoteIP

	render.RenderTemplate(writer, "about", &models.TemplateData{
		StringMap: stringMap,
	})
}
