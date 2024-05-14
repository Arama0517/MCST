VERSION := $(shell go run main.go --version)
NAME := MCSCS
BUILDS_DIR := builds/$(VERSION)

run:
	@go run main.go

test:
	@go test ./... -v

build: 
	go build -o $(NAME) main.go
	mkdir -p $(BUILDS_DIR)
	mv $(NAME) $(BUILDS_DIR)/
	@echo "构建完成, 程序版本: $(VERSION), 输出目录: $(BUILDS_DIR)"

package:
	mkdir -p $(BUILDS_DIR)
	tar -cjf $(BUILDS_DIR)/linux-$(VERSION).tar.bz2 --exclude=.git --exclude=builds .
	tar -cJf $(BUILDS_DIR)/linux-$(VERSION).tar.xz --exclude=.git --exclude=builds .
	tar -czf $(BUILDS_DIR)/linux-$(VERSION).tar.gz --exclude=.git --exclude=builds .
	GOOS=windows
	GOARCH=amd64 go build -o $(NAME).exe main.go
	zip -r $(BUILDS_DIR)/windows-$(VERSION)-32bits.zip MCSCS.exe
	GOARCH=386 go build -o $(NAME).exe main.go 
	zip -r $(BUILDS_DIR)/windows-$(VERSION)-64bits.zip MCSCS.exe 
	rm $(NAME).exe

install: build
	sudo cp $(BUILDS_DIR)/$(NAME) /usr/local/bin/
	@echo "安装完成, 程序版本: $(VERSION)"

uninstall:
	sudo rm /usr/local/bin/$(NAME)
	@echo "卸载完成"

clean:
	rm -rf builds/*

deps:
	go mod tidy