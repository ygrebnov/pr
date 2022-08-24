package pipelinerunner

import (
	"github.com/google/uuid"
)

type Pipeline struct {
	ID       string
	Name     string
	Actions  []*Action
	Modified string
}

func NewPipeline(name string, actions []*Action) *Pipeline {
	p := new(Pipeline)
	p.ID = uuid.New().String()
	p.Name = name
	p.Actions = actions
	return p
}
