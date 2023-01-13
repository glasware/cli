package main

import (
	"bytes"
	"context"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/glasware/glas-core"
	"github.com/spf13/afero"
)

type (
	surface struct {
		m    model
		buf  *bytes.Buffer
		in   chan string
		glas glas.Glas
		prog *tea.Program

		errCh chan error
	}
)

func main() {
	if err := start(); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func start() error {
	s := surface{
		buf:   new(bytes.Buffer),
		in:    make(chan string),
		errCh: make(chan error, 1),
	}

	s.m = initialModel(&s)

	var err error
	s.glas, err = glas.New(s.in, pipe{&s}, "./cfg", glas.OptAfs(afero.NewOsFs())) // FIXME: this path should come from env or flags.
	if err != nil {
		return err
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		if err := s.glas.Start(ctx); err != nil {
			s.errCh <- fmt.Errorf("m.glas.Start -- %w", err)
		}
	}()

	s.prog = tea.NewProgram(s.m)

	go func() {
		if err := s.prog.Start(); err != nil {
			s.errCh <- fmt.Errorf("p.Start -- %w", err)
		}

		close(s.errCh)
	}()

	if err := <-s.errCh; err != nil {
		return err
	}

	return nil
}
