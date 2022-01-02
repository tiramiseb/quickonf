package commands

import "fmt"

var cleaningRegistry = []func() error{}

func registerClean(f func() error) {
	cleaningRegistry = append(cleaningRegistry, f)
}

func Clean() {
	for _, f := range cleaningRegistry {
		if err := f(); err != nil {
			fmt.Println(err)
		}
	}
}
