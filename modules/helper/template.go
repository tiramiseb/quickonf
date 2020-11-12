package helper

import (
	"bytes"
	"os/exec"
	"strings"
	"text/template"
)

type templateValues struct {
	DistributionCodename string
}

var tmplValues *templateValues

func init() {
	cmdout, err := exec.Command("lsb_release", "--codename", "--short").Output()
	if err != nil {
		panic(err)
	}
	tmplValues = &templateValues{
		DistributionCodename: strings.TrimSpace(string(cmdout)),
	}
}

func Template(src string) ([]byte, error) {
	tmpl, err := template.New("quickonf").Parse(src)
	if err != nil {
		return nil, err
	}
	var result bytes.Buffer
	if err := tmpl.Execute(&result, tmplValues); err != nil {
		return nil, err
	}
	return result.Bytes(), nil
}

