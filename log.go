package main

import (
	"fmt"
	"os"
)

func fatal(err error) {
	fmt.Println(err.Error())
	os.Exit(1)
}

func (s surface) printErr(err error) {
	if s.tuiRunning.Load() {
		_, werr := s.buf.WriteString(err.Error())
		if werr != nil {
			fatal(werr)
		}
	} else {
		fatal(err)
	}
}
