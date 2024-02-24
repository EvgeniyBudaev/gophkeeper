.PHONY: all
all: ;

.PHONY: build-client
build-client:
	go build -o ./cmd/client/bin/gclient ./cmd/client

.PHONY: build
build:
	go build -o ./cmd/gophkeeper/gophkeeper ./cmd/gophkeeper

.PHONY: run
run: build
	./cmd/gophkeeper/gophkeeper -r :8080