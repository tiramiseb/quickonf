all: build

build:
	CGO_ENABLED=0 go build -o quickonf cmd/main.go
