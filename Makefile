.PHONY: build-all clean

CGO ?=0
GOOS ?= linux

# now only support amd64
ARCH ?= $(shell go env GOARCH)
#GOFLAGS ?=""
RELEASE ?=v0.1.0

echo:
	ROOT=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)/..;
	echo $(ROOT)

clean:
	@./scripts/clean.sh $(RELEASE)

build-all: clean
	./scripts/build.sh $(CGO)  amd64 $(GOOS)


copyright:
	./scripts/verify-copyright.sh

build-image:
	docker build  --build-arg CGO=${CGO} --build-arg ARCH=${ARCH} --build-arg GOOS=${GOOS} -t kubearbiter/observer-metric-server:$(RELEASE) -f build/dockerfile.metric-server .
	docker build  --build-arg CGO=${CGO} --build-arg ARCH=${ARCH} --build-arg GOOS=${GOOS} -t kubearbiter/observer-prometheus-server:$(RELEASE) -f build/dockerfile.prometheus .
	docker build  --build-arg CGO=${CGO} --build-arg ARCH=${ARCH} --build-arg GOOS=${GOOS} -t kubearbiter/executor-resource-tagger:$(RELEASE)  -f build/dockerfile.resource-tagger .
