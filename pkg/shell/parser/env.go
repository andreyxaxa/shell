package parser

import (
	"os"
	"strings"
)

// Смотрим, есть ли переменные окружения
func ParseEnv(args []string) []string {
	var enved []string

	for _, arg := range args {
		envedArg := parseEnvArg(arg)
		enved = append(enved, envedArg)
	}

	return enved
}

// Если аргумент содержит "$", то пробуем достать переменную, иначе возвращаем аргумент как есть
func parseEnvArg(arg string) string {
	if !strings.Contains(arg, "$") {
		return arg
	}

	argLen := len(arg)

	if argLen == 1 {
		return "$"
	}

	var res strings.Builder

	for i := 0; i < argLen; i++ {
		if arg[i] == '$' && i+1 < argLen {
			name := arg[i+1 : argLen]
			value := os.Getenv(name)
			res.WriteString(value)
			return res.String()
		}
		res.WriteByte(arg[i])
	}

	return res.String()
}
