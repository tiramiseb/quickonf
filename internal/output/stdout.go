package output

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gookit/color"
)

// Stdout is an output writer which writes on stdout
type Stdout struct {
	report *Report

	loaderPosition int
	stopLoading    chan bool
}

// NewStdout returns a new output writer which writes on stdout
func NewStdout() Output {
	return &Stdout{
		report:      &Report{},
		stopLoading: make(chan bool),
	}
}

var loaderImages = []string{
	"\r       [==        ]",
	"\r       [===       ]",
	"\r       [ ===      ]",
	"\r       [  ===     ]",
	"\r       [   ===    ]",
	"\r       [    ===   ]",
	"\r       [     ===  ]",
	"\r       [      === ]",
	"\r       [       ===]",
	"\r       [        ==]",
	"\r       [       ===]",
	"\r       [      === ]",
	"\r       [     ===  ]",
	"\r       [    ===   ]",
	"\r       [   ===    ]",
	"\r       [  ===     ]",
	"\r       [ ===      ]",
	"\r       [===       ]",
}

const loaderClear = "\r                   \r"

var (
	stdoutStepTitle        = color.New(color.FgYellow, color.Bold)
	stdoutInstructionTitle = color.New(color.FgCyan, color.Bold)
	stdoutReportTitle      = color.New(color.FgWhite, color.BgBlue, color.Bold)
	stdoutReportAlertTitle = color.New(color.FgWhite, color.BgYellow, color.Bold)
	stdoutReportErrorTitle = color.New(color.FgWhite, color.BgRed, color.Bold)
)

// StepTitle writes a step title
func (s *Stdout) StepTitle(str string) {
	stdoutStepTitle.Printf("\n" + str + "\n" + strings.Repeat("=", len(str)) + "\n")
	s.report.step(str)
}

// InstructionTitle writes an instruction title
func (s *Stdout) InstructionTitle(str string) {
	stdoutInstructionTitle.Println("  " + str)
	s.report.instruction(str)
}

// Info writes an informational message
func (s *Stdout) Info(str string) {
	color.Info.Println("    ⬩  " + str)
	s.report.info(str)
}

// Infof writes an informational message with format
func (s *Stdout) Infof(format string, a ...interface{}) {
	s.Info(fmt.Sprintf(format, a...))
}

// Success writes a successful message
func (s *Stdout) Success(str string) {
	color.Info.Println("    ✓  " + str)
	s.report.success(str)
}

// Successf writes an informational message with format
func (s *Stdout) Successf(format string, a ...interface{}) {
	s.Success(fmt.Sprintf(format, a...))
}

// Alert writes an alert message
func (s *Stdout) Alert(str string) {
	color.Danger.Println("    ⊖  " + str)
	s.report.alert(str)
}

// Error writes an error message
func (s *Stdout) Error(err error) {
	fmt.Print("     ")
	color.Error.Println(" " + err.Error() + " ")
	s.report.error(err.Error())
}

// ShowLoader displays the loader
func (s *Stdout) ShowLoader() {
	tick := time.Tick(100 * time.Millisecond)
	go func() {
		for {
			select {
			case <-tick:
				s.loaderPosition++
				if s.loaderPosition == 18 {
					s.loaderPosition = 0
				}
				color.Danger.Print(loaderImages[s.loaderPosition])
			case <-s.stopLoading:
				return
			}
		}
	}()
}

// HideLoader hides the loader
func (s *Stdout) HideLoader() {
	s.stopLoading <- true
	fmt.Print(loaderClear)
}

// ShowPercentage displays the loader with percentage
func (s *Stdout) ShowPercentage(p int) {
	switch {
	case p < 5:
		color.Danger.Printf("\r       [    %d%%    ]", p)
	case p < 10:
		color.Danger.Printf("\r       [=   %d%%    ]", p)
	case p < 15:
		color.Danger.Printf("\r       [=  %d%%    ]", p)
	case p < 25:
		color.Danger.Printf("\r       [== %d%%    ]", p)
	case p < 65:
		color.Danger.Printf("\r       [===%d%%    ]", p)
	case p < 75:
		color.Danger.Printf("\r       [===%d%%=   ]", p)
	case p < 85:
		color.Danger.Printf("\r       [===%d%%==  ]", p)
	case p < 95:
		color.Danger.Printf("\r       [===%d%%=== ]", p)
	case p < 100:
		color.Danger.Printf("\r       [===%d%%====]", p)
	case p == 100:
		color.Danger.Print("\r       [===100%===]")
	default:
		color.Danger.Print("\r       [    ??    ]")
	}
}

// HidePercentage hides the loader with percentage
func (s *Stdout) HidePercentage() {
	fmt.Print(loaderClear)
}

// ShowXonY displays the loader with "X/Y" information
func (s *Stdout) ShowXonY(x, y int) {
	xS := strconv.Itoa(x)
	yS := strconv.Itoa(y)
	if len(xS) > 4 {
		xS = xS[len(xS)-5 : len(xS)-1]
	} else {
		xS = fmt.Sprintf("%4s", xS)
	}
	if len(yS) > 5 {
		yS = yS[0:5]
	} else {
		yS = fmt.Sprintf("%-5s", yS)
	}
	color.Danger.Print("\r       [" + xS + "/" + yS + "]")
}

// HideXonY hides the loader with "X/Y" information
func (s *Stdout) HideXonY() {
	fmt.Print(loaderClear)
}

// Report writes the summary
func (s *Stdout) Report() {
	fmt.Println("")
	stdoutReportTitle.Println(" Report ")
	stdoutReportTitle.Println(" ====== ")
	alerts, errors := s.report.alertsAndErrors()

	if len(alerts) > 0 {
		fmt.Print("  ")
		stdoutReportAlertTitle.Println(" Alerts ")
		fmt.Print("  ")
		stdoutReportAlertTitle.Println(" ------ ")
		for _, step := range alerts {
			stdoutStepTitle.Println("    " + step.title)
			if len(step.alerts) > 0 {
				for _, alert := range step.alerts {
					color.Danger.Println("      " + alert)
				}
			}
			for _, module := range step.instructions {
				stdoutInstructionTitle.Println("      " + module.name)
				for _, alert := range module.alerts {
					color.Danger.Println("        " + alert)
				}
			}
		}
	}
	if len(errors) > 0 {
		fmt.Print("  ")
		stdoutReportErrorTitle.Println(" Errors ")
		fmt.Print("  ")
		stdoutReportErrorTitle.Println(" ------ ")
		for _, step := range errors {
			stdoutStepTitle.Println("    " + step.title)
			if len(step.errors) > 0 {
				for _, err := range step.errors {
					color.Danger.Println("      " + err)
				}
			}
			for _, module := range step.instructions {
				stdoutInstructionTitle.Println("      " + module.name)
				for _, err := range module.errors {
					color.Danger.Println("        " + err)
				}
			}
		}
	}
}
