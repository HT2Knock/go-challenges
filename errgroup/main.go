package main

import (
	"context"
	"fmt"
	"net/http"

	"golang.org/x/sync/errgroup"
)

func main() {
	g, ctx := errgroup.WithContext(context.Background())
	urls := []string{
		"http://www.somethingstupidname.com/",
		"http://www.golang.org/",
		"http://www.google.com/",
	}

	for _, url := range urls {
		url := url
		g.Go(func() error {
			req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
			if err != nil {
				return err
			}

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				return fmt.Errorf("failed calling %v: %w", url, err)
			}
			defer resp.Body.Close()

			fmt.Printf("Finished getting url=%v\n", url)
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		fmt.Printf("Finished with error: %v\n", err)
	}
}
