package executors

import (
	"fmt"
	"os"
	"os/exec"
)

// RunWithExitCode runs external command and RETURNS ERROR if return code != 0.
// Used in command chaining
func RunWithExitCode(args []string, inputFile, outputFile string) error {
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if inputFile != "" {
		file, err := os.Create(inputFile)
		if err != nil {
			return fmt.Errorf("cannot create or open file %s: %w", outputFile, err)
		}
		defer file.Close()
		cmd.Stdin = file
	} else {
		cmd.Stdin = os.Stdin
	}

	if outputFile != "" {
		file, err := os.Create(outputFile)
		if err != nil {
			return fmt.Errorf("cannot create or open file %s: %w", outputFile, err)
		}
		defer file.Close()
		cmd.Stdout = file
	} else {
		cmd.Stdout = os.Stdout
	}

	return cmd.Run()
}

// Run runs external command and DISPLAYS ERROR without returning it.
// Because we use Run directly in shell, just displays error to avoid shell terminating
func Run(args []string) {
	if len(args) == 0 {
		return
	}

	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout

	if err := cmd.Run(); err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			_ = exitError
		} else {
			fmt.Fprintln(os.Stderr, "exec error:", err)
		}
	}
}
