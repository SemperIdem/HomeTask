package controllers

import (
	"bytes"
	"context"
	"github.com/golang/mock/gomock"
	"homeTask/config"
	"homeTask/controllers/mocks"
	"io"
	"net/http"
	"testing"
)

func TestNumbersFetcher_FetchNumbers(t *testing.T) {
	tests := []struct {
		name               string
		tuneMockHTTPClient func(m *mocks.MockHTTPClient)
		tuneMockSender     func(m *mocks.MockSender)
		expectedNums       []int
		isTimeOut          bool
	}{
		{
			name: "default",
			tuneMockHTTPClient: func(m *mocks.MockHTTPClient) {
				m.EXPECT().Get(gomock.Any()).Return(
					&http.Response{
						StatusCode: http.StatusOK,
						Body: io.NopCloser(bytes.NewReader([]byte(`{"numbers": [1,2,3,4,5], 
							"strings": ["one", "two", "three", "four", "five"}}`))),
					}, nil)
			},
			tuneMockSender: func(m *mocks.MockSender) {
				m.EXPECT().Send(gomock.Any()).Times(1)
			},
			expectedNums: []int{1, 2, 3, 4, 5},
		},
		{
			name: "timeout",
			tuneMockHTTPClient: func(m *mocks.MockHTTPClient) {
				m.EXPECT().Get(gomock.Any()).Return(
					&http.Response{
						StatusCode: http.StatusOK,
						Body: io.NopCloser(bytes.NewReader([]byte(`{"numbers": [1,2,3,4,5], 
							"strings": ["one", "two", "three", "four", "five"}}`))),
					}, nil)
			},
			tuneMockSender: func(m *mocks.MockSender) {
				m.EXPECT().Send(gomock.Any()).Times(1)
			},
			expectedNums: []int{1, 2, 3, 4, 5},
		},
	}

	for _, tt := range tests {
		m := gomock.NewController(t)
		httpClient := mocks.NewMockHTTPClient(m)
		sender := mocks.NewMockSender(m)

		tt.tuneMockSender(sender)
		tt.tuneMockHTTPClient(httpClient)

		numbersFetcher := New(httpClient, config.NumWorkers, config.NumbJobs)
		numbersFetcher.Sender = sender

		numbersFetcher.FetchNumbers("Test", context.Background())

	}
}
