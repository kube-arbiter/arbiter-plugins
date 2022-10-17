REPOPROJECT := kubearbiter
EXEC := observer-prometheus-server
RELEASETAG := v0.1.0

clean:
	rm -rf $(EXEC) server main


build: clean main.go
	CGO_ENABLED=0 GOOS=linux go build -o observer-prometheus-server main.go

build-image: build
	docker build -t $(REPOPROJECT)/$(EXEC):$(RELEASETAG)  .

.DEFAULT_GOAL := build
