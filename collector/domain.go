package collector

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/nrdcg/porkbun"

	"github.com/prometheus/client_golang/prometheus"
)

type DomainCollector struct {
	client  *porkbun.Client
	domains []string

	// Metrics
	DNSTypes *prometheus.Desc
}

func NewDomainCollector(apikey, secret string, domains []string) *DomainCollector {
	client := porkbun.New(secret, apikey)

	return &DomainCollector{
		client:  client,
		domains: domains,

		DNSTypes: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "dns_type"),
			"A metric that totals a domain's DNS record types",
			[]string{
				"domain",
				"type",
			},
			nil,
		),
	}
}

// Collect implements Prometheus' Collector interface and is used to collect metrics
func (c *DomainCollector) Collect(ch chan<- prometheus.Metric) {
	ctx := context.Background()
	var wg sync.WaitGroup
	for _, domain := range c.domains {
		wg.Add(1)
		go func(domain string) {
			defer wg.Done()
			records, err := c.client.RetrieveRecords(ctx, domain)
			if err != nil {
				msg := fmt.Sprintf("unable to retrieve records for domain (%s)", domain)
				log.Print(msg)
				return
			}
			dnsTypes := make(map[string]uint16)

			// Enumerate records
			for _, record := range records {
				dnsTypes[record.Type]++
			}

			// Enumerate types
			for dnsType, count := range dnsTypes {
				ch <- prometheus.MustNewConstMetric(
					c.DNSTypes,
					prometheus.GaugeValue,
					float64(count),
					[]string{
						domain,
						dnsType,
					}...,
				)
			}
		}(domain)
	}
	wg.Wait()
}

// Describe implements Prometheus' Collector interface and is used to describe metrics
func (c *DomainCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.DNSTypes
}
