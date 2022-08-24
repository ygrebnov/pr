package pipelinerunner

import (
	"bytes"
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"testing"
)

func TestCreateAction(t *testing.T) {

	testActionService := NewActionService()

	var tests = []struct {
		input string
		want  Action
	}{
		{
			"RUN echo \"What a lovely day!\" > whatalovalyday",
			Action{Command: "echo \"What a lovely day!\" > whatalovalyday", Async: false, Actiontype: RUN},
		},
		{
			"RUN echo \"What a lovely day!\" > whatalovalyday ASYNC",
			Action{Command: "echo \"What a lovely day!\" > whatalovalyday", Async: true, Actiontype: RUN},
		},
		{
			"GIT https://github.com/google/re2.git BRANCH main DEST $GOPATH/github.com/google ASYNC",
			Action{
				Command:    "git clone -b main https://github.com/google/re2.git $GOPATH/github.com/google",
				Async:      true,
				Actiontype: GIT,
			},
		},
	}

	for _, tt := range tests {
		testname := fmt.Sprintf("%s,%v", tt.want.Actiontype.String(), tt.want.Async)
		t.Run(testname, func(t *testing.T) {
			a, _ := testActionService.CreateAction(tt.input)

			if !a.Equal(tt.want) {
				t.Errorf("got: %v, want: %v", *a, tt.want)
			}
		})
	}
}

func TestParseStringAction(t *testing.T) {
	var tests = []struct {
		name       string
		actionType ActionType
		input      string
		want       map[string]string
	}{
		{
			"run",
			RUN,
			"echo \"What a lovely day!\" > whatalovelyday",
			map[string]string{"Base": "echo \"What a lovely day!\" > whatalovelyday"},
		},
		{
			"async_run",
			RUN,
			"echo \"What a lovely day!\" > whatalovelyday ASYNC",
			map[string]string{"Base": "echo \"What a lovely day!\" > whatalovelyday"},
		},
		{
			"async_run2",
			RUN,
			"touch step2; sleep 2; echo \"step2 async, after sleeping for 2s\" > step2 ASYNC",
			map[string]string{
				"Base": "touch step2; sleep 2; echo \"step2 async, after sleeping for 2s\" > step2",
			},
		},
		{
			"async_git_with_branch_dest",
			GIT,
			"https://github.com/google/re2.git BRANCH main DEST $GOPATH/github.com/google ASYNC",
			map[string]string{
				"Base":   "https://github.com/google/re2.git",
				"Branch": "main",
				"Dest":   "$GOPATH/github.com/google",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			re := regexp.MustCompile("^(.*?)(" + strings.Join(tt.actionType.Attributes(), "|") + ")")
			attr := make(map[string]string)
			parseActionString(tt.input, re, "Base", attr)
			if !reflect.DeepEqual(attr, tt.want) {
				t.Errorf("got: %v, want: %v", attr, tt.want)
			}
		})
	}
}

func TestExecuteAction(t *testing.T) {
	testActionService := NewActionService()

	a := Action{Command: "echo hello world", Async: false, Actiontype: RUN}

	var aOut, aErr bytes.Buffer

	if err := testActionService.ExecuteAction(&a, &aOut, &aErr); err != nil {
		t.Errorf("error executing Action: %v", err)
	}

	if strings.TrimSpace(aOut.String()) != "hello world" {
		t.Errorf("incorrect stdout, got: [%s], want: [%s]", aOut.String(), "hello world")
	}

	if strings.TrimSpace(aErr.String()) != "" {
		t.Errorf("incorrect stderr, got: [%s], want: [%s]", aErr.String(), "")
	}

}
