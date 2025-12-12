package shell

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"os/signal"
	"slices"
	"strings"
	"syscall"

	"github.com/andreyxaxa/shell/pkg/shell/builtins"
	"github.com/andreyxaxa/shell/pkg/shell/executors"
	"github.com/andreyxaxa/shell/pkg/shell/executors/chaining"
	"github.com/andreyxaxa/shell/pkg/shell/executors/pipeline"
	"github.com/andreyxaxa/shell/pkg/shell/executors/redirects"
	"github.com/andreyxaxa/shell/pkg/shell/parser"
)

func Start() {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT)
	go func() {
		for range sig {
			fmt.Println("^C")
		}
	}()

	reader := bufio.NewReader(os.Stdin)

	for {
		wd, _ := os.Getwd()
		fmt.Printf("desktop minishell %s\n$ ", wd)

		line, err := reader.ReadString('\n')
		if err != nil {
			if errors.Is(err, io.EOF) {
				return
			}
			fmt.Fprintln(os.Stderr, "error while reading args:", err)
			continue
		}
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Проверяем, есть ли оператор &&
		if strings.Contains(line, "||") {
			if err := chaining.RunCondOr(line); err != nil {
				fmt.Fprintln(os.Stderr, "conditional execution error:", err)
			}
			continue
		}

		// Проверяем, есть ли оператор &&
		if strings.Contains(line, "&&") {
			if err := chaining.RunCondAnd(line); err != nil {
				fmt.Fprintln(os.Stderr, "conditional execution error:", err)
			}
			continue
		}

		// Проверяем, есть ли пайплайны
		if strings.Contains(line, "|") {
			if err := pipeline.RunPipeline(line); err != nil {
				fmt.Fprintln(os.Stderr, "pipeline error:", err)
			}
			continue
		}

		args, err := parser.ParseArgs(line)
		if err != nil {
			fmt.Fprintln(os.Stderr, "error while parsing args:", err)
			continue
		}

		if len(args) == 0 {
			continue
		}

		args = parser.ParseEnv(args)

		// редирект ?
		inputStr, outputStr := "", ""
		redirect := slices.Contains(args, ">") || slices.Contains(args, "<")

		if redirect {
			args, inputStr, outputStr, err = parser.ParseRedirects(args)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				continue
			}
		}

		switch args[0] {
		case "cd":
			if err := builtins.Cd(args); err != nil {
				fmt.Fprintln(os.Stderr, "cd:", err)
			}
			continue
		case "pwd":
			if err := builtins.Pwd(outputStr); err != nil {
				fmt.Fprintln(os.Stderr, "pwd:", err)
			}
			continue
		case "echo":
			builtins.Echo(args, outputStr)
			continue
		case "ps":
			builtins.Ps(outputStr)
			continue
		case "kill":
			builtins.Kill(args)
			continue
		}
		if redirect {
			redirects.RunWithRedirects(args, inputStr, outputStr)
		} else {
			executors.Run(args)
		}
	}
}
