package builtins

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Echo prints text into os.Stdout or into file if it exists
func Echo(args []string, outputFile string) {
	if len(args) <= 1 {
		fmt.Println()
		return
	}

	output := strings.Join(args[1:], " ")

	if outputFile != "" {
		file, err := os.Create(outputFile)
		if err != nil {
			fmt.Fprintln(os.Stderr, "echo: cannot create file:", err)
			return
		}
		defer file.Close()
		_, err = file.WriteString(output + "\n")
		if err != nil {
			fmt.Fprintln(os.Stderr, "echo: cannot write to file:", err)
		}
	} else {
		fmt.Println(output)
	}
}

// Pwd prints absolute path of current directory.
func Pwd(outputFile string) error {
	dir, err := os.Getwd()
	if err != nil {
		return err
	}

	if outputFile != "" {
		file, err := os.Create(outputFile)
		if err != nil {
			return err
		}
		defer file.Close()
		_, err = file.WriteString(dir + "\n")
		return err
	}

	fmt.Println(dir)
	return nil
}

// Cd changes the current working directory to the named directory
func Cd(args []string) error {
	dir := ""

	if len(args) < 2 {
		home := os.Getenv("HOME")
		if home == "" {
			return errors.New("HOME not set")
		}
		dir = home
	} else {
		dir = args[1]
	}

	return os.Chdir(dir)
}

// Ps shows running processes (ps -ef)
func Ps(outputFile string) {
	cmd := exec.Command("ps", "-ef")
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr

	if outputFile != "" {
		file, err := os.Create(outputFile)
		if err != nil {
			fmt.Fprintln(os.Stderr, "ps: cannot create file:", err)
			return
		}
		defer file.Close()
		cmd.Stdout = file
	} else {
		cmd.Stdout = os.Stdout
	}

	if err := cmd.Run(); err != nil {
		fmt.Fprintln(os.Stderr, "ps:", err)
	}
}

// Kill causes the process to exit immediately
func Kill(args []string) {
	if len(args) < 2 {
		fmt.Fprintln(os.Stderr, "kill: missing pid")
		return
	}

	pid, err := strconv.Atoi(args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, "invalid pid:", err)
		return
	}

	process, err := os.FindProcess(pid)
	if err != nil {
		fmt.Fprintln(os.Stderr, "kill: process not found:", err)
		return
	}

	if err := process.Kill(); err != nil {
		fmt.Fprintln(os.Stderr, "kill:", err)
		return
	}
}
