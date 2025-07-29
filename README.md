# Web Surfer
A Go library/application that crawls the web starting from seed URL, following hyperlinks to discover connected pages.  
  
## Example
```go
package main

import (
	"fmt"

	web_surfer "github.com/rst0070/web-surfer"
)

func main() {
	stream := web_surfer.SurfWebLinksStream("https://go.dev/", 3, 100)
	// start url, maximum depth, max concurrency

	for link := range stream {
		fmt.Println(*(link.Source), *(link.Target))
	}
}

```