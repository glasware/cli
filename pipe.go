package main

import (
	"io"
)

// TODO: replace with a buffer that deletes the front of it after
// it reaches max size and use in all write operations.
type pipe struct {
	surface *surface
}

var _ io.Writer = new(pipe)

func (p pipe) Write(b []byte) (int, error) {
	l, err := p.surface.buf.Write(b)
	if err != nil {
		return l, err
	}

	p.surface.prog.Send(forceUpdateMsg{})
	return l, nil
}
