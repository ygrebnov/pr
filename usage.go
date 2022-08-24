package pipelinerunner

type Usage struct {
	PR     string
	CONFIG string
}

var USAGE Usage = Usage{
	PR: `
Pipelinerunner is a tool to run pipelines defined as sequences of actions in pipeline files (pfiles).

Usage: 
	pr <command> [arguments]

The commands are:

	ls                      list pipelines
	new                     create a new pipeline from pfile. Example: pr new --name Example --file ./pfiles/example0.pfile
	view                    view pipeline content. Example: pr view <pipeline id>
	rm                      delete pipeline. Example: pr rm <pipeline id>
	run                     run pipeline. Example: pr run <pipeline id>
	log                     view pipeline runs log. For all pipelines: pr log. For some particular pipeline: pr log <pipeline id>
	config                  configuration commands. Example: pr config view
`,
	CONFIG: `
usage: pr config [subcommand] [arguments]

The subcommands are:

	set <file path>			set configuration from path. See conf.json for an example
	view                    view current configuration
`,
}

func GetUsage(s string) string {
	switch s {
	case "config":
		return USAGE.CONFIG
	default:
		return USAGE.PR
	}
}
