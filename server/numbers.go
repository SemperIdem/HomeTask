package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"time"
)

func (srv *Server) Numbers() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		ctx, cancel := context.WithTimeout(r.Context(), 7000*time.Millisecond)
		defer cancel()

		query, ok := r.URL.Query()["u"]
		if !ok || len(query[0]) < 1 {
			return
		}

		for _, url := range query {
			srv.ctrl.ProcessUrl(url)
		}

		go func(ctx context.Context, out chan []int) {
			done := 0
			var source, received []int
			fmt.Println("blala")
			for {
				select {
				case <-ctx.Done():
					fmt.Println("timeout")
					close(out)
					return
				case received = <-srv.ctrl.Results:
					source = merge(source, received)
					done++
					fmt.Println("job done")
					if done == len(query) {
						close(out)
						return
					}
				}
			}
		}()

		//for {
		//	select {
		//	case <-ctx.Done():
		//		fmt.Println("timeout")
		//		RespondJSON(w, source)
		//		return
		//	case received = <-srv.result:
		//		source = merge(source, received)
		//		done++
		//		fmt.Println("job done")
		//		if done == len(query) {
		//			RespondJSON(w, source)
		//			return
		//		}
		//	}
		//}
	}
}

func RespondJSON(w http.ResponseWriter, responseObject interface{}) error {
	return RespondWithJSON(w, http.StatusOK, responseObject)
}

func RespondWithJSON(w http.ResponseWriter, statusCode int, responseObject interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	return json.NewEncoder(w).Encode(responseObject)
}

func merge(a []int, b []int) []int {
	sort.Ints(b)

	if len(a) == 0 {
		b = deduplicate(b)
		return b
	}

	i := 0
	j := 0
	var res []int

	for i < len(a) && j < len(b) {
		if a[i] < b[j] {
			if len(res) == 0 || res[len(res)-1] != a[i] {
				res = append(res, a[i])
			}
			i++
		} else {
			res = append(res, b[j])
			if len(res) == 0 || res[len(res)-1] != b[j] {
				res = append(res, b[j])
			}
			j++
		}
	}
	for ; i < len(a); i++ {
		res = append(res, a[i])
	}
	for ; j < len(b); j++ {
		res = append(res, b[j])
	}
	return res
}

func deduplicate(s []int) []int {
	if len(s) < 2 {
		return s
	}

	e := 1
	for i := 1; i < len(s); i++ {
		if s[i] == s[i-1] {
			continue
		}
		s[e] = s[i]
		e++
	}

	return s[:e]
}
