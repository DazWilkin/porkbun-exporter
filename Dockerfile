ARG GOLANG_VERSION=1.23.0

ARG GOOS=linux
ARG GOARCH=amd64

ARG COMMIT
ARG VERSION

FROM docker.io/golang:${GOLANG_VERSION} as build

WORKDIR /porkbun-exporter

COPY go.* ./
COPY main.go .
COPY collector ./collector

ARG GOOS
ARG GOARCH

ARG VERSION
ARG COMMIT

RUN CGO_ENABLED=0 GOOS=${GOOS} GOARCH=${GOARCH} \
    go build \
    -ldflags "-X main.OSVersion=${VERSION} -X main.GitCommit=${COMMIT}" \
    -a -installsuffix cgo \
    -o /go/bin/porkbun-exporter \
    ./main.go

FROM gcr.io/distroless/static-debian11:latest

LABEL org.opencontainers.image.description "Prometheus Exporter for Porkbun"
LABEL org.opencontainers.image.source https://github.com/DazWilkin/porkbun-exporter

COPY --from=build /go/bin/porkbun-exporter /

EXPOSE 8080

ENTRYPOINT ["/porkbun-exporter"]
