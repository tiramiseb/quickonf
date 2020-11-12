package output

import "fmt"

type reportPartType int

const (
	reportPartTypeStepTitle reportPartType = iota
	reportPartTypeModuleName
	reportPartTypeMessage
	reportPartTypeWarn
	reportPartTypeAlert
	reportPartTypeError
)

type reportModule struct {
	name      string
	infos     []string
	successes []string
	alerts    []string
	errors    []string
}

type reportStep struct {
	title      string
	modules    []*reportModule
	lastModule *reportModule
	infos      []string
	successes  []string
	alerts     []string
	errors     []string
}

type Report struct {
	steps    []*reportStep
	lastStep *reportStep
}

func (r *Report) Step(title string) {
	step := &reportStep{title: title}
	r.steps = append(r.steps, step)
	r.lastStep = step
}

func (r *Report) Module(name string) {
	module := &reportModule{name: name}
	r.lastStep.modules = append(r.lastStep.modules, module)
	r.lastStep.lastModule = module
}

func (r *Report) Info(str string) {
	if r.lastStep.lastModule != nil {
		r.lastStep.lastModule.infos = append(r.lastStep.lastModule.infos, str)
	} else {
		r.lastStep.infos = append(r.lastStep.infos, str)
	}
}

func (r *Report) Success(str string) {
	if r.lastStep.lastModule != nil {
		r.lastStep.lastModule.successes = append(r.lastStep.lastModule.successes, str)
	} else {
		r.lastStep.successes = append(r.lastStep.successes, str)
	}
}

func (r *Report) Alert(str string) {
	if r.lastStep.lastModule != nil {
		r.lastStep.lastModule.alerts = append(r.lastStep.lastModule.alerts, str)
	} else {
		r.lastStep.alerts = append(r.lastStep.alerts, str)
	}
}

func (r *Report) Error(str string) {
	fmt.Println("---------")
	fmt.Println(str)
	if r.lastStep.lastModule != nil {
		r.lastStep.lastModule.errors = append(r.lastStep.lastModule.errors, str)
	} else {
		r.lastStep.errors = append(r.lastStep.errors, str)
	}
}

func (r *Report) AlertsAndErrors() ([]*reportStep, []*reportStep) {
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
		for _, module := range step.modules {
			if len(module.alerts) > 0 {
				alertStep.modules = append(alertStep.modules, module)
			}
			if len(module.errors) > 0 {
				errorStep.modules = append(errorStep.modules, module)
			}
		}
		if len(alertStep.alerts) > 0 || len(alertStep.modules) > 0 {
			alertSteps = append(alertSteps, alertStep)
		}
		if len(errorStep.errors) > 0 || len(errorStep.modules) > 0 {
			errorSteps = append(errorSteps, errorStep)
		}

	}
	return alertSteps, errorSteps
}
