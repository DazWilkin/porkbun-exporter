package collector

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/nrdcg/porkbun"
	"github.com/prometheus/client_golang/prometheus"
	"golang.org/x/time/rate"
)

// SSLCollector is a Promethues collector that represents the Porkbun DNS functionality
type SSLCollector struct {
	client         *porkbun.Client
	sslRateLimiter *rate.Limiter

	// Configuration
	domains []string

	// Metrics
	Bundle *prometheus.Desc
}

// NewSSLCollector is a function that creates a new SSLCollector
func NewSSLCollector(apikey, secret string, domains []string) *SSLCollector {
	client := porkbun.New(secret, apikey)

	// Porkbun API /ssl endpoint has a 1 qps rate limit (per API key)
	sslRateLimiter := rate.NewLimiter(rate.Every(time.Second), 1)

	return &SSLCollector{
		client:         client,
		sslRateLimiter: sslRateLimiter,

		domains: domains,

		Bundle: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "ssl_bundle"),
			"A metric with a constant value of 1 if bundle exists",
			[]string{
				"domain",
			},
			nil,
		),
	}
}

// Collect implements Prometheus' Collector interface and is used to collect metrics
func (c *SSLCollector) Collect(ch chan<- prometheus.Metric) {
	method := "SSLCollector"

	ctx := context.Background()
	var wg sync.WaitGroup
	for _, domain := range c.domains {
		wg.Add(1)
		go func(domain string) {
			defer wg.Done()

			log.Printf("[%s:go] Domain: %s", method, domain)

			// Porkbun API /ssl endpoint has a 1 qps rate limit
			// Before making requests on /ssl endpoint, wait on the limiter
			if err := c.sslRateLimiter.Wait(ctx); err != nil {
				msg := "Porkbun API rate limit exceeded"
				log.Printf("[%s:go] %s for domain (%s)", method, msg, domain)
				return
			}

			_, err := c.client.RetrieveSSLBundle(ctx, domain)
			if err != nil {
				msg := "unable to retrieve SSL bundle for domain"
				log.Printf("[%s:go] %s (%s)\n%s", method, msg, domain, err)
				return
			}

			log.Printf("[%s:go] Domain (%s) contains SSL record", method, domain)

			ch <- prometheus.MustNewConstMetric(
				c.Bundle,
				prometheus.GaugeValue,
				1.0,
				[]string{
					domain,
				}...,
			)
		}(domain)
	}
	wg.Wait()
}

// Describe implements Prometheus' Collector interface and is used to describe metrics
func (c *SSLCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.Bundle
}
