package output

import (
	"fmt"
	"strings"
	"time"

	"github.com/gookit/color"
)

type Stdout struct {
	report *Report

	loaderPosition int
	stopLoading    chan bool
}

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
	stdoutModuleName       = color.New(color.FgCyan, color.Bold)
	stdoutReportTitle      = color.New(color.FgWhite, color.BgBlue, color.Bold)
	stdoutReportAlertTitle = color.New(color.FgWhite, color.BgYellow, color.Bold)
	stdoutReportErrorTitle = color.New(color.FgWhite, color.BgRed, color.Bold)
)

func (s *Stdout) StepTitle(str string) {
	stdoutStepTitle.Printf("\n" + str + "\n" + strings.Repeat("=", len(str)) + "\n")
	s.report.Step(str)
}

func (s *Stdout) ModuleName(str string) {
	stdoutModuleName.Println("  " + str)
	s.report.Module(str)
}

func (s *Stdout) Info(str string) {
	color.Info.Println("    ⚠  " + str)
	s.report.Info(str)
}

func (s *Stdout) Success(str string) {
	color.Info.Println("    ✓  " + str)
	s.report.Success(str)
}

func (s *Stdout) Alert(str string) {
	color.Danger.Println("    ⊖  " + str)
	s.report.Alert(str)
}

func (s *Stdout) Error(err error) {
	fmt.Print("     ")
	color.Error.Println(" " + err.Error() + " ")
	s.report.Error(err.Error())
}

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

func (s *Stdout) HideLoader() {
	s.stopLoading <- true
	fmt.Print(loaderClear)
}

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

func (s *Stdout) HidePercentage() {
	fmt.Print(loaderClear)
}

func (s *Stdout) Report() {
	fmt.Println("")
	stdoutReportTitle.Println(" Report ")
	stdoutReportTitle.Println(" ====== ")
	alerts, errors := s.report.AlertsAndErrors()

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
			for _, module := range step.modules {
				stdoutModuleName.Println("      " + module.name)
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
			for _, module := range step.modules {
				stdoutModuleName.Println("      " + module.name)
				for _, err := range module.errors {
					color.Danger.Println("        " + err)
				}
			}
		}
	}

}
