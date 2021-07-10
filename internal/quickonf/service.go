package quickonf

import (
	"path"
	"strings"

	"github.com/tiramiseb/quickonf/internal/helper"
	"github.com/tiramiseb/quickonf/internal/modules"
	"github.com/tiramiseb/quickonf/internal/output"
)

// Service is a quickonf service
type Service struct {
	steps  []Step
	output output.Output
}

// New creates a new quickonf service containing the given steps
func New(steps []Step, outputName string) (*Service, error) {
	out, err := output.New(outputName)
	if err != nil {
		return nil, err
	}
	return &Service{
		steps:  steps,
		output: out,
	}, nil
}

// Run runs only the selected steps, and the steps marked as "always"
func (s *Service) Run(filter []string, dryrun bool) {
	modules.Dryrun = dryrun
	helper.Dryrun = dryrun
	if len(filter) > 0 {
		s.output.StepTitle("Running steps matching:")
		for _, step := range filter {
			s.output.Info(step)
		}
	}
	for _, step := range s.steps {
		if step.Always() {
			step.run(s.output)
			continue
		}
		if len(filter) > 0 {
			for _, wanted := range filter {
				if ok, _ := path.Match(wanted, strings.ReplaceAll(strings.ToLower(step.Name()), "/", " ")); ok {
					step.run(s.output)
					break
				}
			}
		} else {
			step.run(s.output)
		}
	}
	s.output.Report()
}

// List lists the steps contained in the quickonf service
func (s *Service) List() {
	s.output.StepTitle("List of steps")
	for _, step := range s.steps {
		step.list(s.output)
	}
}
