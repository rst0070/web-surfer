package main

import (
	"fmt"

	web_surfer "github.com/rst0070/web-surfer"
)

func main() {
	stream := web_surfer.SurfWebLinksStream("https://go.dev/", 3, 100)

	for link := range stream {
		fmt.Println(*(link.Source), *(link.Target))
	}
}
