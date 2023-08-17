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

podman run \
--interactive --tty --rm \
--name=porkbun-exporter \
--env=APIKEY=${APIKEY} \
--env=SECRET=${SECRET} \
--publish=${HOST_PORT}:${CONT_PORT}/tcp \
ghcr.io/dazwilkin/porkbun-exporter:3157c4ced91f8ed6b9cae2c71c96121fad25e3d5 \
--domains=example.com,example,org \
--endpoint=:${CONT_PORT} \
--path=/metrics
```

## Build



## Metrics

All metrics are prefixed `porkbun_exporter_`

|Name|Type|Description|
|----|----|-----------|
|`porkbun_exporter_build_info`|Counter||
|`porkbun_exporter_dns_type`|Gauge||
|`porkbun_exporter_start_time`|Gauge||

## [Sigstore](https://www.sigstore.dev/)

`porkbun-exporter` container images are being signed by [Sigstore](https://www.sigstore.dev/) and may be verified:

```bash
cosign verify \
--key=./cosign.pub \
ghcr.io/dazwilkin/porkbun-exporter:3157c4ced91f8ed6b9cae2c71c96121fad25e3d5
```

> **NOTE** `cosign.pub` may be downloaded [here](https://github.com/DazWilkin/porkbun-exporter/blob/master/cosign.pub)

To install `cosign`, e.g.:
```bash
go install github.com/sigstore/cosign/cmd/cosign@latest
```