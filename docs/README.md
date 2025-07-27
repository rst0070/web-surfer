# Web Surfer - Go Learning Project Development Guide

## Project Overview
A Go library/application that crawls the web starting from one or more seed URLs, following hyperlinks to discover connected pages. This project serves as a practical learning tool for mastering Go programming concepts.

## Goals 🎯

### Primary Goals
- **Learn Go fundamentals**: goroutines, channels, interfaces, error handling, HTTP client
- **Build a functional web crawler**: Start from seed URLs and discover connected pages
- **Implement concurrent processing**: Use Go's concurrency primitives effectively
- **Practice clean architecture**: Separation of concerns, testable code
- **Handle real-world challenges**: Rate limiting, robots.txt compliance, error recovery

### Learning Goals
- Master Go's concurrency model (goroutines + channels)
- Understand HTTP client programming in Go
- Learn HTML parsing and link extraction
- Practice Go testing and benchmarking
- Implement design patterns in Go (worker pools, producer-consumer)
- Handle structured data with Go structs and interfaces

## Non-Goals ❌

### What We Won't Build
- **Full search engine**: No indexing, ranking, or search functionality
- **Distributed crawler**: Single-node operation only
- **JavaScript execution**: Static HTML parsing only
- **Advanced content analysis**: No content classification or NLP
- **Web scraping framework**: Focus on link discovery, not data extraction
- **Production-scale system**: Learning-focused, not enterprise-ready

### Technical Limitations
- No database integration (use in-memory storage or simple files)
- No web UI (CLI interface sufficient)
- No complex configuration management
- No Docker containerization (keep it simple)

## Architecture Overview 🏗️

### Core Components
```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   URL Queue     │◄───│   Crawler       │◄───│   Link Parser   │
│   (Channel)     │    │   (Workers)     │    │   (HTML Parser) │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                       │                       │
         ▼                       ▼                       ▼
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Visited Set   │    │   HTTP Client   │    │   Results       │
│   (Map/Set)     │    │   (Rate Limited)│    │   (Output)      │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

### Key Packages to Learn
- `net/http`: HTTP client operations
- `html`: HTML parsing and tokenization
- `net/url`: URL manipulation and validation
- `sync`: Mutexes, WaitGroups for concurrency
- `context`: Request cancellation and timeouts
- `time`: Rate limiting and delays

## Action Items & Implementation Roadmap 📋

### Phase 1: Foundation (Week 1-2)
#### Setup & Basic Structure
- [ ] **Clean up main.go**: Remove example code
- [ ] **Define core structs**: `WebSurfer`, `CrawlResult`, `CrawlConfig`
- [ ] **Create basic CLI**: Accept seed URLs as command-line arguments
- [ ] **Implement URL validation**: Basic URL parsing and validation
- [ ] **Add logging**: Use standard `log` package for debugging

#### Learning Focus
- Go project structure and packages
- Command-line argument parsing (`flag` package)
- Basic struct definitions and methods

#### Example Structure
```go
type WebSurfer struct {
    config    CrawlConfig
    visited   map[string]bool
    queue     chan string
    results   chan CrawlResult
}

