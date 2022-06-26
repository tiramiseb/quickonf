package main

//go:generate go run .

import (
	"bytes"
	_ "embed"
	"io"
	"os"
)

const target = "../../instructions/embedded-cookbook.go"

func main() {
	var buf bytes.Buffer
	e := &embeder{&buf}
	e.make()

	f, err := os.Create(target)
	if err != nil {
		panic(err)
	}

	if _, err := io.Copy(f, &buf); err != nil {
		panic(err)
	}

	if err := f.Close(); err != nil {
		panic(err)
	}
}
