package pipelinerunner

import (
	"bytes"
	"fmt"
	"os/exec"
	"regexp"
	"strings"
)

const actionsTemplate = "ActionID\tType\tAsync\tCommand\n{{range .}}{{.ID}}\t{{.Actiontype}}\t{{.Async}}\t{{.Command}}\n{{end}}"

func NewActionService() ActionService {
	return defaultActionService{}
}

type defaultActionService struct{}

// Create a new Action from a Pfile line. Insert Action data into the database
func (defaultActionService) CreateAction(s string) (*Action, error) {
	l := strings.Split(s, " ")
	at, err := actionTypeFromString(l[0])
	if err != nil {
		return nil, fmt.Errorf("cannot infer Action type")
	}
	// Create regexp pattern using ActionType attributes
	re := regexp.MustCompile("^(.*?)(" + strings.Join(at.Attributes(), "|") + ")")
	attr := make(map[string]string)
	parseActionString(s[len(l[0])+1:], re, "Base", attr)
	if len(attr) > 0 {
		return NewAction(at, strings.Contains(s, "ASYNC"), attr), nil
	} else {
		return nil, fmt.Errorf("invalid Action command")
	}
}

// Execute given Action
func (defaultActionService) ExecuteAction(a *Action, aOut *bytes.Buffer, aErr *bytes.Buffer) error {
	cmd := exec.Command("sh", "-c", a.Command)
	cmd.Stdout = aOut
	cmd.Stderr = aErr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error executing Action: %v", err)
	}
	return nil
}

// Recursively get Action attributes from a Pfile line
func parseActionString(input string, re *regexp.Regexp, currAttr string, attr map[string]string) {
	if len(input) == 0 {
		return
	} else {
		s := re.FindStringSubmatch(input)
		if s != nil {
			l := strings.Split(s[0], " ")
			attr[currAttr] = strings.Join(l[:len(l)-1], " ")
			parseActionString(strings.TrimSpace(input[len(s[0]):]), re, strings.Title(strings.ToLower(l[len(l)-1])), attr)
		} else {
			attr[currAttr] = input
		}
	}
}