type CrawlConfig struct {
    MaxDepth    int
    MaxPages    int
    Delay       time.Duration
    MaxWorkers  int
}
```

### Phase 2: HTTP & HTML Processing (Week 3-4)
#### Core Crawling Logic
- [ ] **HTTP client setup**: Configure timeouts, user agent
- [ ] **Fetch web pages**: GET requests with error handling
- [ ] **Parse HTML**: Extract links using `golang.org/x/net/html`
- [ ] **URL resolution**: Handle relative URLs correctly
- [ ] **Basic duplicate detection**: Track visited URLs

#### Learning Focus
- HTTP client programming
- HTML parsing and tree traversal
- Error handling patterns in Go
- URL manipulation

#### Key Implementation
```go
func (ws *WebSurfer) fetchPage(url string) (*html.Node, error) {
    resp, err := ws.client.Get(url)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    
    return html.Parse(resp.Body)
}
```

### Phase 3: Concurrency (Week 5-6)
#### Concurrent Processing
- [ ] **Worker pool pattern**: Multiple goroutines processing URLs
- [ ] **Channel communication**: URL queue and result collection
- [ ] **Synchronization**: Proper coordination between workers
- [ ] **Rate limiting**: Respect server resources
- [ ] **Graceful shutdown**: Handle interrupts cleanly

#### Learning Focus
- Goroutines and channels
- Worker pool patterns
- `sync.WaitGroup` and mutexes
- Context cancellation

#### Worker Implementation
```go
func (ws *WebSurfer) worker(ctx context.Context, wg *sync.WaitGroup) {
    defer wg.Done()
    for {
        select {
        case url := <-ws.queue:
            ws.processURL(url)
        case <-ctx.Done():
            return
        }
    }
}
```

### Phase 4: Advanced Features (Week 7-8)
#### Production-Like Features
- [ ] **Robots.txt support**: Basic robots.txt parsing
- [ ] **Depth tracking**: Limit crawl depth
- [ ] **Statistics**: Track pages crawled, errors, etc.
- [ ] **Output formats**: JSON, CSV export options
- [ ] **Configuration**: File-based configuration

#### Learning Focus
- File I/O operations
- JSON marshaling/unmarshaling
- Configuration management
- Testing and benchmarking

### Phase 5: Testing & Polish (Week 9-10)
#### Quality & Testing
- [ ] **Unit tests**: Test individual components
- [ ] **Integration tests**: Test full crawl scenarios
- [ ] **Benchmarks**: Performance testing
- [ ] **Documentation**: Comprehensive README and code docs
- [ ] **Error handling review**: Robust error handling throughout

#### Learning Focus
- Go testing framework
- Benchmark writing
- Documentation best practices
- Code review and refactoring

## Technical Requirements 🔧

### Dependencies (Minimal)
```go
// External dependencies
golang.org/x/net/html  // HTML parsing
golang.org/x/time/rate // Rate limiting (optional)

// Standard library (no external deps needed)
net/http, net/url, html, sync, context, time, etc.
```

### Project Structure
```
web-surfer/
├── main.go              # CLI entry point
├── crawler/
│   ├── surfer.go       # Main WebSurfer struct
│   ├── worker.go       # Worker pool implementation
│   ├── parser.go       # HTML parsing logic
│   └── config.go       # Configuration handling
├── internal/
│   ├── urlutil/        # URL utilities
│   └── storage/        # Visited set management
├── cmd/
│   └── websurfer/      # CLI commands
└── examples/           # Usage examples
```

## Success Metrics 📊

### Functional Goals
- [ ] Successfully crawl a simple website (e.g., local static site)
- [ ] Handle 1000+ pages without crashing
- [ ] Respect rate limits (configurable delay between requests)
- [ ] Export results in structured format
- [ ] Handle common errors gracefully

### Learning Goals
- [ ] Understand and explain Go's concurrency model
- [ ] Write idiomatic Go code following conventions
- [ ] Implement proper error handling patterns
- [ ] Create comprehensive tests
- [ ] Use channels and goroutines effectively

## Getting Started 🚀

### Immediate Next Steps
1. **Clean up `main.go`**: Remove the channel example
2. **Define basic structures**: Start with `WebSurfer` struct
3. **Implement basic CLI**: Parse command-line arguments
4. **Test with simple HTTP request**: Fetch a single page
5. **Add basic HTML parsing**: Extract first link from a page

### Example Starting Point
```go
package main

import (
    "flag"
    "fmt"
    "log"
    "net/http"
    "net/url"
)

func main() {
    var seedURL = flag.String("url", "", "Starting URL for web surfing")
    flag.Parse()
    
    if *seedURL == "" {
        log.Fatal("Please provide a starting URL with -url flag")
    }
    
    // Validate URL
    _, err := url.Parse(*seedURL)
    if err != nil {
        log.Fatalf("Invalid URL: %v", err)
    }
    
    fmt.Printf("Starting web surf from: %s\n", *seedURL)
    
    // TODO: Implement crawling logic
}
```
