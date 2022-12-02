package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {
	mux := http.NewServeMux()
	mux.Handle(
		"/story/",
		NewHandler(
			parseJSON(),
			WithTemplate(tmp),           // with options
			WithPathFn(defaultPathFunc), // with options
		),
	)
	log.Fatal(http.ListenAndServe(":8080", mux))
}

var tmp *template.Template

func init() {
	b, err := os.ReadFile("default.html")
	if err != nil {
		log.Fatal(err)
	}
	tmp = template.Must(template.New("").Parse(string(b)))
}

type Story map[string]Chapter

type Chapter struct {
	Title      string   `json:"title"`
	Paragraphs []string `json:"story"`
	Options    []struct {
		Text    string `json:"text"`
		Chapter string `json:"arc"`
	} `json:"options"`
}

func parseJSON() Story {
	d, e := os.ReadFile("gopher.json")
	if e != nil {
		log.Fatal("parseJSON: ", e)
	}

	var v Story
	e = json.Unmarshal(d, &v)
	if e != nil {
		log.Fatal("parseJSON: ", e)
	}

	return v
}

type handler struct {
	story  Story
	tmp    *template.Template
	pathFn func(r *http.Request) string
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := h.pathFn(r)
	title := path[1:]
	if _, ok := h.story[title]; ok {
		err := h.tmp.Execute(w, h.story[title])
		if err != nil {
			log.Fatal(err)
		}
		return
	}
	http.Error(w, "chapter not found", 404)
}

func NewHandler(s Story, opts ...HandlerOptions) http.Handler {
	h := handler{story: s, tmp: tmp, pathFn: defaultPathFunc}
	for _, opt := range opts {
		opt(&h)
	}
	return h
}

func defaultPathFunc(r *http.Request) string {
	path := strings.TrimSpace(r.URL.Path)
	if path == "/story/" || path == "/story" {
		path = "/story/intro"
	}
	return path[len("/story"):]
}

type HandlerOptions func(h *handler)

func WithTemplate(t *template.Template) HandlerOptions {
	return func(h *handler) {
		h.tmp = t
	}
}

func WithPathFn(fn func(r *http.Request) string) HandlerOptions {
	return func(h *handler) {
		h.pathFn = fn
	}
}
