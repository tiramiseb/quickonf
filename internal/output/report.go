package output

type reportPartType int

const (
	reportPartTypeStepTitle reportPartType = iota
	reportPartTypeInstructionTitle
	reportPartTypeMessage
	reportPartTypeWarn
	reportPartTypeAlert
	reportPartTypeError
)

type reportInstruction struct {
	name      string
	infos     []string
	successes []string
	alerts    []string
	errors    []string
}

type reportStep struct {
	title           string
	instructions    []*reportInstruction
	lastInstruction *reportInstruction
	infos           []string
	successes       []string
	alerts          []string
	errors          []string
}

// Report gathers data for reporting
type Report struct {
	steps    []*reportStep
	lastStep *reportStep
}

func (r *Report) step(title string) {
	step := &reportStep{title: title}
	r.steps = append(r.steps, step)
	r.lastStep = step
}

func (r *Report) instruction(name string) {
	instruction := &reportInstruction{name: name}
	if r.lastStep != nil {
		r.lastStep.instructions = append(r.lastStep.instructions, instruction)
		r.lastStep.lastInstruction = instruction
	}
}

func (r *Report) info(str string) {
	if r.lastStep != nil {
		if r.lastStep.lastInstruction != nil {
			r.lastStep.lastInstruction.infos = append(r.lastStep.lastInstruction.infos, str)
		} else {
			r.lastStep.infos = append(r.lastStep.infos, str)
		}
	}
}

func (r *Report) success(str string) {
	if r.lastStep != nil {
		if r.lastStep.lastInstruction != nil {
			r.lastStep.lastInstruction.successes = append(r.lastStep.lastInstruction.successes, str)
		} else {
			r.lastStep.successes = append(r.lastStep.successes, str)
		}
	}
}

func (r *Report) alert(str string) {
	if r.lastStep != nil {
		if r.lastStep.lastInstruction != nil {
			r.lastStep.lastInstruction.alerts = append(r.lastStep.lastInstruction.alerts, str)
		} else {
			r.lastStep.alerts = append(r.lastStep.alerts, str)
		}
	}
}

func (r *Report) error(str string) {
	if r.lastStep != nil {
		if r.lastStep.lastInstruction != nil {
			r.lastStep.lastInstruction.errors = append(r.lastStep.lastInstruction.errors, str)
		} else {
			r.lastStep.errors = append(r.lastStep.errors, str)
		}
	}
}

func (r *Report) alertsAndErrors() ([]*reportStep, []*reportStep) {
	alertSteps := []*reportStep{}
	errorSteps := []*reportStep{}
	for _, step := range r.steps {
		alertStep := &reportStep{title: step.title}
		errorStep := &reportStep{title: step.title}
		if len(step.alerts) > 0 {
			alertStep.alerts = append(alertStep.alerts, step.alerts...)
		}
		if len(step.errors) > 0 {
			errorStep.errors = append(errorStep.errors, step.errors...)
		}
		for _, instruction := range step.instructions {
			if len(instruction.alerts) > 0 {
				alertStep.instructions = append(alertStep.instructions, instruction)
			}
			if len(instruction.errors) > 0 {
				errorStep.instructions = append(errorStep.instructions, instruction)
			}
		}
		if len(alertStep.alerts) > 0 || len(alertStep.instructions) > 0 {
			alertSteps = append(alertSteps, alertStep)
		}
		if len(errorStep.errors) > 0 || len(errorStep.instructions) > 0 {
			errorSteps = append(errorSteps, errorStep)
		}

	}
	return alertSteps, errorSteps
}
