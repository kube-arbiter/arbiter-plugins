REPOPROJECT := kubearbiter
EXEC := observer-metric-server
RELEASETAG := v0.1.0

clean:
	rm -rf $(EXEC)


build: clean
	CGO_ENABLED=0 GOOS=linux go build -o observer-metric-server main.go

build-image:
	docker build -t $(REPOPROJECT)/$(EXEC):$(RELEASETAG)  .

.DEFAULT_GOAL := build
