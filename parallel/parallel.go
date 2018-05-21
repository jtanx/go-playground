package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

type Settings struct {
	Jobs       int
	Executable string
	StaticArgs []string
	Inputs     []string

	Verbose bool
}

type JobState struct {
	Log      Logger
	Settings Settings
	Jobs     chan string
	JobErrs  chan error // nil error == goroutine done
}

type Logger struct {
	File    *os.File
	Verbose bool
}

func (l *Logger) Log(prefix, format string, args ...interface{}) {
	tm := time.Now().Format("2006-01-02 15:04:05.000 ")
	fmt.Fprintf(l.File, tm+prefix+format+"\n", args...)
}
func (l *Logger) Debug(format string, args ...interface{}) {
	if l.Verbose {
		l.Log("[Debug ] ", format, args...)
	}
}
func (l *Logger) Info(format string, args ...interface{})  { l.Log("[Info  ] ", format, args...) }
func (l *Logger) Error(format string, args ...interface{}) { l.Log("[Error ] ", format, args...) }

func parseCmdline(settings *Settings) error {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [flags] executable <[static args...] -- [input inputs...]>|<[input inputs...]>\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.IntVar(&settings.Jobs, "j", 2, "Number of parallel jobs")
	flag.BoolVar(&settings.Verbose, "v", false, "Verbose")
	flag.Parse()

	args := flag.Args()
	if settings.Jobs <= 0 {
		return fmt.Errorf("Invalid job count: %d", settings.Jobs)
	} else if len(args) <= 0 {
		return fmt.Errorf("Need an executable")
	}

	settings.Executable = args[0]
	inputs := args[1:]
	for i, arg := range inputs {
		if arg == "--" {
			settings.StaticArgs = inputs[:i]
			inputs = inputs[i+1:]
			break
		}
	}

	for _, input := range inputs {
		files, err := filepath.Glob(input)
		if err != nil {
			return err
		} else if len(files) == 0 {
			settings.Inputs = append(settings.Inputs, input)
		} else {
			settings.Inputs = append(settings.Inputs, files...)
		}
	}

	if len(settings.Inputs) > 0 && settings.Jobs > len(settings.Inputs) {
		settings.Jobs = len(settings.Inputs)
	}

	return nil
}

func worker(id int, state *JobState) {
	state.Log.Debug("Worker %d init", id)

	for j := range state.Jobs {
		args := append(state.Settings.StaticArgs, j)
		state.Log.Debug("Worker %d: Job %s %v", id, state.Settings.Executable, args)
		cmd := exec.Command(state.Settings.Executable, args...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			state.JobErrs <- fmt.Errorf("%s %v: %v", state.Settings.Executable, args, err)
		}
	}
	state.Log.Debug("Worker %d done", id)
	state.JobErrs <- nil
}

func main() {
	state := JobState{Log: Logger{File: os.Stderr}}
	if err := parseCmdline(&state.Settings); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		flag.Usage()
		os.Exit(1)
	}

	state.Log.Verbose = state.Settings.Verbose
	state.Log.Info("Settings: %#v", state.Settings)

	state.Jobs = make(chan string, state.Settings.Jobs)
	state.JobErrs = make(chan error)
	for j := 0; j < state.Settings.Jobs; j++ {
		go worker(j+1, &state)
	}

	for _, input := range state.Settings.Inputs {
		state.Jobs <- input
	}
	close(state.Jobs)

	var errors []error
	for state.Settings.Jobs > 0 {
		err := <-state.JobErrs
		if err == nil {
			state.Settings.Jobs--
		} else {
			errors = append(errors, err)
		}
	}

	if len(errors) > 0 {
		state.Log.Error("The following jobs exited with an error:")
		for _, err := range errors {
			state.Log.Error("%v", err)
		}
		os.Exit(len(errors))
	}
}
