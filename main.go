package main

import (
	"context"
	"homeTask/cache"
	"homeTask/config"
	"homeTask/controllers"
	"homeTask/server"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	log "github.com/sirupsen/logrus"
)

func main() {
	defer log.Println("gracefully shutdown")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	httpCache := cache.New()
	client := &http.Client{
		Transport: &cache.CachedTransport{
			Cache:         httpCache,
			BaseTransport: http.DefaultTransport,
		},
	}

	sender := controllers.NewSenderNumbers(config.NumbJobs)
	t := controllers.New(client, sender, config.NumWorkers)

	out := make(chan []int)

	srv, srvCancel := server.New(ctx, t, out)
	defer srvCancel()

	go t.StartTasks(ctx)

	srv.SetUpRoutes()
	err := srv.ListenAndServeHTTP(":" + config.DefaultPort)
	if err != nil {
		log.Panicf("Error http.ListenAndServe failed: " + err.Error())
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	stopSignal := <-stop
	log.Println("Service received signal: ", stopSignal)
}
