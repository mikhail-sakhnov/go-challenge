package challenge

import (
	"context"
	"github.com/soider/go-challenge/challenge/tree"
	"github.com/soider/go-challenge/logger"
	"sync"
)

type FetchFunction func(ctx context.Context, target string, resultChan chan []int)

type NumberService struct {
	tick               chan []int
	done               chan struct{}
	perTargetResultsCh chan []int
	sanitizer          UrlsSanitizer
	fetcher            FetchFunction
	data               *tree.Tree
	wg                 sync.WaitGroup
}

func NewService(sanitizer UrlsSanitizer) *NumberService {
	return &NumberService{
		sanitizer:          sanitizer,
		tick:               make(chan []int),
		done:               make(chan struct{}),
		perTargetResultsCh: make(chan []int),
		data:               tree.New(),
		fetcher:            Fetch,
	}
}

func (ns *NumberService) Start(ctx context.Context) {
	urlsCh := ns.sanitizer.SanitizedUrls(ctx)
	go ns.responseCollector(ctx)
	for target := range urlsCh {
		ns.wg.Add(1)
		go ns.fetch(ctx, target)
	}
	go ns.waiter(ctx)
}

func (ns *NumberService) waiter(ctx context.Context) {
	ns.wg.Wait()
	close(ns.done)
}

func (ns *NumberService) fetch(ctx context.Context, target string) {
	defer ns.wg.Done()
	ns.fetcher(ctx, target, ns.perTargetResultsCh)
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

func (ns *NumberService) Tick() chan []int {
	return ns.tick
}
func (ns *NumberService) Done() chan struct{} {
	return ns.done
}
