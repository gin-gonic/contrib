package multitemplate

import (
	"html/template"

	"github.com/gin-gonic/gin/render"
)

type htmlRender struct {
	templates map[string]*template.Template
}

var _ render.HTMLRender = &htmlRender{}

func New() *htmlRender {
	return &htmlRender{
		templates: make(map[string]*template.Template),
	}
}

func (r *htmlRender) Add(name string, tmpl *template.Template) {
	if tmpl == nil {
		panic("template can not be nil")
	}
	if len(name) == 0 {
		panic("template name cannot be empty")
	}
	r.templates[name] = tmpl
}

func (r *htmlRender) AddFromFiles(name string, files ...string) *template.Template {
	tmpl := template.Must(template.ParseFiles(files...))
	r.Add(name, tmpl)
	return tmpl
}

func (r *htmlRender) AddFromGlob(name, glob string) *template.Template {
	tmpl := template.Must(template.ParseGlob(glob))
	r.Add(name, tmpl)
	return tmpl
}

func (r *htmlRender) AddFromString(name, templateString string) *template.Template {
	tmpl := template.Must(template.New("").Parse(templateString))
	r.Add(name, tmpl)
	return tmpl
}

func (r *htmlRender) Instance(name string, data interface{}) render.Render {
	return render.HTML{
		Template: r.templates[name],
		Data:     data,
	}
}
