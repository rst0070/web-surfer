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

	waitGroup sync.WaitGroup
}

func (surfer *BFSSurfer) SurfWeb() []types.WebLink {

	// init
	surfer.nodeQ = make(chan types.WebNode)
	surfer.nextNodeQ = make(chan types.WebNode)
	surfer.resultQ = make(chan types.WebLink)

	surfer.visitUrlMap = make(map[string]bool)
	surfer.visitUrlMap[surfer.StartUrl] = true

	surfer.numWorker = 0

	surfer.waitGroup.Add(1)

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

	// Result
	var resultList []types.WebLink
	func() {
		for link := range surfer.resultQ {
			resultList = append(resultList, link)
		}
	}()

	return resultList
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
			surfer.visitMutex.Lock()
			if surfer.visitUrlMap[url] {
				surfer.visitMutex.Unlock()
				continue
			} else {
				surfer.visitUrlMap[url] = true
				surfer.visitMutex.Unlock()
			}

			neighbour := types.WebNode{
				Url:   url,
				Depth: node.Depth + 1,
			}

			if node.Depth < surfer.MaxDepth-1 {
				surfer.nextNodeQ <- neighbour
				// if enqueue to nodeQ, it can cause situation:
				// all worker tries to enqueue at the same time with waiting forever
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
