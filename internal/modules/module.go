package modules

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/tiramiseb/quickonf/internal/helper"
	"github.com/tiramiseb/quickonf/internal/output"
	"gopkg.in/yaml.v3"
)

// Instruction is an instruction. Returns true if succeeds, false if there has been an error
type Instruction func(in interface{}, out output.Output) error

var registry = map[string]Instruction{}

// Dryrun allows running instructions without system modification
var Dryrun = false

// Register adds an instruction
func Register(name string, instruction Instruction) {
	registry[name] = instruction
}

// Get gets an instruction
func Get(name string) Instruction {
	return registry[name]
}

// RunAs runs the given instruction with the given data as the given user.
// It is not a regular module, but executed explicitly in the step.
func RunAs(user, instruction string, in interface{}, out output.Output) error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	fullPath := filepath.Join(cwd, os.Args[0])
	// Remove automatic store values from the forwarded values
	allStore := helper.AllStore()
	delete(allStore, "home")
	delete(allStore, "oscodename")
	yamlSource := []map[string]interface{}{}
	if len(allStore) > 0 {
		yamlSource = append(yamlSource, map[string]interface{}{
			"STORE": []map[string]interface{}{
				{
					"store": allStore,
				},
			},
		})
	}
	yamlSource = append(yamlSource, map[string]interface{}{
		"EXECUTE": []map[string]interface{}{
			{
				instruction: in,
			},
		},
	})
	f, err := ioutil.TempFile("", fmt.Sprintf("quickonf-tempconf-%s-%s-", instruction, user))
	if err != nil {
		return err
	}
	e := yaml.NewEncoder(f)
	if err := e.Encode(yamlSource); err != nil {
		return err
	}
	e.Close()
	f.Close()
	defer os.Remove(f.Name())
	if err := os.Chmod(f.Name(), 0644); err != nil {
		return err
	}

	command := exec.Command("sudo", "--reset-timestamp", "--prompt=", "--stdin", "-u", user, fullPath, "--config", f.Name(), "--output", "prog")
	command.Env = append(os.Environ(), "LANG=C")
	command.Stdin = strings.NewReader(helper.SudoPassword)
	outPipe, err := command.StdoutPipe()
	if err != nil {
		return err
	}
	var stderr bytes.Buffer
	command.Stderr = &stderr
	outputScanner := bufio.NewScanner(outPipe)
	done := make(chan bool)

	go func() {
		var inExecute bool
		for outputScanner.Scan() {
			txt := strings.SplitN(outputScanner.Text(), ":", 2)
			prefix := txt[0]
			var content string
			if len(txt) == 2 {
				content = txt[1]
			}
			switch prefix {
			case "STEP":
				if content == "EXECUTE" {
					inExecute = true
				}
			case "INSTRUCTION":
				if !inExecute {
					break
				}
				out.InstructionTitle("[as " + user + "] " + content)
			case "INFO":
				if !inExecute {
					break
				}
				out.Info(content)
			case "SUCCESS":
				if !inExecute {
					break
				}
				out.Success(content)
			case "ALERT":
				if !inExecute {
					break
				}
				out.Alert(content)
			case "ERROR":
				if !inExecute {
					break
				}
				out.Error(errors.New(content))
			case "LOADER":
				if !inExecute {
					break
				}
				out.ShowLoader()
			case "HIDE LOADER":
				if !inExecute {
					break
				}
				out.HideLoader()
			case "PERCENTAGE":
				if !inExecute {
					break
				}
				pc, err := strconv.Atoi(content)
				if err != nil {
					out.Error(err)
					break
				}
				out.ShowPercentage(pc)
			case "HIDE PERCENTAGE":
				if !inExecute {
					break
				}
				out.HidePercentage()
			case "XY":
				if !inExecute {
					break
				}
				xy := strings.SplitN(content, ":", 2)
				if len(xy) != 2 {
					out.Error(errors.New("Does not contain 2 values: " + content))
				}
				x, err := strconv.Atoi(xy[0])
				if err != nil {
					out.Error(err)
					break
				}
				y, err := strconv.Atoi(xy[1])
				if err != nil {
					out.Error(err)
					break
				}
				out.ShowXonY(x, y)
			case "HIDE XY":
				if !inExecute {
					break
				}
				out.HideXonY()
			}
		}
		err = outputScanner.Err()
		done <- true
	}()
	if err := command.Start(); err != nil {
		return err
	}
	<-done
	if err != nil {
		return err
	}
	if err := command.Wait(); err != nil {
		if stderr.Len() > 0 {
			return errors.New(stderr.String())
		}
		return err
	}
	return nil
}
