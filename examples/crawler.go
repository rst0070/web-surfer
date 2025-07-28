package main

import (
	"fmt"

	"github.com/rst0070/web-surfer/pkg/adapter"
	"github.com/rst0070/web-surfer/pkg/port"
)

func main() {
	var crawler port.WebCrawler = &adapter.SimpleWebCrawler{}

	links, _ := crawler.ExtractLinks("https://go.dev/")

	for _, link := range links {
		fmt.Println(link)
	}

	metadata, _ := crawler.ExtractMetadata("https://go.dev/")
	fmt.Println(metadata)
}
