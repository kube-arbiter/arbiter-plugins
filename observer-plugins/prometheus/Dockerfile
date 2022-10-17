
FROM golang:1.18 as builder

ARG ARCH=amd64
ARG CGO=0
ARG GOOS=linux


WORKDIR /workspace
# Copy the Go Modules manifests
COPY . .

## Copy the go source and vendor
## Build binary
RUN ./build-prometheus.sh $CGO $ARCH $GOOS
#RUN GOPROXY=https://goproxy.io CGO_ENABLED=${CGO} GOOS=${GOOS} GOARCH=${ARCH} go build -o arbiter-prometheus-server .

FROM alpine:3.15

MAINTAINER 0xff-dev

COPY --from=builder /workspace/observer-prometheus-server /usr/local/bin


CMD ["observer-prometheus-server"]
