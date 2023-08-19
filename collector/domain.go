package collector

import (
	"context"
	"log"
	"sync"

	"github.com/nrdcg/porkbun"

	"github.com/prometheus/client_golang/prometheus"
)

// DomainCollector is a Promethues collector that represents the Porkbun SSL functionality
type DomainCollector struct {
	client  *porkbun.Client
	domains []string

	// Metrics
	DNSTypes *prometheus.Desc
}

// NewDomainCollector is a function that creates a new DomainCollector
func NewDomainCollector(apikey, secret string, domains []string) *DomainCollector {
	client := porkbun.New(secret, apikey)

	return &DomainCollector{
		client:  client,
		domains: domains,

		DNSTypes: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "dns_type"),
			"A metric that totals a domain's DNS records by type",
			[]string{
				"domain",
				"type",
				"name",
			},
			nil,
		),
	}
}

// Collect implements Prometheus' Collector interface and is used to collect metrics
func (c *DomainCollector) Collect(ch chan<- prometheus.Metric) {
	method := "DomainCollector"

	ctx := context.Background()
	var wg sync.WaitGroup
	for _, domain := range c.domains {
		wg.Add(1)
		go func(domain string) {
			defer wg.Done()

			log.Printf("[%s:go] Domain: %s", method, domain)

			// Porkbun API has a 1 query/second rate limit
			// To effect this apiRateLimiter is shared across collectors
			// Before making Porkbun API requests, wait on the limiter
			if err := apiRateLimiter.Wait(ctx); err != nil {
				msg := "Porkbun API rate limit exceeded"
				log.Printf("[%s:go] %s for domain (%s)", method, msg, domain)
				return
			}
			records, err := c.client.RetrieveRecords(ctx, domain)
			if err != nil {
				msg := "unable to retrieve records for domain"
				log.Printf("[%s:go] %s (%s)", method, msg, domain)
				return
			}

			// DNS records comprise
			// Type (e.g. A)
			// Name (e.g. www.google.com)
			// Data (e.g.192.168.1.1)
			// The Data value may be extensive and is unlikely to be useful as a label
			// The rather awkward map, maps type:name:count(name)
			// For each Type, we count the number of occurrences of each Data value
			// (A,192.168.1.1,1)
			dnsTypesByNameByCount := make(map[string]map[string]uint16)

			// Enumerate records
			// Generating Type*Name*count values
			log.Printf("[%s:go] Domain (%s) contains %d records", method, domain, len(records))
			for _, record := range records {
				if _, ok := dnsTypesByNameByCount[record.Type]; !ok {
					dnsTypesByNameByCount[record.Type] = make(map[string]uint16)
				}
				dnsTypesByNameByCount[record.Type][record.Name]++
			}

			// Enumerate Type*Name label pairs with count as metric value
			for dnsType, NameByCount := range dnsTypesByNameByCount {
				for dnsName, count := range NameByCount {
					ch <- prometheus.MustNewConstMetric(
						c.DNSTypes,
						prometheus.GaugeValue,
						float64(count),
						[]string{
							domain,
							dnsType,
							dnsName,
						}...,
					)
				}
			}
		}(domain)
	}
	wg.Wait()
}

// Describe implements Prometheus' Collector interface and is used to describe metrics
func (c *DomainCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.DNSTypes
}
