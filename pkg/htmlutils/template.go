package htmlutils

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"sync"
)

type TemplateGroup struct {
	mu sync.Mutex
	t  *template.Template

	debug bool
}

func (g *TemplateGroup) SetConfig(debug bool) {
	g.mu.Lock()
	defer g.mu.Unlock()

	g.debug = debug
	g.t = nil
}

func (g *TemplateGroup) Render(w http.ResponseWriter, name string, context interface{}) {
	t := g.GetTemplate(name)
	buf := &bytes.Buffer{}
	err := t.Execute(buf, context)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	w.Header().Set("Content-Length", strconv.Itoa(buf.Len()))
	w.WriteHeader(http.StatusOK)
	buf.WriteTo(w)
}

func (g *TemplateGroup) GetTemplate(name string) *template.Template {
	g.mu.Lock()
	defer g.mu.Unlock()

	if g.t == nil || g.debug {
		g.t = g.LoadTemplates()
	}
	t := g.t.Lookup(name)
	if t == nil {
		panic(fmt.Sprintf("Could not load template %v", name))
	}
	return t
}

// LoadTemplates loads the templates for our toy server
func (g *TemplateGroup) LoadTemplates() *template.Template {
	tData, err := sitedata.LoadFilesInDir("templates")
	if err != nil {
		log.Printf("Error loading template files: %v", err)
	}

	t := template.New("").Funcs(FuncMap())

	for f, fData := range tData {
		log.Printf("Loading template for %v", f)
		_, err := t.New(f).Parse(string(fData))
		if err != nil {
			log.Printf("ERROR: Could parse template %v: %v", f, err)
		}
	}
	return t
}
