package quickonf

import "github.com/tiramiseb/quickonf/output"

type Service struct {
	steps  []Step
	output output.Output
}

func New(steps []Step) (*Service, error) {
	return &Service{
		steps:  steps,
		output: output.NewStdout(),
	}, nil
}

func (s *Service) Run() {
	for _, step := range s.steps {
		step.Run(s.output)
	}
	s.output.Report()
}
