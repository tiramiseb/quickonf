default:
	CGO_ENABLED=0 go build -o quickonf main.go

all: x86 x64

x64:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o quickonf main.go

x86:
	GOOS=linux GOARCH=386 CGO_ENABLED=0 go build -o quickonf-32 main.go
