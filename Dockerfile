ARG GOLANG_VERSION=1.24.3

ARG TARGETOS
ARG TARGETARCH

ARG COMMIT
ARG VERSION

FROM --platform=${TARGETARCH} docker.io/golang:${GOLANG_VERSION} AS build

WORKDIR /porkbun-exporter

COPY go.* ./
COPY main.go .
COPY collector ./collector

ARG TARGETOS
ARG TARGETARCH

ARG VERSION
ARG COMMIT

RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} \
    go build \
    -ldflags "-X main.OSVersion=${VERSION} -X main.GitCommit=${COMMIT}" \
    -a -installsuffix cgo \
    -o /go/bin/porkbun-exporter \
    ./main.go

FROM --platform=${TARGETARCH} gcr.io/distroless/static-debian12:latest

LABEL org.opencontainers.image.description="Prometheus Exporter for Porkbun"
LABEL org.opencontainers.image.source=https://github.com/DazWilkin/porkbun-exporter

COPY --from=build /go/bin/porkbun-exporter /

EXPOSE 8080

ENTRYPOINT ["/porkbun-exporter"]
