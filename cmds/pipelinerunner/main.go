package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"pipelinerunner"
)

func main() {
	newCmd := flag.NewFlagSet("new", flag.ExitOnError)
	newPfilePath := newCmd.String("file", "", "file")
	newName := newCmd.String("name", "", "name")

	if len(os.Args) < 2 {
		log.Fatal(pipelinerunner.USAGE.PR)
	}

	if err := pipelinerunner.Configure(); err != nil {
		log.Fatal("error configuring application")
	}

	pipelineService := pipelinerunner.NewPipelineService()
	pipelineRunService := pipelinerunner.NewPipelineRunService()

	switch os.Args[1] {
	case "help":
		switch len(os.Args) {
		case 3:
			fmt.Println(pipelinerunner.GetUsage(os.Args[2]))
		default:
			fmt.Println(pipelinerunner.GetUsage(""))
		}
	case "config":
		if len(os.Args) < 3 {
			log.Fatal("please specify the subcommand. Use \"pr help config\" for more information.")
		}
		switch os.Args[2] {
		case "view":
			fmt.Printf("%q\n", pipelinerunner.Config)
		case "set":
			if len(os.Args) < 4 {
				log.Fatal("please specify the configuration file path. Command should have syntax: pr config set <file path>")
			}
			if err := pipelinerunner.ParseConfiguration(os.Args[3]); err != nil {
				log.Fatalf("cannot parse configuration from file: %s", os.Args[3])
			}
			fmt.Printf("Applied configuration from file: %s", os.Args[3])
		default:
			log.Fatal(pipelinerunner.USAGE.CONFIG)
		}
	case "ls":
		if err := pipelineService.PrintPipelines(); err != nil {
			log.Fatal("error fetching pipelines data")
		}
	case "new":
		newCmd.Parse(os.Args[2:])
		pipelineId, err := pipelineService.CreatePipeline(*newPfilePath, *newName)
		if err != nil {
			log.Fatal("error creating a new pipeline")
		}
		fmt.Printf("Created a new pipeline: %s\n", pipelineId)
	case "view":
		if len(os.Args) < 3 {
			log.Fatal("please specify the pipeline id. Command should have syntax: pr view <pipeline id>")
		}
		pipelineId := os.Args[2]
		if err := pipelineService.PrintPipeline(pipelineId); err != nil {
			log.Fatalf("error getting data for pipeline id: %s", pipelineId)
		}
	case "rm":
		if len(os.Args) < 3 {
			log.Fatal("please specify the pipeline id. Command should have syntax: pr rm <pipeline id>")
		}
		pipelineId := os.Args[2]
		err := pipelineService.DeletePipeline(pipelineId)
		if err != nil {
			log.Fatalf("error deleting pipeline with id: %s", pipelineId)
		}
		fmt.Printf("Deleted pipeline with id: %s\n", pipelineId)
	case "run":
		if len(os.Args) < 3 {
			log.Fatal("please specify the pipeline id. Command should have syntax: pr run <pipeline id>")
		}
		pipelineId := os.Args[2]
		err := pipelineService.RunPipeline(pipelineId)
		if err != nil {
			log.Fatalf("error running pipeline with id: %s", pipelineId)
		}
		fmt.Printf("Running pipeline with id: %s\nIts status can be checked using the log command: pr log %s\n", pipelineId, pipelineId)
	case "log":
		switch len(os.Args) {
		case 3:
			if err := pipelineRunService.PrintPipelineRuns(os.Args[2]); err != nil {
				log.Fatalf("error getting runs data for pipeline id: %s", os.Args[2])
			}
		default:
			if err := pipelineRunService.PrintAllPipelineRuns(); err != nil {
				log.Fatal("error getting pipeline runs data")
			}
		}

	default:
		log.Fatalf("unknown command\n%s", pipelinerunner.USAGE.PR)
	}
}
