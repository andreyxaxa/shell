package parser

import "errors"

// ["echo", "123", ">", "out.txt"]
// ["echo", "123"], inputFile, outputFile, nil
func ParseRedirects(args []string) ([]string, string, string, error) {
	var resArgs []string
	inputFile := ""
	outputFile := ""

	for i := 0; i < len(args); i++ {
		if args[i] == "<" {
			if i+1 >= len(args) {
				return nil, "", "", errors.New("no input file after <")
			}
			inputFile = args[i+1]
			i += 2
			continue
		}

		if args[i] == ">" {
			if i+1 >= len(args) {
				return nil, "", "", errors.New("no output file after >")
			}
			outputFile = args[i+1]
			i += 2
			continue
		}

		resArgs = append(resArgs, args[i])
	}

	return resArgs, inputFile, outputFile, nil
}
