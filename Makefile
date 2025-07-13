.PHONY: all build install

all: build

build:
	go build -o gofetch main.go

install: build
	sudo mv gofetch /usr/bin/
