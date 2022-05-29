
run:
	go run cmd/server.go
.PHONY: run

build-docker:
	docker build -t mud:latest .
.PHONY: build-docker

run-docker:
	docker run -it mud:latest
.PHONY: run-docker

test:
	go test ./...
.PHONY: test