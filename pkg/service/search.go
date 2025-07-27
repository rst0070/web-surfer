package service

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/rst0070/web-surfer/pkg/types"
)

func ConcurrentBFS(
	workCtx context.Context,
	cancelWork context.CancelFunc,
	waitGroup *sync.WaitGroup,
	maxDepth int,
	visitLock *sync.Mutex,
	visitMap map[string]bool,
	searchQueue chan types.WebNode,
	resultQueue chan types.WebLink,
) {
	waitGroup.Add(1)
	defer waitGroup.Done()

	for {
		select {
		case <-workCtx.Done():
			return
		case node, ok := <-searchQueue:
			if !ok {
				return
			}

			if node.Depth >= maxDepth {
				cancelWork()
				return
			}

			neighbours, err := CrawlWebLinks(node.Url)

			if err != nil {
				continue
			}

			for _, url := range neighbours {
				visitLock.Lock()
				if visitMap[url] {
					visitLock.Unlock()
					continue
				} else {
					visitMap[url] = true
					visitLock.Unlock()
				}

				neighbour := types.WebNode{
					Url:   url,
					Depth: node.Depth + 1,
				}

				select {
				case <-workCtx.Done():
					return
				case searchQueue <- neighbour:
				}

				select {
				case <-workCtx.Done():
					return
				case resultQueue <- types.WebLink{
					Source: &node,
					Target: &neighbour,
				}:
				}
			}
		case <-time.After(2 * time.Second):
			fmt.Println("Timeout on search queue")
			return
		}
	}

}
