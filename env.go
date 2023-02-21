package main

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/spf13/afero"
	"github.com/spf13/pflag"
)

var (
	configPath *string
)

func loadEnv() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	defaultConfigPath := filepath.Join(home, ".glas")

	configPath = pflag.String("config", defaultConfigPath, "define the configuration path")

	pflag.Parse()

	if *configPath == defaultConfigPath {
		if path, ok := os.LookupEnv(`GLAS_CONFIG`); ok {
			configPath = &path
		}
	}

	abs, err := filepath.Abs(*configPath)
	if err != nil {
		return err
	}

	afs = afero.NewOsFs()

	if err := afs.Mkdir(abs, 0765); err != nil && !errors.Is(err, os.ErrExist) {
		return err
	}

	return nil
}
