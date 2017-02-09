package challenge

import (
	"context"
	"encoding/json"
	"github.com/soider/go-challenge/challenge/entities"
	"github.com/soider/go-challenge/logger"
	"net/http"
)


func Fetch(ctx context.Context, target string, resultChan chan []int) {
	ctx, cancel := context.WithCancel(ctx)
	errCh := make(chan error, 1)

	req, err := http.NewRequest("GET", target, nil)
	req = req.WithContext(ctx) // request would be canceled if deadline exceed
	if err != nil {
		return
	}
	client := &http.Client{}

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
		resultChan <- respObj.Numbers
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
