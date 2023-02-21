package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"sync/atomic"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/glasware/glas-core"
	"github.com/spf13/afero"
)

var afs afero.Fs

type (
	surface struct {
		m    tea.Model
		buf  *bytes.Buffer
		in   chan string
		glas glas.Glas

		prog       *tea.Program
		tuiRunning atomic.Bool

		errCh chan error
	}
)

func init() {
	if err := loadEnv(); err != nil {
		fatal(err)
	}
}

func main() {
	if err := start(); err != nil {
		fatal(err)
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
	s.glas, err = glas.New(s.in, pipe{&s}, *configPath, glas.OptAfs(afs))
	if err != nil {
		return err
	}

	s.prog = tea.NewProgram(s.m)

	go func() {
		s.tuiRunning.Store(true)
		s.m, err = s.prog.Run()
		s.tuiRunning.Store(false)
		if err != nil {
			s.errCh <- fmt.Errorf("p.Run -- %w", err)
		}

		close(s.errCh)
	}()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		if err := s.glas.Start(ctx); err != nil {
			s.errCh <- err
		}
	}()

	if err := <-s.errCh; err != nil {
		if errors.Is(err, glas.ErrExit) {
			cancel()
			s.prog.Send(tea.Quit)
		} else {
			// Attempt to print to the TUI if the TUI is not running errors are fatal.
			s.printErr(err)
		}
	}

	return nil
}
