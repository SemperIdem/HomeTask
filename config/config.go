package config

import "time"

const (
	DefaultPort     = "8080"
	DefaultCacheTTL = 5 * time.Minute
	DefaultTimeout  = 500 * time.Millisecond
	NumWorkers      = 64
	NumbJobs        = 128
)
