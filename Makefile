VERSION := @go run main.go
NAME := MCSCS
BUILDS_DIR := builds/$(VERSION)

run:
	@go run main.go

version: 
	@go run main.go -v

build: 
	go build -o $(NAME) main.go
	mv $(NAME) $(BUILDS_DIR)
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

clean:
	rm -rf builds/*

deps:
	go mod tidy