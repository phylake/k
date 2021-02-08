.PHONY: default prebuild install

default: install

prebuild:
	gofmt -s -w .

install: prebuild
	go build -o ~/bin/k main.go
