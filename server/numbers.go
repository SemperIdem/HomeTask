package server

import (
	"context"
	"homeTask/config"
	"homeTask/controllers"
	"homeTask/numbers"
	"homeTask/utils"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

// Numbers parses URL's from query paramaters, sends them to the controller
// getting back slices of ints and process them
func Numbers(ctrl NumbersFetcher) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		defer timeTrack(start, "Numbers")
		ctx, cancel := context.WithTimeout(r.Context(), config.DefaultTimeout)
		defer cancel()

		urls, ok := r.URL.Query()["u"]
		if !ok || len(urls[0]) < 1 {
			utils.RespondJSON(w, &controllers.NumbersResponse{Numbers: []int{}})
			return
		}

		ctrl.ProcessUrls(urls)
		done := len(urls)

		source, received := []int{}, []int{}
		for {
			select {
			case <-ctx.Done():
				log.Println("Timeout")
				utils.RespondJSON(w, &controllers.NumbersResponse{Numbers: source})
				return
			case received = <-ctrl.Receive():
				source = numbers.ProcessNumbers(source, received)
				done--
				if done == 0 {
					log.Println("All jobs finished")
					utils.RespondJSON(w, &controllers.NumbersResponse{Numbers: source})
					return
				}
			}
		}
	}
}

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("%s took %s", name, elapsed)
}
