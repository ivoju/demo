package main

import (
	"context"
	"fmt"

	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/demo/config"
	"github.com/demo/router"
	"golang.org/x/sync/errgroup"
)

func main() {
	var err error

	conf, logger, err := config.New()
	if err != nil {
		fmt.Println(err)
		return
	}

	logger.Infof("[TOWAMS] Environment %s is ready", conf.Env.Env())

	// create run group
	g, _ := errgroup.WithContext(context.Background())

	// create connection
	var servers []*http.Server

	// goroutine to check for signals to gracefully finish all functions
	g.Go(func() error {
		signalChannel := make(chan os.Signal, 1)
		signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)

		select {
		case sig := <-signalChannel:
			logger.Infof("received signal: %s\n", sig)

			for i, s := range servers {
				if err := s.Shutdown(context.Background()); err != nil {
					if err == nil {
						logger.Infof("error shutting down server %d: %v", i, err)
						return err
					}
				}
			}
			os.Exit(1)
		}
		return nil
	})

	g.Go(func() error { return router.NewGRPCServer(conf, logger) })
	g.Go(func() error { return router.NewHTTPServer(conf, logger) })

	if err := g.Wait(); !router.IgnoreErr(err) {
		logger.Fatal(err)
	}
	logger.Infoln("done.")
}
