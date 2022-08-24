package pipelinerunner

import (
	"bytes"
	"fmt"
	"text/template"
)

// A type to identify Actions by their type
type ActionType uint8

const (
	GIT ActionType = iota
	ENV
	RUN
	INVALID
)

func actionTypeFromString(s string) (ActionType, error) {
	switch s {
	case "GIT":
		return GIT, nil
	case "ENV":
		return ENV, nil
	case "RUN":
		return RUN, nil
	default:
		return INVALID, fmt.Errorf("unknown Action type: %s", s)
	}
}

func (a ActionType) String() string {
	return []string{"GIT", "ENV", "RUN"}[a]
}

func (a ActionType) Command(attr map[string]string) string {
	tmpl := []string{
		"git clone {{if .Branch}}-b {{.Branch}}{{end}} {{.Base}} {{if .Dest}}{{.Dest}}{{end}}",
		"export {{.Base}}{{if .Value}}={{.Value}}{{end}}",
		"{{.Base}}",
	}[a]
	command_tmpl := template.Must(template.New("command_tmpl").Parse(tmpl))
	var buff bytes.Buffer
	if err := command_tmpl.Execute(&buff, attr); err != nil {
		panic(fmt.Sprintf("Cannot create command from template: %v", err))
	}
	return buff.String()
}

func (a ActionType) Attributes() []string {
	return [][]string{
		{"BRANCH", "DEST", "ASYNC"},
		{"VALUE", "ASYNC"},
		{"ASYNC"},
	}[a]
}

type Action struct {
	ID          string     `db:"id"`
	Pipelineref string     `db:"pipelineref"`
	Actiontype  ActionType `db:"actiontype"`
	Async       bool       `db:"async"` // A flag to indicate whether the action can be executed asynchronously
	Command     string     `db:"command"`
	Modified    string     `db:"modified"`
}

// Compare an Action with other Action
func (a Action) Equal(other Action) bool {
	return a.Actiontype == other.Actiontype && a.Async == other.Async && a.Command == other.Command
}

func NewAction(at ActionType, async bool, attr map[string]string) *Action {
	a := new(Action)
	a.ID = getUUID()
	a.Actiontype = at
	a.Async = async
	a.Command = a.Actiontype.Command(attr)
	return a
}
