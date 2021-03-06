package main

import (
	"io"
	"os"
	"strings"

	"github.com/creativeprojects/clog"
	"github.com/creativeprojects/resticprofile/shell"
)

type shellCommandDefinition struct {
	command  string
	args     []string
	env      []string
	useStdin bool
	stdout   io.Writer
	stderr   io.Writer
	dryRun   bool
	sigChan  chan os.Signal
}

// newShellCommand creates a new shell command definition
func newShellCommand(command string, args, env []string, dryRun bool, sigChan chan os.Signal) shellCommandDefinition {
	if env == nil {
		env = make([]string, 0)
	}
	return shellCommandDefinition{
		command:  command,
		args:     args,
		env:      env,
		useStdin: false,
		stdout:   os.Stdout,
		stderr:   os.Stderr,
		dryRun:   dryRun,
		sigChan:  sigChan,
	}
}

// runShellCommand instantiates a shell.Command and sends the information to run the shell command
func runShellCommand(command shellCommandDefinition) error {
	var err error

	if command.dryRun {
		clog.Infof("dry-run: %s %s", command.command, strings.Join(command.args, " "))
		return nil
	}

	shellCmd := shell.NewSignalledCommand(command.command, command.args, command.sigChan)

	shellCmd.Stdout = command.stdout
	shellCmd.Stderr = command.stderr

	if command.useStdin {
		shellCmd.Stdin = os.Stdin
	}

	shellCmd.Environ = os.Environ()
	if command.env != nil && len(command.env) > 0 {
		shellCmd.Environ = append(shellCmd.Environ, command.env...)
	}

	err = shellCmd.Run()
	if err != nil {
		return err
	}
	return nil
}
