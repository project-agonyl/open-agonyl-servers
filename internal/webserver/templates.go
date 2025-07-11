package webserver

import (
	"embed"
	"html/template"
	"io"
	"strings"

	"github.com/labstack/echo/v4"
)

//go:embed templates/*.html templates/pages/*.html templates/partials/*.html
var templateFS embed.FS

//go:embed assets
var staticFS embed.FS

type Templates struct {
	templates map[string]*template.Template
}

func NewTemplates() *Templates {
	baseTemplate, err := template.ParseFS(templateFS, "templates/base.html", "templates/partials/*.html")
	if err != nil {
		panic(err)
	}

	templates := make(map[string]*template.Template)

	pages := []string{"index", "about", "login", "characters"}
	for _, page := range pages {
		pageTmpl, err := baseTemplate.Clone()
		if err != nil {
			panic(err)
		}

		_, err = pageTmpl.ParseFS(templateFS, "templates/pages/"+page+".html")
		if err != nil {
			panic(err)
		}

		templates[page] = pageTmpl
	}

	return &Templates{
		templates: templates,
	}
}

func (t *Templates) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	if template, exists := t.templates[name]; exists {
		return template.ExecuteTemplate(w, "base", data)
	}

	if strings.HasPrefix(name, "login-") || strings.HasPrefix(name, "csrf-") {
		tmpl, err := template.ParseFS(templateFS, "templates/partials/"+name+".html")
		if err != nil {
			return err
		}
		return tmpl.ExecuteTemplate(w, name+".html", data)
	}

	baseTemplate, err := template.ParseFS(templateFS, "templates/base.html", "templates/partials/*.html")
	if err != nil {
		return err
	}
	return baseTemplate.ExecuteTemplate(w, "base", data)
}

func GetStaticFS() embed.FS {
	return staticFS
}
