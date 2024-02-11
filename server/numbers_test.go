package server

import (
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"homeTask/config"
	"homeTask/server/mocks"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var ch chan []int

func TestNumbers(t *testing.T) {

	tests := []struct {
		name                   string
		tuneMockNumbersFetcher func(m *mocks.MockNumbersFetcher)
		urlParams              string
		chanOutput             [][]int
		expectedBody           string
	}{
		{
			name:      "succes case with two parameters",
			urlParams: "?u=primes&u=even",
			tuneMockNumbersFetcher: func(m *mocks.MockNumbersFetcher) {
				m.EXPECT().ProcessUrls(gomock.Any()).Times(1)
				m.EXPECT().Receive().AnyTimes().Return(ch)
			},
			chanOutput: [][]int{
				{2, 3, 5, 7, 11, 13},
				{0, 1, 1, 2, 3, 5, 8, 13, 21},
			},
			expectedBody: `{"numbers":[0,1,2,3,5,7,8,11,13,21]}`,
		},
		{
			name:      "only duplicates",
			urlParams: "?u=1&u=2",
			tuneMockNumbersFetcher: func(m *mocks.MockNumbersFetcher) {
				m.EXPECT().ProcessUrls(gomock.Any()).Times(1)
				m.EXPECT().Receive().AnyTimes().Return(ch)
			},
			chanOutput: [][]int{
				{8, 8, 8, 8, 8, 8, 8},
				{8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8},
			},
			expectedBody: `{"numbers":[8]}`,
		},
		{
			name:      "one get empty (could be error)",
			urlParams: "?u=1&u=2",
			tuneMockNumbersFetcher: func(m *mocks.MockNumbersFetcher) {
				m.EXPECT().ProcessUrls(gomock.Any()).Times(1)
				m.EXPECT().Receive().AnyTimes().Return(ch)
			},
			chanOutput: [][]int{
				{8, 7, 5, 4, 3, 1, 2},
				{},
			},
			expectedBody: `{"numbers":[1,2,3,4,5,7,8]}`,
		},
		{
			name: "empty without params",
			tuneMockNumbersFetcher: func(m *mocks.MockNumbersFetcher) {
				m.EXPECT().ProcessUrls(gomock.Any()).Times(0)
				m.EXPECT().Receive().Times(0)
			},
			expectedBody: `{"numbers":[]}`,
		},
	}
	for _, tt := range tests {
		m := gomock.NewController(t)
		nf := mocks.NewMockNumbersFetcher(m)

		rr := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodGet, tt.urlParams, nil)
		if err != nil {
			require.NoError(t, err)
		}

		ch = make(chan []int, config.NumbJobs)

		for _, arr := range tt.chanOutput {
			ch <- arr
		}
		close(ch)

		tt.tuneMockNumbersFetcher(nf)
		handlerFunc := Numbers(nf)
		handlerFunc(rr, req)

		if rr.Result().StatusCode != http.StatusOK {
			require.NoError(t, err)
		}
		b, err := io.ReadAll(rr.Result().Body)
		if err != nil {
			require.NoError(t, err)
		}

		require.Equal(t, tt.expectedBody, strings.TrimSpace(string(b)))
	}

}
