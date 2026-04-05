package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"
)

type config struct {
	host string
	port string
}

func main() {
	ctx := context.Background()
	if err := run(ctx, os.Stdout, os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "%s", err)
		os.Exit(1)
	}
}

func run(ctx context.Context, w io.Writer, args []string) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	var cfg config

	fs := flag.NewFlagSet(args[0], flag.ContinueOnError)

	fs.StringVar(&cfg.host, "host", "localhost", "server host")
	fs.StringVar(&cfg.port, "port", "3000", "server port")

	if err := fs.Parse(os.Args[1:]); err != nil {
		return err
	}

	srv := NewServer()

	httpServer := &http.Server{
		Addr:    net.JoinHostPort(cfg.host, cfg.port),
		Handler: srv,
	}

	go func() {
		log.Printf("Listening on port %v \n", cfg.port)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Fprintf(os.Stderr, "error listening and serving: %s \n", err)
		}
	}()

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		<-ctx.Done()

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			fmt.Fprintf(os.Stderr, "error shutting down server: %s \n", err)
		}
	}()

	wg.Wait()
	return nil
}
