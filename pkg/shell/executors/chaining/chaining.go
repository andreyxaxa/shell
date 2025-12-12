package chaining

import (
	"errors"
	"fmt"
	"strings"

	"github.com/andreyxaxa/shell/pkg/shell/executors/chaining/helpers"
)

// Для ||
func RunCondOr(line string) error {
	commands := strings.Split(line, "||")
	if len(commands) == 0 {
		return errors.New("empty command line")
	}

	for i, cmdStr := range commands {
		args, inputFile, outputFile, err := helpers.TakeArgsFromCmd(i, cmdStr)
		if err != nil {
			return fmt.Errorf("parse error in command %d: %w", i+1, err)
		}

		lastErr := helpers.CheckCmd(args, inputFile, outputFile)

		if lastErr == nil {
			break
		} else {
			fmt.Printf("command failed: %s\n", cmdStr)
			continue
		}
	}

	return nil
}

// Для &&
func RunCondAnd(line string) error {
	commands := strings.Split(line, "&&")
	if len(commands) == 0 {
		return errors.New("empty command line")
	}

	for i, cmdStr := range commands {
		args, inputFile, outputFile, err := helpers.TakeArgsFromCmd(i, cmdStr)
		if err != nil {
			return fmt.Errorf("parse error in command %d: %w", i+1, err)
		}

		lastErr := helpers.CheckCmd(args, inputFile, outputFile)

		if lastErr != nil {
			return fmt.Errorf("command failed: %s", cmdStr)
		}
	}

	return nil
}
