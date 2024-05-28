all: run

build:
	go build main.go

run:
	go run main.go

test:
	go test -v ./...

deps:
	go mod tidy

install: build
	cp main /usr/local/bin/MCSCS
	rm main

uninstall:
	rm /usr/local/bin/MCSCS