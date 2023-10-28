package link

import (
	"io"
	"strings"

	"golang.org/x/net/html"
)

type Link struct {
	Href string
	Text string
}

// Parse will accept an HTML doc and return []Link, error
func Parse(r io.Reader) ([]Link, error) {
	node, err := html.Parse(r)
	if err != nil {
		return nil, err
	}

	nodes := linkNodes(node)

	var links []Link
	for _, node = range nodes {
		links = append(links, linkFromNode(node))
	}
	return links, nil
}

func linkFromNode(node *html.Node) (l Link) {
	for _, attr := range node.Attr {
		if attr.Key == "href" {
			l.Href = attr.Val
			break
		}
	}
	l.Text = strings.Join(strings.Fields(textFromNode(node)), " ")
	return
}

func textFromNode(node *html.Node) string {
	if node.Type == html.TextNode {
		return node.Data
	}
	if node.Type != html.ElementNode {
		return ""
	}

	var ret string
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		ret += textFromNode(c)
	}
	return ret
}

func linkNodes(node *html.Node) (ret []*html.Node) {
	if node.Type == html.ElementNode && node.Data == "a" {
		return []*html.Node{node}
	}

	for c := node.FirstChild; c != nil; c = c.NextSibling {
		ret = append(ret, linkNodes(c)...)
	}
	return
}
