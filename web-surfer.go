package web_surfer

import (
	"github.com/rst0070/web-surfer/pkg/adapter"
	"github.com/rst0070/web-surfer/pkg/types"
	"github.com/rst0070/web-surfer/pkg/usecase"
)

func SurfWebLinks(url string, maxDepth int, maxConcurrency int) []types.WebLink {
	surfer := usecase.BFSSurfer{
		StartUrl:       url,
		MaxDepth:       maxDepth,
		MaxConcurrency: maxConcurrency,
		Crawler:        &adapter.SimpleWebCrawler{},
	}

	return surfer.SurfWeb()
}

func SurfWebLinksStream(url string, maxDepth int, maxConcurrency int) <-chan types.WebLink {
	surfer := usecase.BFSSurfer{
		StartUrl:       url,
		MaxDepth:       maxDepth,
		MaxConcurrency: maxConcurrency,
		Crawler:        &adapter.SimpleWebCrawler{},
	}
	return surfer.SurfWebStream()
}
