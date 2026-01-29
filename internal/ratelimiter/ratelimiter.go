package ratelimiter

import "time"

type Limiter interface {
	// The duration is how to wait before can make another request
	// if not allowed return 429 -> header with duration before next request
	Allow(ip string) (bool, time.Duration)
}

type Config struct {
	RequestsPerTimeFrame int
	TimeFrame            time.Duration
	Enabled              bool
}
