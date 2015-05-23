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

func (n HTML) Write(w http.ResponseWriter) error {
	file := n.Name
	ctx := n.Data.(pongo2.Context)

	var t *pongo2.Template

	if tmpl, ok := n.Template[file]; ok {
		t = tmpl
	} else {
		tmpl, err := pongo2.FromCache(file)
		if err != nil {
			return err
		}
		n.Template[file] = tmpl
		t = tmpl
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	return t.ExecuteWriter(ctx, w)
}

type PongRender struct {
	Template map[string]*pongo2.Template
}

func (n *PongRender) Instance(name string, data interface{}) render.Render {
	return HTML{
		Template: n.Template,
		Name:     name,
		Data:     data,
	}
}
func NewPongRender() *PongRender {
	return &PongRender{Template: map[string]*pongo2.Template{}}
}
