.PHONY: docs

default: intructions/embedded-cookbook.go
	CGO_ENABLED=0 go build -ldflags "-s -w" -o quickonf cmd/*.go
	upx -9 quickonf

intructions/embedded-cookbook.go:
	go generate ./embeddedcookbook/embeder

quick: intructions/embedded-cookbook.go
	CGO_ENABLED=0 go build -o quickonf cmd/*.go

all: x86 x64

x64: intructions/embedded-cookbook.go
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags "-s -w" -o quickonf cmd/*.go
	upx -9 quickonf

x86: intructions/embedded-cookbook.go
	GOOS=linux GOARCH=386 CGO_ENABLED=0 go build -ldflags "-s -w" -o quickonf-32 cmd/*.go
	upx -9 quickonf-32

docs:
	go generate ./docs

svglogo:
	openscad --export-format svg -o - logo.scad | sed 's/stroke-width="0.5"/stroke-width="4"/;s/fill="lightgray"/fill="#800080"/' > docs/assets/logo.svg

pnglogo:
	openscad --export-format svg -o - logo.scad | sed 's/stroke-width="0.5"/stroke-width="4"/;s/fill="lightgray"/fill="#800080"/' | convert -background none - vscode-extension/icon.png
