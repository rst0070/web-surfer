package main

import (
	"fmt"

	web_surfer "github.com/rst0070/web-surfer"
)

func main() {
	links := web_surfer.SurfWebLinks(
		"https://go.dev/",
		3,
		100,
	)

	for _, link := range links {
		fmt.Println(*(link.Source), *(link.Target))
	}
}
