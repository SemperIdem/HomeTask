package controllers

import (
	"encoding/json"
	"fmt"
	"homeTask/cache"
	"io"
	"net/http"
)

type TaskManager struct {
	urls   chan string
	cache  cache.MemoryCache
	client *http.Client

	numWorkers int
	numbJobs   int
	Jobs       chan string
	Results    chan []int
}

func New(client *http.Client) *TaskManager {
	return &TaskManager{
		urls:       make(chan string),
		client:     client,
		numWorkers: 10,
		numbJobs:   30,

		Jobs:    make(chan string, 30),
		Results: make(chan []int, 30),
	}
}

func (t *TaskManager) FetchData(out chan []int) {
	for i := 0; i < t.numWorkers; i++ {
		go t.task()
	}
}

func (t *TaskManager) ProcessUrl(url string) {
	t.Jobs <- url
}

type NumbersResponse struct {
	Numbers []int `json:"numbers"`
}

func (t *TaskManager) task() {
	for job := range t.Jobs {
		numbers := t.fetchNumbers(job)
		t.Results <- numbers
	}
}

func (t *TaskManager) fetchNumbers(url string) []int {
	var body []byte
	var err error
	var result NumbersResponse
	resp, _ := t.client.Get(url)
	defer resp.Body.Close()
	body, err = io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("err handle it")
		return []int{}
	}
	err = json.Unmarshal(body, &result)
	if err != nil {
		fmt.Println("err handle json")
		return []int{}
	}
	return result.Numbers

}
