default:
	CGO_ENABLED=0 go build -ldflags "-s -w" -o quickonf cmd/*.go
	upx -9 quickonf

all: x86 x64

x64:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags "-s -w" -o quickonf cmd/*.go
	upx -9 quickonf

x86:
	GOOS=linux GOARCH=386 CGO_ENABLED=0 go build -ldflags "-s -w" -o quickonf-32 cmd/*.go
	upx -9 quickonf-32
