# Bouncer

A request queue middleware written in Go. It rate-limits and queues incoming requests using a two-tier queue system — a **main queue** (fixed-size, channel-backed) and a **wait queue** (dynamic slice) for overflow and re-prioritization.

## Current Status

Early development. The core data structures and a basic worker pool are in place. The main queue is a Go channel (fixed capacity), and the wait queue is a slice that collects requests that exceed the per-IP limit or couldn't enter the main queue.

## Architecture

```
Incoming Requests
       │
       ▼
┌──────────────────┐
│  Per-IP Throttle  │  ← drops / wait-queues if too many from same address
└────────┬─────────┘
         │ (accepted)
         ▼
┌──────────────────┐
│   Main Queue      │  ← fixed-size channel (FIFO)
│   (chan Request)  │
└────────┬─────────┘
         │ (popped by worker)
         ▼
┌──────────────────┐
│   Worker Pool     │  ← goroutines that simulate processing
└──────────────────┘
         │
         ▼
   Wait Queue (slice)  ← overflow / IP-limited requests
```

### Components

| File | Purpose |
|---|---|
| [`models.go`](models.go) | `Request` struct and `MiddlewareQueue` type (main queue, wait queue, IP counters, mutex). |
| [`processor.go`](processor.go) | `ProcessRequest` (simulated work), `StartWorkerPool` (goroutine workers), `ProcessRequestRoutine` (enqueue logic with IP check). |
| [`main.go`](main.go) | Entry point — currently just prints "starting bouncer...". |
| [`primary_queue.go`](primary_queue.go) | (Placeholder) Main queue operations. |
| [`wait_queue.go`](wait_queue.go) | (Placeholder) Wait queue operations. |

## Planned Features

- **Wait-queue promotion** — move requests from the wait queue into the main queue when it has capacity, with priority for long-waiting requests.
- **Dynamic worker scaling** — spawn / kill worker goroutines based on queue pressure.
- **Concurrent simulation** — multiple simulated clients generating requests to test throughput.
- **Configurable limits** — IP limit, queue sizes, timeouts via flags or config file.

## Running

```bash
go run .
```
