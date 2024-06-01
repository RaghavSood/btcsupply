package templates

import (
	"embed"
	"html/template"
	"io"
	"time"
)

//go:embed *
var Templates embed.FS

type Template struct {
	templates *template.Template
}

func New() *Template {
	funcMap := template.FuncMap{
		"now": func() interface{} { return time.Now() },
	}

	templates := template.Must(template.New("").Funcs(funcMap).ParseFS(Templates, "footer.tmpl", "header.tmpl", "base.tmpl"))
	return &Template{
		templates: templates,
	}
}

// func (t *Template) Render(w io.Writer, name string, data interface{}) error {
// 	tmpl := template.Must(t.templates.Clone())
// 	tmpl = template.Must(tmpl.ParseFS(Templates, name))
// 	return tmpl.ExecuteTemplate(w, name, data)
// }

func (t *Template) Render(w io.Writer, contentTemplate string, data interface{}) error {
	tmpl, err := t.templates.Clone()
	if err != nil {
		return err
	}

	// Parse the specific content template
	_, err = tmpl.ParseFS(Templates, contentTemplate)
	if err != nil {
		return err
	}

	return tmpl.ExecuteTemplate(w, "base.tmpl", data)
}
