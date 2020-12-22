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
func New(steps []Step) (*Service, error) {
	return &Service{
		steps:  steps,
		output: output.NewStdout(),
	}, nil
}

// Run runs the steps contained in the quickonf service
func (s *Service) Run(dryrun bool) {
	modules.Dryrun = dryrun
	helper.Dryrun = dryrun
	for _, step := range s.steps {
		step.run(s.output)
	}
	s.output.Report()
}

// Steps runs the selected steps
func (s *Service) Steps(steps []string, dryrun bool) {
	modules.Dryrun = dryrun
	helper.Dryrun = dryrun
	s.output.StepTitle("Running steps matching:")
	for _, step := range steps {
		s.output.Info(step)
	}
	for _, step := range s.steps {
		for _, wanted := range steps {
			if ok, _ := path.Match(wanted, strings.ReplaceAll(strings.ToLower(step.Name()), "/", " ")); ok {
				step.run(s.output)
			}
		}
	}
}

// List lists the steps contained in the quickonf service
func (s *Service) List() {
	s.output.StepTitle("List of steps")
	for _, step := range s.steps {
		step.list(s.output)
	}
}
