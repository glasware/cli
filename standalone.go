package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"sync"

	glas "github.com/glasware/glas-core"
	pb "github.com/glasware/glas-core/proto"
	"github.com/ingcr3at1on/x/sigctx"
)

func startStandalone(
	ctx context.Context,
	cancel context.CancelFunc,
) error {
	return sigctx.StartWithContext(ctx, func(ctx context.Context) error {
		var wg sync.WaitGroup
		errCh := make(chan error, 1)
		inCh := make(chan *pb.Input)
		outCh := make(chan *pb.Output)

		// Don't put this in our waitgroup, it will never finish.
		go func() {
			for {
				out := <-outCh
				if out != nil {
					n, err := os.Stdout.WriteString(out.Data)
					if err != nil {
						errCh <- err
						return
					}

					if n != len(out.Data) {
						errCh <- io.ErrShortWrite
						return
					}
				}
			}
		}()

		g, err := glas.New(&glas.Config{
			Input:  inCh,
			Output: outCh,
		})
		if err != nil {
			return err
		}

		wg.Add(1)
		go func() {
			defer wg.Done()

			if err := g.Start(ctx, cancel); err != nil {
				errCh <- err
				return
			}
		}()

		// Don't put this in the waitgroup because it can and will continue running
		// until we stop it.
		go func() {
			// FIXME: I don't think we can actually use a scanner here,
			// better to detect the enter/return key somehow; the issue
			// is we need to be able to send an empty string in some
			// cases.
			scanner := bufio.NewScanner(os.Stdin)
			for scanner.Scan() {
				inCh <- &pb.Input{
					Data: scanner.Text(),
				}
			}

			if err := scanner.Err(); err != nil {
				if err != io.EOF {
					errCh <- err
				}
			}
		}()

		select {
		case <-ctx.Done():
			break
		case err := <-errCh:
			if err != nil {
				return err
			}
		}

		wg.Wait()
		fmt.Println("exiting")
		return nil
	})
}
