package main

import (
	"context"
	"fmt"
	"net/url"
	"os"

	"github.com/ingcr3at1on/x/sigctx"
	"github.com/spf13/cobra"
)

var (
	root = &cobra.Command{
		Use:           "glas",
		SilenceErrors: true,
		SilenceUsage:  true,
		Hidden:        true,
		RunE:          runE,
	}
)

const (
	glasWeb = `glasweb`
)

func init() {
	flags := root.Flags()
	flags.StringP(glasWeb, "w", "", "A websocket server address for a remote Glas Gateway")
}

func runE(cmd *cobra.Command, args []string) error {
	ctx, cancel := sigctx.WithCancel(context.Background())

	web, err := cmd.Flags().GetString(glasWeb)
	if err != nil {
		return err
	}
	if web != `` {
		_url, err := url.Parse(web)
		if err != nil {
			return err
		}

		return startWeb(ctx, cancel, _url.String())
	}

	return startStandalone(ctx, cancel)
}

func main() {
	if err := root.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
