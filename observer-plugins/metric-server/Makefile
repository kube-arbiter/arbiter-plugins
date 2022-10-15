REPOPROJECT := github.com/kube-arbiter/arbiter-plugins
EXEC := arbiter-metric-server
RELEASETAG := v0.0.1

clean:
	rm -rf $(EXEC)


build: clean main.go server.go
	CGO_ENABLED=0 GOOS=linux go build -o arbiter-metric-server main.go server.go metric_types.go

build-image: build
	docker build -t $(REPOPROJECT)/$(EXEC):$(RELEASETAG)  .

.DEFAULT_GOAL := build
