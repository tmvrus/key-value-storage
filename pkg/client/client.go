package client

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log/slog"
)

const defaultReadBufferSize = 1024

type Client struct {
	socket readerWriter
	log    *slog.Logger
}

func (c *Client) execute(cmd []byte) ([]byte, error) {
	_, err := c.socket.Write(cmd)
	if err != nil {
		return nil, fmt.Errorf("write command: %w", err)
	}

	result := make([]byte, defaultReadBufferSize)
	n, err := c.socket.Read(result)
	if err != nil {
		return nil, fmt.Errorf("read result: %w", err)
	}

	return result[:n], nil
}

func (c *Client) StartInteractionLoop(in io.Reader, out io.Writer) {
	const eolByte byte = '\n'

	input := bufio.NewReader(in)

	_, err := out.Write([]byte("Waiting for command\n"))
	if err != nil {
		c.log.Error("write to out", "error", err.Error())
		return
	}

	for {
		cmd, err := input.ReadBytes(eolByte)
		if err != nil {
			if errors.Is(err, io.EOF) {
				c.log.Debug("got EOL from cmd input")
				return
			}
			c.log.Error("read cmd new comment", "error", err.Error())
			continue
		}

		result, err := c.execute(cmd)
		if err != nil {
			c.log.Error("execute command", "error", err.Error(), "command", string(cmd))
			continue
		}

		_, err = out.Write(result)
		if err != nil {
			c.log.Error("write result", "error", err.Error())
			continue
		}
		_, err = out.Write([]byte{eolByte})
		if err != nil {
			c.log.Error("write eol", "error", err.Error())
			continue
		}
	}
}

func NewClient(i readerWriter, log *slog.Logger) *Client {
	return &Client{socket: i, log: log}
}
