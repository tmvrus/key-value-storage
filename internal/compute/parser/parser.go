package parser

import (
	"fmt"
	"strings"

	"github.com/tmvrus/key-value-storage/internal/domain"
)

func Parse(s string) (cmd domain.Command, err error) {
	parts := strings.Split(s, " ")
	t := domain.CommandType(parts[0])

	switch len(parts) {
	case 2:
		if t != domain.CommandGet && t != domain.CommandDelete {
			err = fmt.Errorf("invalid command %s for size %d", t, 2)
			return
		}
		if parts[1] == "" {
			err = fmt.Errorf("empty key")
			return
		}
		cmd.Type = t
		cmd.Key = parts[1]
		return
	case 3:
		if t != domain.CommandSet {
			err = fmt.Errorf("invalid command %s for size %d", t, 3)
			return
		}
		if parts[1] == "" || parts[2] == "" {
			err = fmt.Errorf("empty key or value")
			return
		}
		cmd.Type = t
		cmd.Key = parts[1]
		cmd.Value = parts[2]
		return
	default:
		err = fmt.Errorf("invalid parts count")
		return
	}
}
