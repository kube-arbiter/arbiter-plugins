REPOPROJECT := github.com/kube-arbiter/arbiter-plugins
EXEC := arbiter-prometheus-server
RELEASETAG := v0.0.1

clean:
	rm -rf $(EXEC) server main


build: clean main.go
	CGO_ENABLED=0 GOOS=linux go build -o arbiter-prometheus-server main.go

build-image: build
	docker build -t $(REPOPROJECT)/$(EXEC):$(RELEASETAG)  .

.DEFAULT_GOAL := build