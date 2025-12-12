package helpers

import (
	"errors"
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/andreyxaxa/shell/pkg/shell/builtins"
	"github.com/andreyxaxa/shell/pkg/shell/executors"
	"github.com/andreyxaxa/shell/pkg/shell/executors/pipeline"
	"github.com/andreyxaxa/shell/pkg/shell/parser"
)

// TakeArgsFromCmd takes string of commands and returns a slice of commands, input and output files.
// " cd 22" -> "cd 22" -> ["cd", "22"]
func TakeArgsFromCmd(idx int, line string) ([]string, string, string, error) {
	line = strings.TrimSpace(line)
	if line == "" {
		return nil, "", "", errors.New("empty command in chain")
	}

	if strings.Contains(line, "|") {
		if err := pipeline.Run(line); err != nil {
			return nil, "", "", fmt.Errorf("pipeline failed in conditional execution: %w", err)
		}
		return nil, "", "", nil
	}

	var args []string
	inputFile, outputFile := "", ""
	var err error

	args, err = parser.ParseArgs(line)
	if err != nil {
		return nil, "", "", fmt.Errorf("parse error in command %d: %w", idx+1, err)
	}

	if slices.Contains(args, "<") || slices.Contains(args, ">") {
		args, inputFile, outputFile, err = parser.ParseRedirects(args)
	}

	args = parser.ParseEnv(args)

	if len(args) == 0 {
		return nil, "", "", fmt.Errorf("empty command at position %d", idx+1)
	}

	return args, inputFile, outputFile, err
}

// CheckCmd checks whether a command is builtin or external.
// If this is builtin -> we call it, else we call it via exec.Command.Run()
func CheckCmd(args []string, inputFile, outputFile string) error {
	var lastErr error

	switch args[0] {
	case "cd":
		lastErr = builtins.Cd(args)
		if lastErr != nil {
			fmt.Fprintln(os.Stderr, "cd:", lastErr)
		}
	case "pwd":
		lastErr = builtins.Pwd(outputFile)
		if lastErr != nil {
			fmt.Fprintln(os.Stderr, "pwd:", lastErr)
		}
	case "echo":
		builtins.Echo(args, outputFile)
	case "ps":
		builtins.Ps(outputFile)
	case "kill":
		builtins.Kill(args)
	default:
		lastErr = executors.RunWithExitCode(args, inputFile, outputFile)
	}

	return lastErr
}
