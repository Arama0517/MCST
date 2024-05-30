.PHONY: all build run test deps install uninstall

BINARY_NAME=MCST
BUILD_PATH=cmd/$(BINARY_NAME)/main.go
INSTALL_PATH=/usr/local/bin

all: run

build:
	go build -o $(BINARY_NAME) $(BUILD_PATH)

run:
	go run $(BUILD_PATH)

test:
	go test -v ./pkg/...

deps:
	go mod tidy

install: build
	@if [ -f "$(INSTALL_PATH)/$(BINARY_NAME)" ]; then \
		echo "$(BINARY_NAME) 已安装"; \
	else \
		cp $(BINARY_NAME) $(INSTALL_PATH)/$(BINARY_NAME); \
		rm $(BINARY_NAME) -f; \
	fi

uninstall:
	@if [ -f "$(INSTALL_PATH)/$(BINARY_NAME)" ]; then \
		read -p "你确认要卸载 $(BINARY_NAME) 吗? [y/N] " confirm && [ $$confirm == y ] || [ $$confirm == Y ] && rm $(INSTALL_PATH)/$(BINARY_NAME); \
	else \
		echo "$(BINARY_NAME) 还没有安装"; \
	fi