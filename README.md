# A Prometheus Exporter for [Porkbun](https://porkbun.com)

[![build](https://github.com/DazWilkin/porkbun-exporter/actions/workflows/build.yml/badge.svg)](https://github.com/DazWilkin/porkbun-exporter/actions/workflows/build.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/DazWilkin/porkbun-exporter.svg)](https://pkg.go.dev/github.com/DazWilkin/porkbun-exporter)
[![Go Report Card](https://goreportcard.com/badge/github.com/DazWilkin/porkbun-exporter)](https://goreportcard.com/report/github.com/DazWilkin/porkbun-exporter)

## Run

The exporter requires Porkbun API key and secret values provided by the environment. The exports accepts flags for `domains` to be queried, `endpoint` of the exporter, and metrics `path`:

```bash
export APIKEY="..."
export SECRET="..."

HOST_PORT="8080"
CONT_PORT="8080"

# Replace with your list of Porbun API-enabled domains
DOMAINS="example.com,example,org"

podman run \
--interactive --tty --rm \
--name=porkbun-exporter \
--env=APIKEY=${APIKEY} \
--env=SECRET=${SECRET} \
--publish=${HOST_PORT}:${CONT_PORT}/tcp \
ghcr.io/dazwilkin/porkbun-exporter:c43902431f69124f68883c07c3f9977813026566 \
--domains=${DOMAINS} \
--endpoint=:${CONT_PORT} \
--path=/metrics
```

## Build

## Prometheus

```bash
VERS="v2.46.0"

# Binds to host network to scrape Porkbun Exporter
podman run \
--interactive --tty --rm \
--net=host \
--volume=${PWD}/prometheus.yml:/etc/prometheus/prometheus.yml \
--volume=${PWD}/rules.yml:/etc/alertmanager/rules.yml \
quay.io/prometheus/prometheus:${VERS} \
  --config.file=/etc/prometheus/prometheus.yml \
  --web.enable-lifecycle
```

## Metrics

All metrics are prefixed `porkbun_exporter_`

|Name|Type|Description|
|----|----|-----------|
|`porkbun_exporter_build_info`|Counter|A metric with a constant '1' value labeled by OS version, Go version, and the Git commit of the exporter|
|`porkbun_exporter_ssl_bundle`|Gauge|A metric with a constant value of 1 if bundle exists|
|`porkbun_exporter_dns_type`|Gauge|A metric that totals a domain's DNS records by type|
|`porkbun_exporter_start_time`|Gauge|Exporter start time in Unix epoch seconds|

## Alertmanager

```YAML
groups:
  - name: porkbun-exporter
    rules:
    - alert: porkbun_ssl_certs
      expr: sum without(domain) (porkbun_exporter_ssl_bundle{}) != X
      for: 1h
      labels:
        severity: warning
      annotations:
        summary: "Porkbun SSL certificates {{ $value }} (expect X)"
    - alert: porkbun_dns_records
      expr: |
            sum without(name,type) (porkbun_exporter_dns_type{domain="D1"}) != X or
            sum without(name,type) (porkbun_exporter_dns_type{domain="D2"}) != Y
      for: 1h
      labels:
        severity: page
      annotations:
        summary: "Porkbun DNS records changed for {{ $labels.domain }} now {{ $value }} records"
```

## [Sigstore](https://www.sigstore.dev/)

`porkbun-exporter` container images are being signed by [Sigstore](https://www.sigstore.dev/) and may be verified:

```bash
cosign verify \
--key=./cosign.pub \
ghcr.io/dazwilkin/porkbun-exporter:c43902431f69124f68883c07c3f9977813026566
```

> **NOTE** `cosign.pub` may be downloaded [here](https://github.com/DazWilkin/porkbun-exporter/blob/master/cosign.pub)

To install `cosign`, e.g.:
```bash
go install github.com/sigstore/cosign/cmd/cosign@latest
```

## Similar Exporters

+ [Prometheus Exporter for Azure](https://github.com/DazWilkin/azure-exporter)
+ [Prometheus Exporter for crt.sh](https://github.com/DazWilkin/crtsh-exporter)
+ [Prometheus Exporter for Fly.io](https://github.com/DazWilkin/fly-exporter)
+ [Prometheus Exporter for GoatCounter](https://github.com/DazWilkin/goatcounter-exporter)
+ [Prometheus Exporter for Google Cloud](https://github.com/DazWilkin/gcp-exporter)
+ [Prometheus Exporter for Koyeb](https://github.com/DazWilkin/koyeb-exporter)
+ [Prometheus Exporter for Linode](https://github.com/DazWilkin/linode-exporter)
+ [Prometheus Exporter for PorkBun](https://github.com/DazWilkin/porkbun-exporter)
+ [Prometheus Exporter for updown.io](https://github.com/DazWilkin/updown-exporter)
+ [Prometheus Exporter for Vultr](https://github.com/DazWilkin/vultr-exporter)

<hr/>
<br/>
<a href="https://www.buymeacoffee.com/dazwilkin" target="_blank"><img src="https://cdn.buymeacoffee.com/buttons/default-orange.png" alt="Buy Me A Coffee" height="41" width="174"></a>
