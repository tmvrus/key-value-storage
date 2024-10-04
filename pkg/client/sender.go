package client

import (
	"context"
	"fmt"
)

const readBufferSize = 1024

type netSender struct {
	conn     readerWriter
	readBuff []byte
}

func newSender(conn readerWriter) netSender {
	return netSender{
		conn:     conn,
		readBuff: make([]byte, readBufferSize),
	}
}
func (s netSender) Send(ctx context.Context, cmd string) (string, error) {
	ch := s.send(cmd)

	select {
	case <-ctx.Done():
		return "", ctx.Err()
	case r := <-ch:
		return r.result, r.err
	}
}

type sendResult struct {
	result string
	err    error
}

func (s netSender) send(cmd string) <-chan sendResult {
	r := make(chan sendResult, 1)
	go func() {
		defer close(r)

		_, err := s.conn.Write([]byte(cmd))
		if err != nil {
			r <- sendResult{err: fmt.Errorf("write %w: ", err)}
			return
		}

		n, err := s.conn.Read(s.readBuff)
		if err != nil {
			r <- sendResult{err: fmt.Errorf("read %w: ", err)}
			return
		}

		result := string(s.readBuff[:n])
		s.readBuff = s.readBuff[:]

		r <- sendResult{result: result}
	}()
	return r
}
