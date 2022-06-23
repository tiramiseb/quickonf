default:
	CGO_ENABLED=0 go build -ldflags "-s -w --extldflags=-static" -tags osusergo,netgo -o quickonf cmd/*.go

all: x86 x64

x64:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags "-s -w --extldflags=-static" -tags osusergo,netgo -o quickonf cmd/*.go

x86:
	GOOS=linux GOARCH=386 CGO_ENABLED=0 go build -ldflags "-s -w --extldflags=-static" -tags osusergo,netgo -o quickonf-32 cmd/*.go
