# A Prometheus Exporter for [Porkbun](https://porkbun.com)

## Run

The exporter requires Porkbun API key and secret values provided by the environment:

```bash
export APIKEY="..."
export SECRET="..."
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