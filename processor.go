package processor

import (
	"fmt"
	"math/rand"
	"time"

	. "bouncer/models"
)

func ProcessRequest(req Request) {
	fmt.Printf("Processing request ID: %s from IP: %s\n", req.ID, req.IPAddress)
	// Simulate processing time
	// random processing time between 1-3 seconds
	processTime := time.Duration(1+rand.Intn(3)) * time.Second
	time.Sleep(processTime)
	fmt.Printf("Finished processing request ID: %s\n", req.ID)
}

func StartWorkerPool(queue *MiddlewareQueue) {
	for i := 0; i < queue.workerCount; i++ {
		go func(workerID int) {
			for req := range queue.mainQueue {
				fmt.Printf("Worker %d picked up request ID: %s\n", workerID, req.ID)
				ProcessRequest(req)
				queue.mu.Lock()
				queue.ipCounts[req.IPAddress]--
				queue.mu.Unlock()
			}
		}(i)
	}
}

// process request go routine
func ProcessRequestRoutine(queue *MiddlewareQueue, req Request) {
	queue.mu.Lock()
	if queue.ipCounts[req.IPAddress] < queue.ipLimit {
		queue.ipCounts[req.IPAddress]++
		queue.mu.Unlock()
		queue.mainQueue <- req
	} else {
		queue.waitQueue = append(queue.waitQueue, req)
		queue.mu.Unlock()
		fmt.Printf("Request ID: %s from IP: %s added to wait queue\n", req.ID, req.IPAddress)
	}
}
