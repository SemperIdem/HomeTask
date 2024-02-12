package controllers

import (
	"context"
	"encoding/json"
	"homeTask/config"
	"io"
	"log"
	"net/http"
)

// NumbersFetcher implements WorkerPool pattern, processing URLs,
// fetching the data and send it back to the number handler
type NumbersFetcher struct {
	urls   chan string
	client HTTPClient

	numWorkers int
	numbJobs   int
	jobs       chan string
	Sender
}

// Sender send retrived nums to Handler
type Sender interface {
	Send(data []int)
	Receive() chan []int
}

//go:generate mockgen -destination=./mocks/mocks_sender.go --build_flags=--mod=mod -package=mocks homeTask/controllers Sender
type SenderNumbers struct {
	results chan []int
}

func (s *SenderNumbers) Send(nums []int) {
	s.results <- nums
}

func (s *SenderNumbers) Receive() chan []int {
	return s.results
}

//go:generate mockgen -destination=./mocks/mocks_httpclient.go --build_flags=--mod=mod -package=mocks homeTask/controllers HTTPClient
type HTTPClient interface {
	Get(url string) (resp *http.Response, err error)
}

// NumbersResponse struct to unmarshall numbers from any json, contains the field
type NumbersResponse struct {
	Numbers []int `json:"numbers"`
}

func New(client HTTPClient, numWorkers, numJobs int) *NumbersFetcher {
	return &NumbersFetcher{
		urls:       make(chan string),
		client:     client,
		numWorkers: numWorkers,
		numbJobs:   numJobs,

		jobs:   make(chan string, numWorkers),
		Sender: &SenderNumbers{make(chan []int, numJobs)},
	}
}

// StartTasks start worker pool
func (t *NumbersFetcher) StartTasks(ctx context.Context) {
	for i := 0; i < t.numWorkers; i++ {
		go t.task(ctx)
	}
}

// task start worker
func (t *NumbersFetcher) task(ctx context.Context) {
	for job := range t.jobs {
		t.FetchNumbers(job, ctx)
	}
}

// StartTasks send urls to worker pool, launch in handler
func (t *NumbersFetcher) ProcessUrls(urls []string) {
	for i := range urls {
		t.jobs <- urls[i]

	}
}

// FetchNumbers request external URL, retrive a number and send to handler
// if there is timeout, http client will finish and cache succesful request, but don't send data to handler this time
func (t *NumbersFetcher) FetchNumbers(url string, ctx context.Context) {

	taskCtx, cancel := context.WithTimeout(ctx, config.DefaultTimeout)
	defer cancel()

	var body []byte
	var err error
	var result NumbersResponse
	resp, err := t.client.Get(url)
	if err == nil {
		defer resp.Body.Close()

		body, err = io.ReadAll(resp.Body)
		if err != nil {
			log.Println("err handle it ,", resp.Status)
		}
		err = json.Unmarshal(body, &result)
		if err != nil {
			log.Println("err handle json")
		}
	}
	select {
	case <-taskCtx.Done():
	default:
		log.Println("sent ", result.Numbers)
		t.Send(result.Numbers)
	}
}
