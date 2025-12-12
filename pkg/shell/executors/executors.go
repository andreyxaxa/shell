package executors

import (
	"fmt"
	"os"
	"os/exec"
)

// runWithExitCode выполняет внешнюю команду и возвращает ошибку если код возврата != 0
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

// Используем в самой оболочке - main, поэтому просто показываем ошибку, не возвращая ее, чтобы не завершать main.
func Run(args []string) {
	if len(args) == 0 {
		return
	}

	// TODO: отладочная инфа, потом удалить
	fmt.Println("run", args)

	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout

	if err := cmd.Run(); err != nil {
		// TODO: подумать над участком кода
		if exitError, ok := err.(*exec.ExitError); ok {
			_ = exitError
		} else {
			fmt.Fprintln(os.Stderr, "exec error:", err)
		}
	}
}
