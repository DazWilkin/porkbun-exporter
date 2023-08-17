package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/DazWilkin/porkbun-exporter/collector"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	APIKey string = "APIKEY"
	Secret string = "SECRET"
)
const (
	sRoot string = `
<h2>A Prometheus Exporter for <a href="https://porkbun.com">Porkbun</a></h2>
<ul>
	<li><a href="{{ .Metrics }}">metrics</a></li>
	<li><a href="/healthz">healthz</a></li>
</ul>`
)

var (
	// GitCommit is the git commit value and is expected to be set during build
	GitCommit string
	// GoVersion is the Golang runtime version
	GoVersion = runtime.Version()
	// OSVersion is the OS version (uname --kernel-release) and is expected to be set during build
	OSVersion string
	// StartTime is the start time of the exporter represented as a UNIX epoch
	StartTime = time.Now().Unix()
)
var (
	domainList  = flag.String("domains", "", "Comma-separated list of domains")
	endpoint    = flag.String("endpoint", ":8080", "The endpoint of the HTTP server")
	metricsPath = flag.String("path", "/metrics", "The path on which Prometheus metrics will be served")
)
var (
	tRoot = template.Must(template.New("root").Parse(sRoot))
)

func handleHealthz(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}
func handleRoot(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=UTF-8")
	fmt.Fprint(w)
	tRoot.Execute(w, struct {
		Metrics string
	}{
		Metrics: *metricsPath,
	})
}
func main() {
	apikey := os.Getenv(APIKey)
	if apikey == "" {
		log.Fatalf("Expected environment to contain  '%s' variable", APIKey)
	}
	secret := os.Getenv(Secret)
	if secret == "" {
		log.Fatalf("Expected environment to contain '%s' variable", secret)
	}

	// For domains,endpoint and metricsPath
	flag.Parse()
	domains := strings.Split(*domainList, ",")
	if len(domains) == 0 {
		log.Fatal("Need at least one domain")
	}

	if GitCommit == "" {
		log.Println("[main] GitCommit value unchanged: expected to be set during build")
	}
	if OSVersion == "" {
		log.Println("[main] OSVersion value unchanged: expected to be set during build")
	}

	registry := prometheus.NewRegistry()
	registry.MustRegister(collector.NewExporterCollector(OSVersion, GoVersion, GitCommit, StartTime))
	registry.MustRegister(collector.NewDomainCollector(apikey, secret, domains))

	mux := http.NewServeMux()
	mux.Handle("/", http.HandlerFunc(handleRoot))
	mux.Handle("/healthz", http.HandlerFunc(handleHealthz))
	mux.Handle(*metricsPath, promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))

	log.Printf("[main] Server starting (%s)", *endpoint)
	log.Printf("[main] metrics served on: %s", *metricsPath)
	log.Fatal(http.ListenAndServe(*endpoint, mux))
}
