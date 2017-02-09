package challenge

import (
	"context"
	"testing"
)

// TestUrlTest tests that UrlSanitizer passes only valid urls
func TestUrlTestPassesOnlyValids(t *testing.T) {
	brokenUrls := []string{
		"asdaseq12",
		"",
		"nothing",
		"@:1",
	}
	urls := []string{
		"http://google.com",
		"http://hostnamewithoutsubname",
		"https://check.me",
		"https://check.me?asd=as123",
	}
	urls = append(urls, brokenUrls...)
	urlList := UrlsSanitizer(urls)
	for url := range urlList.SanitizedUrls(context.Background()) {
		for _, brokenUrl := range brokenUrls {
			if url != brokenUrl {
				continue
			}
			t.Fatalf("Broken url passed the UrlsSanitizer validation: %v", url)
		}
	}
}

// TestUrlTestRemovesDoubles that UrlSanitizer removes doubles
func TestUrlTestRemovesDoubles(t *testing.T) {
	urls := []string{
		"http://google.com",
		"http://google.com",
		"http://hostnamewithoutsubname",
		"http://hostnamewithoutsubname",
		"https://check.me",
		"https://check.me",
		"https://check.me?asd=as123",
		"https://check.me?asd=as123",
	}
	urlList := UrlsSanitizer(urls)
	counter := map[string]int{}
	for url := range urlList.SanitizedUrls(context.Background()) {
		counter[url]++
	}
	for _, url := range urls {
		if counter[url] > 1 {
			t.Fatalf("Url %s appers %d time(s)", url, counter[url])
		}
	}
}
