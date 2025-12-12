package redirects

import (
	"fmt"
	"os"
	"os/exec"
)

// RunWithRedirects runs external command with redirects
func RunWithRedirects(args []string, inputFile, outputFile string) error {
	cmd := exec.Command(args[0], args[1:]...)

	if inputFile != "" {
		input, err := os.Open(inputFile)
		if err != nil {
			return fmt.Errorf("cannot open input file %s: %w", inputFile, err)
		}
		defer input.Close()
		cmd.Stdin = input
	} else {
		cmd.Stdin = os.Stdin
	}

	if outputFile != "" {
		output, err := os.Create(outputFile)
		if err != nil {
			return fmt.Errorf("cannot create/open file %s: %w", outputFile, err)
		}
		defer output.Close()
		cmd.Stdout = output
	} else {
		cmd.Stdout = os.Stdout
	}

	return cmd.Run()
}
