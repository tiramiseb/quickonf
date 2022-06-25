package main

//go:generate go run .

import (
	"bytes"
	_ "embed"
	"fmt"
	"io"
	"os"

	"github.com/tiramiseb/quickonf/conf"
)

const target = "../instructions/embedded-cookbook.go"

func main() {
	f, err := os.Open("embedded-cookbook.qconf")
	if err != nil {
		panic(err)
	}

	groups, errs := conf.Read(f)
	if len(errs) > 0 {
		f.Close()
		for _, err := range errs {
			fmt.Println(err)
		}
		os.Exit(1)
	}

	if err := f.Close(); err != nil {
		panic(err)
	}

	var buf bytes.Buffer
	e := &embeder{groups, &buf}

	e.make()

	f, err = os.Create(target)
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
