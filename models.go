package models

import (
	"sync"
	"time"
)

type Request struct {
	ID        string
	IPAddress string
	Payload   string
	CreatedAt time.Time
}

type MiddlewareQueue struct {
	mainQueue   chan Request   // Go's native fixed-size queue
	waitQueue   []Request      // Slice for the wait queue (allows easy sorting/shuffling for priority)
	ipCounts    map[string]int // Fast O(1) lookup for address limits
	ipLimit     int            // Max requests allowed per IP in flight
	mu          sync.Mutex     // Protects the waitQueue and ipCounts map
	workerCount int
}
