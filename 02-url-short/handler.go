package urlshort

import (
	"log"
	"net/http"

	"gopkg.in/yaml.v2"
)

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		for p, u := range pathsToUrls {
			if r.URL.Path == p {
				http.Redirect(w, r, u, http.StatusTemporaryRedirect)
			}
		}
		fallback.ServeHTTP(w, r)
	}
}

// YAMLHandler will parse the provided YAML and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the YAML, then the
// fallback http.Handler will be called instead.
//
// YAML is expected to be in the format:
//
//   - path: /some-path
//     url: https://www.some-url.com/demo
//
// The only errors that can be returned all related to having
// invalid YAML data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.
func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	v := []pathToUrl{}
	err := yaml.Unmarshal(yml, &v)
	if err != nil {
		log.Fatal("yaml handler: ", err)
	}
	m := make(map[string]string, len(v))
	for _, k := range v {
		m[k.Path] = k.Url
	}
	return MapHandler(m, fallback), err
}

type pathToUrl struct {
	Path string `yaml:"path"`
	Url  string `yaml:"url"`
}
