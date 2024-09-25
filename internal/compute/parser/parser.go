package parser

import (
	"fmt"
	"strings"

	"github.com/tmvrus/key-value-storage/internal/domain"
)

type parserFunc func(string) (domain.Command, error)

func getParser(s string) parserFunc {
	m := map[domain.CommandType]parserFunc{
		domain.CommandGet:    parseGet,
		domain.CommandSet:    parseSet,
		domain.CommandDelete: parseDelete,
	}
	for k, v := range m {
		if strings.HasPrefix(s, string(k)) {
			return v
		}
	}

	return func(s string) (cmd domain.Command, err error) {
		err = fmt.Errorf("unsupporterd operation")
		return
	}
}

func parseGet(s string) (cmd domain.Command, err error) {
	parts := strings.Split(s, " ")
	if len(parts) != 2 {
		err = fmt.Errorf("invalid arguments number for GET command")
		return
	}
	if parts[1] == "" {
		err = fmt.Errorf("empty arguments for GET command")
		return
	}

	cmd.Type = domain.CommandGet
	cmd.Key = parts[1]
	return
}

func parseSet(s string) (cmd domain.Command, err error) {
	parts := strings.Split(s, " ")
	if len(parts) != 3 {
		err = fmt.Errorf("invalid arguments number for SET command")
		return
	}
	if parts[1] == "" || parts[2] == "" {
		err = fmt.Errorf("empty arguments for SET command")
		return
	}

	cmd.Type = domain.CommandSet
	cmd.Key = parts[1]
	cmd.Value = parts[2]
	return
}

func parseDelete(s string) (cmd domain.Command, err error) {
	parts := strings.Split(s, " ")
	if len(parts) != 2 {
		err = fmt.Errorf("invalid arguments number for DELETE command")
		return
	}
	if parts[1] == "" {
		err = fmt.Errorf("empty arguments for DELETE command")
		return
	}

	cmd.Type = domain.CommandDelete
	cmd.Key = parts[1]
	return
}

func Parse(s string) (cmd domain.Command, err error) {
	return getParser(s)(s)
}
