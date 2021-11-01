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

func init() {
	Register("as", As)
}

// As executes instructions as a specific user
func As(in interface{}, out output.Output) error {
	data, err := helper.MapStringInterface(in)
	if err != nil {
		return err
	}

	// Remove automatic store values from the forwarded values, because they will be automatically filled again

	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	fullPath := filepath.Join(cwd, os.Args[0])
	allStore := helper.AllStore()
	delete(allStore, "home")
	delete(allStore, "hostname")
	delete(allStore, "oscodename")
	for user, instructions := range data {
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
			"EXECUTE": instructions,
		})
		f, err := ioutil.TempFile("", fmt.Sprintf("quickonf-tempconf-as-%s-", user))
		if err != nil {
			return err
		}
		e := yaml.NewEncoder(f)
		if err := e.Encode(yamlSource); err != nil {
			f.Close()
			os.Remove(f.Name())
			return err
		}
		e.Close()
		f.Close()
		defer os.Remove(f.Name())
		if err := os.Chmod(f.Name(), 0644); err != nil {
			return err
		}
		args := []string{"--reset-timestamp", "--prompt=", "--stdin", "--user", user, "--", fullPath, "-config", f.Name(), "-output", "prog"}
		if Dryrun {
			args = append(args, "-dry-run")
		}
		command := exec.Command("sudo", args...)
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
						out.Error(fmt.Errorf("does not contain 2 values: %v", content))
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
	}
	return nil
}
