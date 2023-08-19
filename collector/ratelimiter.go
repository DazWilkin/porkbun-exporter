package collector

import (
	"time"

	"golang.org/x/time/rate"
)

var (
	// Porkbun API has a 1 query/second rate limit (per API key)
	// Porkbun API requests across collectors must wait on the limiter
	apiRateLimiter = rate.NewLimiter(rate.Every(time.Second), 1)
)
