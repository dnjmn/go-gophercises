package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"gophercises/link"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	// "gophercises/link"
)

type urlSet struct {
	Urls []loc `xml:"url"`
}

type loc struct {
	Value string `xml:"loc"`
}

func main() {
	site := flag.String("site", "https://gophercises.com", "initial site for sitemap")
	flag.Parse()

	pages := dfs(*site)

	var toXml urlSet

	for _, v := range pages {
		toXml.Urls = append(toXml.Urls, loc{v})
	}
	en := xml.NewEncoder(os.Stdout)
	if e := en.Encode(toXml); e != nil {
		log.Fatal(e)
	}
}

func dfs(urlStr string) (result []string) {
	seen := make(map[string]struct{})
	nq := []string{urlStr}

	for len(nq) > 0 {
		s := nq[0]
		nq = nq[1:]
		if _, ok := seen[s]; !ok {
			nq = append(nq, get(s)...)
			seen[s] = struct{}{}
		}
	}

	for k := range seen {
		result = append(result, k)
	}

	return result
}

func get(urlStr string) []string {
	rsp, err := http.Get(urlStr)
	if err != nil {
		log.Fatal(err)
	}
	defer rsp.Body.Close()

	reqURL := rsp.Request.URL
	baseURL := url.URL{
		Scheme: reqURL.Scheme,
		Host:   reqURL.Host,
	}
	base := baseURL.String()

	return filter(hrefs(rsp.Body, base), withPrefix(base))
}

func hrefs(r io.Reader, base string) []string {
	links, err := link.Parse(r)
	if err != nil {
		log.Fatal(err)
	}
	var results []string
	for _, link := range links {
		switch {
		case strings.HasPrefix(link.Href, "/"):
			results = append(results, base+link.Href)
		case strings.HasPrefix(link.Href, "http"):
			results = append(results, link.Href)
		}
	}

	return results
}

func filter(links []string, keepFn func(string) bool) (result []string) {
	for _, link := range links {
		if keepFn(link) {
			result = append(result, link)
		}
	}

	return result
}

func withPrefix(pfx string) func(string) bool {
	return func(link string) bool {
		return strings.HasPrefix(link, pfx)
	}
}
