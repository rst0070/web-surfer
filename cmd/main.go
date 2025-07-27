package main

import (
	"fmt"

	"github.com/rst0070/web-surfer/pkg"
)

func main() {
	links := pkg.SurfWebLinks("https://gobyexample.com/logging", 10, 100)

	for _, link := range links {
		fmt.Println(*(link.Source), *(link.Target))
	}
}
