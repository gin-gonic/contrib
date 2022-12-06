package pongo2

import (
	"github.com/flosch/pongo2"
	"github.com/gin-gonic/gin/render"
	"net/http"
)

type HTML struct {
	Template map[string]*pongo2.Template
	Name     string
	Data     interface{}
}

func (h HTML) Write(w http.ResponseWriter) error {
	file := h.Name
	ctx := h.Data.(pongo2.Context)

	var t *pongo2.Template

	if tmpl, ok := h.Template[file]; ok {
		t = tmpl
	} else {
		tmpl, err := pongo2.FromCache(file)
		if err != nil {
			return err
		}
		h.Template[file] = tmpl
		t = tmpl
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	return t.ExecuteWriter(ctx, w)
}

type PongRender struct {
	Template map[string]*pongo2.Template
}

func (p *PongRender) Instance(name string, data interface{}) render.Render {
	return HTML{
		Template: p.Template,
		Name:     name,
		Data:     data,
	}
}
func NewPongRender() *PongRender {
	return &PongRender{Template: map[string]*pongo2.Template{}}
}
