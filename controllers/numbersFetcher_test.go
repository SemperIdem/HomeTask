package controllers

import (
	"bytes"
	"context"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"homeTask/config"
	"homeTask/controllers/mocks"
	"io"
	"net/http"
	"testing"
	"time"
)

func TestNumbersFetcher_FetchNumbers(t *testing.T) {
	tests := []struct {
		name               string
		tuneMockHTTPClient func(m *mocks.MockHTTPClient)
		tuneMockSender     func(m *mocks.MockSender)
	}{
		{
			name: "succes when you send the data",
			tuneMockHTTPClient: func(m *mocks.MockHTTPClient) {
				m.EXPECT().Get(gomock.Any()).Return(
					&http.Response{
						StatusCode: http.StatusOK,
						Body: io.NopCloser(bytes.NewReader([]byte(`{"numbers": [1,2,3,4,5], 
							"strings": ["one", "two", "three", "four", "five"]}`))),
					}, nil)
			},
			tuneMockSender: func(m *mocks.MockSender) {
				m.EXPECT().Send(gomock.Any()).Times(1)
			},
		},
		{
			name: "emulate timeout, data is not sending",
			tuneMockHTTPClient: func(m *mocks.MockHTTPClient) {
				m.EXPECT().Get(gomock.Any()).Do(func(_ interface{}) { time.Sleep(1 * time.Second) }).Return(
					&http.Response{
						StatusCode: http.StatusOK,
						Body: io.NopCloser(bytes.NewReader([]byte(`{"numbers": [1,2,3,4,5], 
							"strings": ["one", "two", "three", "four", "five"]}`))),
					}, nil)
			},
			tuneMockSender: func(m *mocks.MockSender) {
				m.EXPECT().Send(gomock.Any()).Times(0)
			},
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

		numbersFetcher.FetchNumbers("Test", context.TODO())

	}
}

func TestReadNumbers(t *testing.T) {
	tests := []struct {
		name         string
		expectedNums []int
		response     *http.Response
	}{
		{
			response: &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewReader([]byte(`{"numbers":[1,2,3,4,5], "strings":["one", "two", "three", "four", "five"]}`))),
			},
			expectedNums: []int{1, 2, 3, 4, 5},
		},
	}

	for _, tt := range tests {
		res := ReadNumbers(tt.response)
		require.Equal(t, tt.expectedNums, res)
	}

}
