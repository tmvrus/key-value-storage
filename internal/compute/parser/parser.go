package parser

import (
	"fmt"
	"strings"

	"github.com/tmvrus/key-value-storage/internal/domain"
)

type parseArgFunc func([]string) (domain.Command, error)

func parseGet(args []string) (cmd domain.Command, err error) {
	if len(args) != 1 {
		err = fmt.Errorf("invalid arguments number for GET command")
		return
	}
	if args[0] == "" {
		err = fmt.Errorf("empty arguments for GET command")
		return
	}

	cmd.Type = domain.CommandGet
	cmd.Key = args[0]
	return
}

func parseSet(args []string) (cmd domain.Command, err error) {
	if len(args) != 2 {
		err = fmt.Errorf("invalid arguments number for SET command")
		return
	}
	if args[0] == "" || args[1] == "" {
		err = fmt.Errorf("empty arguments for SET command")
		return
	}

	cmd.Type = domain.CommandSet
	cmd.Key = args[0]
	cmd.Value = args[1]
	return
}

func parseDelete(args []string) (cmd domain.Command, err error) {
	if len(args) != 1 {
		err = fmt.Errorf("invalid arguments number for DELETE command")
		return
	}
	if args[0] == "" {
		err = fmt.Errorf("empty arguments for DELETE command")
		return
	}

	cmd.Type = domain.CommandDelete
	cmd.Key = args[0]
	return
}

func Parse(s string) (cmd domain.Command, err error) {
	m := map[domain.CommandType]parseArgFunc{
		domain.CommandGet:    parseGet,
		domain.CommandSet:    parseSet,
		domain.CommandDelete: parseDelete,
	}

	args := strings.Split(s, " ")
	if len(args) < 2 {
		err = fmt.Errorf("invalid arguments numbers")
		return
	}

	f, ok := m[domain.CommandType(args[0])]
	if !ok {
		err = fmt.Errorf("unsupporterd operation")
		return
	}

	return f(args[1:])
}
