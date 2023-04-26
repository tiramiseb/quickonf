package conf

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/tiramiseb/quickonf/instructions"
)

func ReadFile(path string) (*instructions.Groups, error) {
	r, err := os.Open(path)
	if err != nil {
		fmt.Println("Could not open configuration file", path)
		fmt.Println(err)
		os.Exit(1)
	}
	defer r.Close()
	absPath, err := filepath.Abs(path)
	if err != nil {
		fmt.Println("Could not get configuration file absolute path")
		fmt.Println(err)
		os.Exit(1)
	}
	instructions.NewGlobalVar("confdir", filepath.Dir(absPath))
	return Read(r)
}

func Read(r io.Reader) (*instructions.Groups, error) {
	lxr := newLexer(r)
	tokens, err := lxr.scan()
	if err != nil {
		return nil, err
	}
	p := newParser(tokens)
	groups := p.parse()
	return instructions.NewGroups(groups), nil
}

func Check(r io.Reader) error {
	lxr := newLexer(r)
	tokens, err := lxr.scan()
	if err != nil {
		return err
	}
	p := newParser(tokens)
	result := p.check()
	json, err := result.ToJSON()
	if err != nil {
		return err
	}
	fmt.Println(json)
	return nil
}
