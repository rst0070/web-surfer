package pkg

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/rst0070/web-surfer/pkg/service"
	"github.com/rst0070/web-surfer/pkg/types"
)

func SurfWebLinks(url string, maxDepth int, maxConcurrency int) []types.WebLink {
	workCtx, cancelWork := context.WithCancel(context.Background())
	defer cancelWork()

	searchQueue := make(chan types.WebNode)
	resultQueue := make(chan types.WebLink)

	waitGroup := sync.WaitGroup{}
	visitMutex := sync.Mutex{}
	visitMap := make(map[string]bool)
	visitMap[url] = true

	go func() {
		noActivityCount := 0
		for {
			select {
			case <-workCtx.Done():
				return
			case <-time.After(100 * time.Millisecond):
				if len(searchQueue) == 0 {
					noActivityCount++
					if noActivityCount > 20 { // 2 seconds of no activity
						fmt.Println("No activity detected - completing search")
						cancelWork()
						return
					}
				} else {
					noActivityCount = 0
				}
			}
		}
	}()

	for i := 0; i < maxConcurrency; i++ {
		go service.ConcurrentBFS(workCtx, cancelWork, &waitGroup, maxDepth, &visitMutex, visitMap, searchQueue, resultQueue)
	}

	searchQueue <- types.WebNode{
		Url:   url,
		Depth: 0,
	}

	resultIdx := -1
	resultSize := 99
	resultList := make([]types.WebLink, resultSize)

	go func() {
		for {
			select {
			case <-workCtx.Done():
				return
			case link, ok := <-resultQueue:
				if !ok {
					return
				}
				resultIdx++
				if resultIdx < resultSize {
					resultList[resultIdx] = link
				} else {
					resultList = append(resultList, link)
					resultSize++
				}

			}
		}
	}()

	waitGroup.Wait()
	cancelWork()
	close(searchQueue)
	close(resultQueue)

	return resultList
}

func SurfWebLinksStream(url string, maxDepth int, maxConcurrency int, resultChannel chan types.WebLink) {

}
