package challenge

import (
	"context"
	"net/url"
)

// UrlsSanitizer used for sanitizing and flattening url slice
type UrlsSanitizer []string

// SanitizedUrls returns channel with sanitized urls, each one would appear only once
func (ul UrlsSanitizer) SanitizedUrls(ctx context.Context) chan string {
	filter := map[string]struct{}{}
	c := make(chan string, len(ul))
	defer close(c)
	for _, urlString := range ul {
		select {
		case <-ctx.Done():
			return c
		default:
			_, err := url.ParseRequestURI(urlString)
			if err != nil {
				continue
			}
			if _, found := filter[urlString]; found {
				continue
			}
			filter[urlString] = struct{}{}
			c <- urlString
		}
	}
	return c
}
