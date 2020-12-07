package quickonf

import (
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

// List lists the steps contained in the quickonf service
func (s *Service) List() {
	s.output.StepTitle("List of steps")
	for _, step := range s.steps {
		step.list(s.output)
	}
}
