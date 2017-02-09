package challenge

import (
	"context"
	"github.com/soider/go-challenge/challenge/tree"
	"github.com/soider/go-challenge/logger"
	"sync"
)

// FetchFunction type for fetching functions
type FetchFunction func(ctx context.Context, target string, resultChan chan []int)

// NumberService fetchs remote numbers. merge them and sorts
type NumberService struct {
	tick               chan []int
	done               chan struct{}
	perTargetResultsCh chan []int
	sanitizer          UrlsSanitizer
	fetcher            FetchFunction
	data               *tree.Tree
	wg                 sync.WaitGroup
	workers            int
}

// NewService constructor
func NewService(sanitizer UrlsSanitizer, workers int) *NumberService {
	return &NumberService{
		sanitizer:          sanitizer,
		tick:               make(chan []int),
		done:               make(chan struct{}),
		perTargetResultsCh: make(chan []int),
		data:               tree.New(),
		fetcher:            Fetch,
		workers:            workers,
	}
}

// Start starts the process
func (ns *NumberService) Start(ctx context.Context) {
	urlsCh := ns.sanitizer.SanitizedUrls(ctx)
	go ns.responseCollector(ctx)
	ns.wg.Add(ns.workers)
	for i := 0; i < ns.workers; i++ {
		go func() {
			ns.fetch(ctx, urlsCh)
			ns.wg.Done()
		}()
	}
	go ns.waiter(ctx)
}

func (ns *NumberService) waiter(ctx context.Context) {
	ns.wg.Wait()
	close(ns.done)
}

func (ns *NumberService) fetch(ctx context.Context, ch chan string) {
	for target := range ch {
		ns.fetcher(ctx, target, ns.perTargetResultsCh)
	}
}

func (ns *NumberService) responseCollector(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case newPart := <-ns.perTargetResultsCh:
			logger.FromContext(ctx).Printf("Got numbers %v", newPart)
			ns.handleNewPart(newPart)
		}
	}
}

func (ns *NumberService) handleNewPart(newPart []int) {
	for _, number := range newPart {
		ns.data.Insert(number)
	}
	ns.tick <- ns.data.ToSlice()
}

// Tick returns tick channel
func (ns *NumberService) Tick() chan []int {
	return ns.tick
}

// Done returns done channel
func (ns *NumberService) Done() chan struct{} {
	return ns.done
}
