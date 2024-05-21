all: run

build:
	go build main.go

run:
	go run main.go

deps:
	go mod tidy

install: build
	cp main /usr/local/bin/MCSCS
	rm main

uninstall:
	rm /usr/local/bin/MCSCS