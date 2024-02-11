package server

import (
	"context"
	"homeTask/config"
	"homeTask/controllers"
	"homeTask/utils"
	"log"
	"net/http"
)

func Numbers(ctrl NumbersFetcher) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), config.DefaultTimeout)
		defer cancel()

		urls, ok := r.URL.Query()["u"]
		if !ok || len(urls[0]) < 1 {
			utils.RespondJSON(w, &controllers.NumbersResponse{Numbers: []int{}})
			return
		}

		ctrl.ProcessUrls(urls)
		done := len(urls)
		var source, received []int
		for {
			select {
			case <-ctx.Done():
				log.Println("Timeout")
				utils.RespondJSON(w, &controllers.NumbersResponse{Numbers: source})
				return
			case received = <-ctrl.Receive():
				source = utils.ProcessNumbers(source, received)
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

//func (srv *Server) HandlerValues() http.HandlerFunc {
//	return func(w http.ResponseWriter, r *http.Request) {
//		ctx, cancel := context.WithTimeout(r.Context(), config.DefaultTimeout)
//		defer cancel()
//
//		vals, ok := r.URL.Query()["vals"]
//		if !ok || len(urls[0]) < 1 {
//			utils.RespondJSON(w, &ValueResponse{Numbers: []int{}})
//		}
//
//		srv.ctrl.ProcessValues(vals)
//		ready := len(urls)
//		var old, new []int
//		for {
//			select {
//			case <-ctx.Done():
//				log.Println("Timeout")
//				utils.RespondJSON(w, &ValueResponse{Numbers: old})
//				return
//			case received = <-srv.ctrl.ReadValuesExternal():
//				old = utils.ProcessNumbers(old, new)
//				ready--
//				if ready == 0 {
//					log.Println("All jobs finished")
//					utils.RespondJSON(w, &ValueResponse{Numbers: old})
//					return
//				}
//			}
//		}
//	}
//}
