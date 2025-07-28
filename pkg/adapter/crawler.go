package adapter

import (
	"fmt"
	"io"
	"net/http"
	"regexp"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

const (
	httpReg = `https?:\/\/(?:www\.)?[-a-zA-Z0-9@:%._\+~#=]{1,256}\.[a-zA-Z0-9()]{1,6}\b(?:[-a-zA-Z0-9()@:%_\+.~#?&//=]*)`
)

type SimpleWebCrawler struct {
}

func (crawler *SimpleWebCrawler) ExtractLinks(url string) ([]string, error) {
	r, _ := regexp.Compile(httpReg)

	resp, err := http.Get(url)
	if err != nil {
		return []string{}, err
	}
	defer resp.Body.Close()

	contentBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return []string{}, err
	}

	contentStr := string(contentBytes)

	urls := r.FindAllString(contentStr, -1)

	fmt.Println(len(urls))

	return urls, nil
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
