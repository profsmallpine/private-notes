package web

import (
	"net/http"
	"text/template"
)

var TemplateDir = "./"

func (c *Controller) respondAuthenticated(w http.ResponseWriter, r *http.Request, templateToRender string, data interface{}) {
	tmpl := template.Must(template.New("authenticated_base.tmpl").ParseFiles(
		TemplateDir+"tmpl/layout/authenticated_base.tmpl",
		TemplateDir+"tmpl/partials/_header.tmpl",
		TemplateDir+"tmpl/"+templateToRender+".tmpl",
	))

	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (c *Controller) respondUnauthenticated(w http.ResponseWriter, r *http.Request, templateToRender string, data interface{}) {
	tmpl := template.Must(template.New("unauthenticated_base.tmpl").ParseFiles(
		TemplateDir+"tmpl/layout/unauthenticated_base.tmpl",
		TemplateDir+"tmpl/"+templateToRender+".tmpl",
	))

	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
