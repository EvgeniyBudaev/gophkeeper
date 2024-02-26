.PHONY: all
all: ;

.PHONY: build-client
build-client:
	go build -o ./bin/gclient ./cmd/client

.PHONY: build
build:
	go build -o ./bin/gophkeeper ./cmd/gophkeeper

.PHONY: run
run: build
	./bin/gophkeeper -r :8080