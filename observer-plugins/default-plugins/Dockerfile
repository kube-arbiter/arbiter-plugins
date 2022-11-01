FROM golang:1.18 as builder

ARG CGO=0
ARG ARCH=amd64
ARG GOOS=linux


WORKDIR /workspace

# Copy the Go Modules manifests
COPY . .

## Copy the go source and vendor
## Build binary
RUN ./build.sh ${CGO} ${ARCH} ${GOOS}

FROM alpine:3.15

MAINTAINER 0xff-dev

COPY --from=builder /workspace/default-plugins /usr/local/bin

CMD ["default-plugins"]
