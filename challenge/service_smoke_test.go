package challenge

import (
	"context"
	"github.com/soider/go-challenge/challenge/tree"
	"github.com/soider/go-challenge/entities"
	"reflect"
	"testing"
	"time"
)

// TestService tests service
func TestService(t *testing.T) {
	urlList := []string{
		"http://google.com",
		"http://gmail.com",
		"http://yandex.ru",
		"http://yandex.ru",
	}
	ctx := context.Background()
	proccessed := map[string]bool{}
	service := &NumberService{
		sanitizer:          UrlsSanitizer(urlList),
		tick:               make(chan []int),
		done:               make(chan struct{}),
		perTargetResultsCh: make(chan []int),
		data:               tree.New(),
		fetcher: func(ctx context.Context, target string, resultChan chan []int) {
			proccessed[target] = true
			resultChan <- []int{1, 2, 3, 4}
		},
		workers: 1,
	}
	service.Start(ctx)
	<-service.Tick()
	<-service.Tick()
	<-service.Tick()
	<-service.Done()
	for _, url := range urlList {
		if !proccessed[url] {
			t.Fatalf("Url %s wasn't processed", url)
		}
	}
}

// TestServiceTimeout tests service
func TestServiceTimeout(t *testing.T) {
	urlList := []string{
		"http://google.com",
		"http://gmail.com",
		"http://yandex.ru",
		"http://yandex.ru",
	}
	ctx := context.Background()
	ctx, _ = context.WithTimeout(ctx, time.Millisecond*450)
	proccessed := map[string]bool{}
	service := &NumberService{
		sanitizer:          UrlsSanitizer(urlList),
		tick:               make(chan []int),
		done:               make(chan struct{}),
		perTargetResultsCh: make(chan []int),
		data:               tree.New(),
		fetcher: func(ctx context.Context, target string, resultChan chan []int) {
			switch target {
			case urlList[0]:
				time.Sleep(time.Second)
				resultChan <- []int{1, 2}
			case urlList[1]:
				resultChan <- []int{1}
			case urlList[2]:
				resultChan <- []int{6, 5}
			}
			proccessed[target] = true
		},
		workers: 2,
	}
	service.Start(ctx)
	var resp entities.Response
	ticks := 0
	start := time.Now()
	timeoutHappened := false
l:
	for {
		select {
		case i := <-service.Tick():
			resp.Numbers = i
			ticks++
		case <-ctx.Done():
			timeoutHappened = true
			break l
		case <-service.Done():
			break l
		}
	}
	end := time.Now()
	if start.Sub(end) > time.Millisecond*500 {
		t.Fatalf("TImeout more than setted")
	}
	// We need to have results from second and third domain but without first one
	if !reflect.DeepEqual(resp.Numbers, []int{1, 5, 6}) {
		t.Fatalf("Wrong result")
	}
	if !timeoutHappened {
		t.Fatalf("Timeout should happened")
	}
}

// TestServiceNoTimeout tests service
func TestServiceNoTimeout(t *testing.T) {
	urlList := []string{
		"http://google.com",
		"http://gmail.com",
		"http://yandex.ru",
		"http://yandex.ru",
	}
	ctx := context.Background()
	ctx, _ = context.WithTimeout(ctx, time.Millisecond*450)
	proccessed := map[string]bool{}
	service := &NumberService{
		sanitizer:          UrlsSanitizer(urlList),
		tick:               make(chan []int),
		done:               make(chan struct{}),
		perTargetResultsCh: make(chan []int),
		data:               tree.New(),
		fetcher: func(ctx context.Context, target string, resultChan chan []int) {
			switch target {
			case urlList[0]:
				resultChan <- []int{1, 2}
			case urlList[1]:
				resultChan <- []int{1}
			case urlList[2]:
				resultChan <- []int{6, 5}
			}
			proccessed[target] = true
		},
		workers: 2,
	}
	service.Start(ctx)
	var resp entities.Response
	ticks := 0
	start := time.Now()
	timeoutHappened := false
l:
	for {
		select {
		case i := <-service.Tick():
			resp.Numbers = i
			ticks++
		case <-ctx.Done():
			timeoutHappened = true
			break l
		case <-service.Done():
			break l
		}
	}
	end := time.Now()
	if start.Sub(end) > time.Millisecond*500 {
		t.Fatalf("TImeout more than setted")
	}
	// We need to have results from each domain
	if !reflect.DeepEqual(resp.Numbers, []int{1, 2, 5, 6}) {
		t.Fatalf("Wrong result")
	}
	if timeoutHappened {
		t.Fatalf("Timeout should not happened")
	}
}
