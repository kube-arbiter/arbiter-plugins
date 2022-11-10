# Global environment variables
IMAGE_REPO_ADDR = docker.io
VERSION = v0.2.0


##@ General

# The help target prints out all targets with their descriptions organized
# beneath their categories. The categories are represented by '##@' and the
# target descriptions by '##'. The awk commands is responsible for reading the
# entire set of makefiles included in this invocation, looking for lines of the
# file as xyz: ## something, and then pretty-format the target and help. Then,
# if there's a line with ##@ something, that gets pretty-printed as a category.
# More info on the usage of ANSI control characters for terminal formatting:
# https://en.wikipedia.org/wiki/ANSI_escape_code#SGR_parameters
# More info on the awk command:
# http://linuxcommand.org/lc3_adv_awk.php

help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Development

install-requirements: _download-protoc _unzip-protoc _locate-protoc _cleanup-tmp _protoc-gen-tools  ## Install protobuf and code generation tools.

_download-protoc:
	mkdir -p tmp && \
	cd tmp && \
	curl -L https://github.com/protocolbuffers/protobuf/releases/download/v21.0-rc2/protoc-21.0-rc-2-osx-universal_binary.zip --output protoc.zip

_unzip-protoc:
	cd tmp && \
	unzip ./protoc.zip -d protoc

_locate-protoc:
	mkdir -p bin && \
	rm -rf bin/protoc && \
	cd tmp && \
	mv -f ./protoc/ ../bin
	chmod u+x bin/protoc/bin/protoc

_cleanup-tmp:
	rm -rf ./tmp

_protoc-gen-tools:
	go mod vendor
	go install github.com/golang/protobuf/protoc-gen-go
	go install github.com/pseudomuto/protoc-gen-doc/cmd/protoc-gen-doc

protoc: _gen-go-out #_gen-proto-doc  ## Generate protobuf go code and documentation.

_gen-go-out:
	mkdir -p internal
	protoc --go-grpc_out=plugins=grpc:./ proto/*.proto

_gen-proto-doc:
	mkdir -p doc/proto
	bin/protoc/bin/protoc --doc_out=html,index.html:./doc/proto proto/*.proto

##@ Local run

run:  ## Run server and client locally.
	go run cmd/server/server.go &
	sleep 3
	true

##@ Build

build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build  -a -o executor-default-plugins cmd/server/server.go

build-image: ## Build executor-default-plugins plugin image.
	docker build -t $(IMAGE_REPO_ADDR)/kubearbiter/executor-default-plugins:$(VERSION) .

updateimg: ## Update executor-default-plugins plugin image.
	kubectl get pod -n arbiter | grep executor | awk '{print $$1}' | xargs kubectl delete pod -n arbiter

##@ Run

kubelog-executor-default-plugins: ## Show executor-default-plugins plugin log in K8S environment.
	kubectl get pod -n arbiter | grep executor | awk '{print $$1}' | xargs kubectl logs -f -n arbiter -c executor-default-plugins

kubeexec-client: ## Test gRPC client in K8S environment.
	kubectl get pod -n arbiter | grep executor | awk '{cmd="kubectl exec -n arbiter "$$1" -c executor-default-plugins -- /client"; system(cmd)}'
