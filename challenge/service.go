package challenge

import (
	"context"
	"encoding/json"
	"github.com/soider/go-challenge/challenge/tree"
	"github.com/soider/go-challenge/logger"
	"net/http"
	"sort"
	"sync"
	"github.com/soider/go-challenge/challenge/entities"
)

func sortedInsert(slice *[]int, number int) {
	s := *slice
	if len(s) == 0 {
		s = []int{number}
		*slice = s
		return
	}
	position := sort.SearchInts(*slice, number)
	s = append(s, 0)
	copy(s[position+1:], s[position:])
	s[position] = number
	*slice = s
}

type NumberService struct {
	tick               chan []int
	done               chan struct{}
	state              []int
	numbersFromTargets chan []int
	sanitizer          UrlsSanitizer
	stateTree          *tree.Tree
	wg                 sync.WaitGroup
}

func NewService(sanitizer UrlsSanitizer) *NumberService {
	return &NumberService{
		sanitizer:          sanitizer,
		tick:               make(chan []int),
		state:              []int{},
		done:               make(chan struct{}),
		numbersFromTargets: make(chan []int),
		stateTree:          tree.New(),
	}
}

func (ns *NumberService) Start(ctx context.Context) {
	urlsCh := ns.sanitizer.SanitizedUrls(ctx)
	tasksCount := len(urlsCh)
	go ns.responseCollector(ctx, tasksCount)
	for target := range urlsCh {
		go ns.fetch(ctx, target)
	}
}

func (ns *NumberService) fetch(ctx context.Context, target string) {
	ctx, cancel := context.WithCancel(ctx)
	req, err := http.NewRequest("GET", target, nil)
	req = req.WithContext(ctx) // request would be canceled if deadline exceed
	if err != nil {
		return
	}
	client := &http.Client{}
	errCh := make(chan error, 1)
	go func() {
		var respObj entities.Response
		resp, err := client.Do(req)
		if err != nil {
			errCh <- err
			return
		}
		if resp.StatusCode != http.StatusOK {
			return
		}
		if err := json.NewDecoder(resp.Body).Decode(&respObj); err != nil {
			errCh <- err
			return
		}
		ns.numbersFromTargets <- respObj.Numbers
		errCh <- nil
	}()
	select {
	case <-ctx.Done():
		cancel() // to stop already sent requests
		return
	case err := <-errCh:
		if err != nil {
			logger.FromContext(ctx).Printf("Have error while fetching %v", err)
		}
		return
	}
}
func (ns *NumberService) responseCollector(ctx context.Context, tasksCount int) {
	part := 0
	for {
		select {
		case <-ctx.Done():
			return
		case newPart := <-ns.numbersFromTargets:
			part++
			logger.FromContext(ctx).Printf("Got numbers %v", newPart)
			ns.handleNewPart(newPart)
			if part == tasksCount {
				close(ns.done)
			}
		}
	}
}

func (ns *NumberService) handleNewPart(newPart []int) {
	for _, number := range newPart {
		ns.stateTree.Insert(number)
	}
	ns.tick <- ns.stateTree.ToSlice()
}

func (ns *NumberService) Tick() chan []int {
	return ns.tick
}
func (ns *NumberService) Done() chan struct{} {
	return ns.done
}
