package challenge

import (
	"context"
	"github.com/soider/go-challenge/challenge/tree"
	"testing"
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
