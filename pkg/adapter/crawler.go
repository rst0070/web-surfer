package adapter

import (
	"net/http"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

type SimpleWebCrawler struct {
}

func (crawler *SimpleWebCrawler) ExtractLinks(url string) ([]string, error) {
	links := []string{}

	resp, err := http.Get(url)
	if err != nil {
		return links, err
	}
	defer resp.Body.Close()

	doc, err := html.Parse(resp.Body)
	if err != nil {
		return links, err
	}

	for n := range doc.Descendants() {
		if n.Type == html.ElementNode && n.DataAtom == atom.A {
			for _, attr := range n.Attr {
				if attr.Key == "href" {
					links = append(links, attr.Val)
				}
			}
		}
	}

	return links, nil
}

func (crawler *SimpleWebCrawler) ExtractMetadata(url string) (map[string]string, error) {
	result := make(map[string]string)
	resp, err := http.Get(url)
	if err != nil {
		return result, err
	}
	defer resp.Body.Close()

	doc, err := html.Parse(resp.Body)
	if err != nil {
		return result, err
	}

	for n := range doc.Descendants() {
		if n.Type == html.ElementNode && n.DataAtom == atom.Meta {
			var metaName string
			var metaContent string

			for _, meta := range n.Attr {
				switch key := meta.Key; key {
				case "name":
					metaName = meta.Val
				case "content":
					metaContent = meta.Val
				}
			}

			if len(metaName) > 0 && len(metaContent) > 0 {
				result[metaName] = metaContent
			}
		}
	}

	return result, nil
}
