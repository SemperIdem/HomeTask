package main

import (
	"context"
	"fmt"
	"homeTask/cache"
	"homeTask/controllers"
	"homeTask/server"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
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
	t := controllers.New(client)

	out := make(chan []int)

	srv, srvCancel := server.New(ctx, t, out)
	defer srvCancel()

	go t.FetchData(out)

	srv.SetUpRoutes()
	err := srv.ListenAndServeHTTP(":" + "8080")
	if err != nil {
		log.Panicf("Error http.ListenAndServe failed: " + err.Error())
	}

	fmt.Println("test")

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	stopSignal := <-stop
	fmt.Println("Classifier Container Logs received signal: ", stopSignal)
}
