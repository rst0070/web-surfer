package usecase

import (
	"sync"
	"time"

	"github.com/rst0070/web-surfer/pkg/port"
	"github.com/rst0070/web-surfer/pkg/types"
)

type BFSSurfer struct {
	StartUrl       string
	MaxDepth       int
	MaxConcurrency int
	Crawler        port.WebCrawler

	nodeQ       chan types.WebNode
	nextNodeQ   chan types.WebNode
	resultQ     chan types.WebLink
	visitUrlMap map[string]bool
	visitMutex  sync.Mutex

	numWorker      int
	numWorkerMutex sync.Mutex
}

func (surfer *BFSSurfer) SurfWeb() []types.WebLink {

	stream := surfer.SurfWebStream()
	// Result
	var resultList []types.WebLink
	func() {
		for link := range stream {
			resultList = append(resultList, link)
		}
	}()

	return resultList
}

func (surfer *BFSSurfer) SurfWebStream() <-chan types.WebLink {
	// init
	surfer.nodeQ = make(chan types.WebNode)
	surfer.nextNodeQ = make(chan types.WebNode)
	surfer.resultQ = make(chan types.WebLink)

	surfer.visitUrlMap = make(map[string]bool)
	surfer.visitUrlMap[surfer.StartUrl] = true

	surfer.numWorker = 0

	// run concurrently
	for i := 0; i < surfer.MaxConcurrency; i++ {
		go surfer.work()
	}

	// Feed queue. This should be after starting recievers
	surfer.nodeQ <- types.WebNode{
		Url:   surfer.StartUrl,
		Depth: 0,
	}

	// Feed next nodes
	go func() {
		for link := range surfer.nextNodeQ {
			surfer.nodeQ <- link
		}
	}()

	// Monitor deadlock caused by nothing to search (all waiting for nodeQ)
	// This algorithm checks deadlock by watching active workers
	go func() {
		defer close(surfer.resultQ)
		defer close(surfer.nodeQ)

		for {
			time.Sleep(100 * time.Millisecond)
			surfer.numWorkerMutex.Lock()

			if surfer.numWorker == 0 {
				surfer.numWorkerMutex.Unlock()
				return
			} else {
				surfer.numWorkerMutex.Unlock()
			}
		}

	}()

	return surfer.resultQ
}

func (surfer *BFSSurfer) work() {
	for node := range surfer.nodeQ {
		surfer.setWorkActive()

		if node.Depth >= surfer.MaxDepth {
			surfer.setWorkInactive() // Cannot stop whole process: maybe other worker are on same or lower level nodes
			return
		}

		neighbours, err := surfer.Crawler.ExtractLinks(node.Url)

		if err != nil {
			surfer.setWorkInactive()
			continue
		}

		for _, url := range neighbours {

			neighbour := types.WebNode{
				Url:   url,
				Depth: node.Depth + 1,
			}

			// this takes too long time
			// metadata, err := surfer.Crawler.ExtractMetadata(url)
			// if err == nil {
			// 	neighbour.Metadata = metadata
			// }

			surfer.visitMutex.Lock()
			if surfer.visitUrlMap[url] {
				surfer.visitMutex.Unlock()
			} else {
				surfer.visitUrlMap[url] = true
				surfer.visitMutex.Unlock()

				if neighbour.Depth < surfer.MaxDepth {
					surfer.nextNodeQ <- neighbour
					// if enqueue to nodeQ, it can cause situation:
					// all worker tries to enqueue at the same time with waiting forever
				}
			}

			surfer.resultQ <- types.WebLink{
				Source: &node,
				Target: &neighbour,
			}
		}

		surfer.setWorkInactive()
	}
}

func (surfer *BFSSurfer) setWorkActive() {
	surfer.numWorkerMutex.Lock()
	surfer.numWorker += 1
	surfer.numWorkerMutex.Unlock()
}

func (surfer *BFSSurfer) setWorkInactive() {
	surfer.numWorkerMutex.Lock()
	surfer.numWorker -= 1
	surfer.numWorkerMutex.Unlock()
}
