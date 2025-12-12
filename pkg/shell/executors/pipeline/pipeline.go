package pipeline

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/andreyxaxa/shell/pkg/shell/parser"
)

// Run splits string into commands by "|" and connects them into pipe
func Run(line string) error {
	// Разбиваем строку на команды по символу |
	// ps | grep go -> ["ps ", " grep go"]
	commands := strings.Split(line, "|")
	if len(commands) == 0 {
		return errors.New("empty pipeline")
	}

	var cmdList []*exec.Cmd
	var pipes []*os.File

	for i, cmdStr := range commands {
		// ["ps ", " grep go"] -> ["ps", "grep go"]
		cmdStr = strings.TrimSpace(cmdStr)
		if cmdStr == "" {
			return errors.New("empty command in pipeline")
		}

		// парсим
		args, err := parser.ParseArgs(cmdStr)
		if err != nil {
			return fmt.Errorf("parse error in command %d: %w", i+1, err)
		}

		args = parser.ParseEnv(args)

		if len(args) == 0 {
			return fmt.Errorf("empty command at position %d", i+1)
		}

		// создаем команды
		cmd := exec.Command(args[0], args[1:]...)
		cmdList = append(cmdList, cmd)
	}

	// Связываем команды через пайпы
	for i := 0; i < len(cmdList)-1; i++ {
		// Создаем пайп между текущей и следующей командой
		r, w, err := os.Pipe()
		if err != nil {
			return fmt.Errorf("pipe creation failed: %w", err)
		}

		// Стандартный вывод текущей команды направляем в пайп
		cmdList[i].Stdout = w
		// Стандартный ввод следующей команды берем из пайпа
		cmdList[i+1].Stdin = r

		// Сохраняем writer-части для закрытия после запуска
		pipes = append(pipes, w)
	}

	// Первая команда получает ввод из stdin, последняя - вывод в stdout
	cmdList[0].Stdin = os.Stdin
	cmdList[len(cmdList)-1].Stdout = os.Stdout

	// Все команды пишут ошибки в stderr
	for _, cmd := range cmdList {
		cmd.Stderr = os.Stderr
	}

	// Запускаем все команды
	for _, cmd := range cmdList {
		if err := cmd.Start(); err != nil {
			return fmt.Errorf("failed to start command: %w", err)
		}
	}

	// Закрываем writer-части пайпов чтобы команды могли завершиться
	for _, pipe := range pipes {
		pipe.Close()
	}

	// Ждем завершения всех команд
	for _, cmd := range cmdList {
		if err := cmd.Wait(); err != nil {
			if exitError, ok := err.(*exec.ExitError); ok {
				_ = exitError
			} else {
				return fmt.Errorf("command wait failed: %w", err)
			}
		}
	}

	return nil
}
