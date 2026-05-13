package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/real-uangi/opencode-orc/buffer"
	"github.com/real-uangi/opencode-orc/config"
	"github.com/real-uangi/opencode-orc/filter"
	"github.com/real-uangi/opencode-orc/output"
	"github.com/real-uangi/opencode-orc/parser"
	"github.com/real-uangi/opencode-orc/process"
)

func main() {
	configPath := flag.String("config", config.DefaultConfigPath(), "Path to config file")
	showVersion := flag.Bool("version", false, "Show version information")
	flag.Parse()

	if *showVersion {
		fmt.Println("opencode-orc version 0.1.0")
		os.Exit(0)
	}

	if err := config.WriteDefaultConfig(*configPath); err != nil {
		output.WriteError("failed to create default config: %v", err)
		os.Exit(1)
	}

	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		output.WriteError("failed to load config: %v", err)
		os.Exit(1)
	}

	args := flag.Args()
	if len(args) == 0 {
		output.WriteError("usage: opencode-orc [flags] <prompt>")
		os.Exit(1)
	}

	opencodeArgs := []string{"run", "--format", "json"}
	opencodeArgs = append(opencodeArgs, args...)

	proc, err := process.New(opencodeArgs)
	if err != nil {
		output.WriteError("failed to create process: %v", err)
		os.Exit(2)
	}

	if err := proc.Start(); err != nil {
		output.WriteError("failed to start process: %v", err)
		os.Exit(2)
	}

	p := parser.NewParser(proc.Stdout())
	f := filter.New(cfg)
	b := buffer.New()
	w := output.New(os.Stdout, cfg.Output.Format, cfg.Output.Pretty)

	for {
		event, err := p.ParseNext()
		if err == io.EOF {
			break
		}
		if err != nil {
			output.WriteError("failed to parse event: %v", err)
			continue
		}

		events := f.Filter(event)
		if len(events) == 0 {
			continue
		}

		for _, filtered := range events {
			b.Update(filtered)

			if err := w.WriteEvent(filtered); err != nil {
				output.WriteError("failed to write event: %v", err)
				os.Exit(3)
			}
		}
	}

	if err := proc.Wait(); err != nil {
		output.WriteError("process failed: %v", err)
		b.UpdateFromError("ProcessError", err.Error())
	}

	doneEvent := b.DoneEvent()
	if err := w.WriteEvent(doneEvent); err != nil {
		output.WriteError("failed to write done event: %v", err)
		os.Exit(3)
	}

	os.Exit(proc.ExitCode())
}
