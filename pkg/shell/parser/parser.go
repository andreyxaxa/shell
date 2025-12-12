package parser

import (
	"errors"
	"strings"
)

// ParseArgs parses arguments and returns a slice of arguments
// "echo hello" -> ["echo", "hello"]
func ParseArgs(line string) ([]string, error) {
	var args []string
	var s strings.Builder
	quoted := rune(0)
	escaped := false

	for _, r := range line {
		// Если экранирование - записываем как есть
		if escaped {
			s.WriteRune(r)
			escaped = false
			continue
		}
		// Если встретили экранирование - помечаем
		if r == '\\' {
			escaped = true
			continue
		}
		// Если мы до этого встретили кавычку
		if quoted != 0 {
			// И она такого же вида
			if r == quoted {
				// То закрываем кавычки
				quoted = 0
				// Если другой символ - записывем его
			} else {
				s.WriteRune(r)
			}
			continue
		}

		// Если встретили кавычку
		if r == '"' || r == '\'' {
			quoted = r
			continue
		}

		// Если встретили пробел/табуляцию - добавляем предыдущее слово как аргумент
		if r == ' ' || r == '\t' {
			if s.Len() > 0 {
				args = append(args, s.String())
				s.Reset()
			}
			continue
		}
		s.WriteRune(r)
	}

	// Если прошли все условия, но остался символ экранирования
	if escaped {
		return nil, errors.New("unfinished scope")
	}
	// Если прошли все условия, но остались незакрытые кавычки
	if quoted != 0 {
		return nil, errors.New("unfinished quote")
	}
	if s.Len() > 0 {
		args = append(args, s.String())
	}

	return args, nil
}
