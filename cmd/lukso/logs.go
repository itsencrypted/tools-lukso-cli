package main

import (
	"bufio"
	"bytes"
	"github.com/urfave/cli/v2"
	"os"
	"os/exec"
	"strings"
)

func (dependency *ClientDependency) Log(logFilePath string) (err error) {
	var commandName string
	var commandArgs []string
	switch systemOs {
	case ubuntu, macos:
		commandName = "tail"
		commandArgs = []string{"-f", "-n", "+0"}
	case windows:
		commandName = "type"
	default:
		commandName = "tail" // For reviewers - do we provide default command? Or omit and return with err?
		commandArgs = []string{"-f", "-n", "+0"}
	}

	command := exec.Command(commandName, append(commandArgs, logFilePath)...)

	command.Stdout = os.Stdout

	err = command.Run()
	if _, ok := err.(*exec.ExitError); ok {
		log.Error("No error logs found")

		return
	}

	// error unrelated to command execution
	if err != nil {
		log.Errorf("There was an error while executing command: %s. Error: %v", commandName, err)
	}

	return
}

// Stat returns whether the client is running or not
func (dependency *ClientDependency) Stat() (isRunning bool, err error) {
	var (
		commandName string
		buf         = new(bytes.Buffer)
	)

	isRunning = false

	switch systemOs {
	case ubuntu, macos:
		commandName = "ps"
	case windows:
		commandName = "tasklist"
	default:
		commandName = "ps"
	}

	command := exec.Command(commandName)
	command.Stdout = buf

	err = command.Run()
	if err != nil {
		log.Errorf("There was an error while executing command: %s. Error: %v", commandName, err)

		return
	}

	scan := bufio.NewScanner(buf)
	for scan.Scan() {
		if strings.Contains(scan.Text(), dependency.name) {
			isRunning = true

			return
		}
	}

	return
}

func logClients(ctx *cli.Context) error {
	log.Info("Please specify your client - run lukso logs help for more info")

	return nil
}

func logClient(dependencyName string, logFileFlag string) func(*cli.Context) error {
	return func(ctx *cli.Context) error {
		logFileDir := ctx.String(logFileFlag)
		if logFileDir == "" {
			return errFlagMissing
		}

		latestFile, err := getLastFile(logFileDir, dependencyName)
		if latestFile == "" && err == nil {
			return nil
		}

		if err != nil {
			return err
		}

		return clientDependencies[dependencyName].Log(logFileDir + "/" + latestFile)
	}
}

func statClients(ctx *cli.Context) (err error) {
	err = statClient(gethDependencyName)(ctx)
	if err != nil {
		return
	}

	err = statClient(prysmDependencyName)(ctx)
	if err != nil {
		return
	}

	err = statClient(validatorDependencyName)(ctx)
	if err != nil {
		return
	}

	return
}

func statClient(dependencyName string) func(*cli.Context) error {
	return func(ctx *cli.Context) error {
		isRunning, err := clientDependencies[dependencyName].Stat()
		if err != nil {
			return err
		}

		if isRunning {
			log.Infof("%s: Running", dependencyName)

			return nil
		}

		log.Warnf("%s: Stopped", dependencyName)

		return nil
	}
}
