package webserver

import (
	"embed"
	"html/template"
	"io"
	"io/fs"
	"path/filepath"

	"github.com/labstack/echo/v4"
)

//go:embed templates/*.html templates/pages/*.html templates/partials/*.html
var templateFS embed.FS

//go:embed assets
var staticFS embed.FS

type Templates struct {
	templates *template.Template
}

func NewTemplates() *Templates {
	tmpl := template.New("")

	err := parseTemplatesFromFS(tmpl, templateFS, "templates")
	if err != nil {
		panic(err)
	}

	return &Templates{
		templates: tmpl,
	}
}

func parseTemplatesFromFS(tmpl *template.Template, fsys embed.FS, dir string) error {
	return fs.WalkDir(fsys, dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() || filepath.Ext(path) != ".html" {
			return nil
		}

		content, err := fsys.ReadFile(path)
		if err != nil {
			return err
		}

		name := filepath.Base(path)
		name = name[:len(name)-len(filepath.Ext(name))]

		_, err = tmpl.New(name).Parse(string(content))
		return err
	})
}

func (t *Templates) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func GetStaticFS() embed.FS {
	return staticFS
}
