package client

import (
	"bufio"
	"fmt"
	"io"
)

const welcomeString = "@> "

type consoleInteractor struct {
	io      readerWriter
	scanner *bufio.Scanner
}

func newConsoleInteractor(io readerWriter) consoleInteractor {
	return consoleInteractor{
		io:      io,
		scanner: bufio.NewScanner(io),
	}
}

func (i consoleInteractor) ReadCommand() (string, error) {
	err := i.WriteResult(welcomeString)
	if err != nil {
		return "", fmt.Errorf("write welcome string: %w", err)
	}

	if i.scanner.Scan() {
		return i.scanner.Text(), nil
	}

	err = i.scanner.Err()
	if err == nil {
		err = io.EOF
	}
	return "", fmt.Errorf("command scanner: %w", err)
}

func (i consoleInteractor) WriteResult(s string) error {
	_, err := i.io.Write([]byte(s))
	if err != nil {
		return fmt.Errorf("write: %w", err)
	}

	return nil
}
