BINARY_NAME=MCST
BUILD_PATH=cmd/$(BINARY_NAME)/main.go

all: run

build:
	go build -o $(BINARY_NAME) $(BUILD_PATH)

run:
	go run $(BUILD_PATH)

test:
	go test -v -parallel 2 ./pkg/...