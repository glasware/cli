package main

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"os"

	pb "github.com/glasware/glas-core/proto"
	"github.com/gogo/protobuf/jsonpb"
	"github.com/gorilla/websocket"
	"github.com/ingcr3at1on/x/sigctx"
)

func startWeb(
	ctx context.Context,
	cancel context.CancelFunc,
	url string,
) error {
	return sigctx.StartWithContext(ctx, func(ctx context.Context) error {
		errCh := make(chan error, 1)

		conn, _, err := websocket.DefaultDialer.Dial(url, nil)
		if err != nil {
			return err
		}
		defer conn.Close()

		go func() {
			for {
				_, byt, err := conn.ReadMessage()
				if err != nil {
					errCh <- err
					return
				}

				if string(byt) == `closing connection` {
					fmt.Println(string(byt))
					cancel()
					return
				} else {
					var output pb.Output
					if err := jsonpb.Unmarshal(bytes.NewReader(byt), &output); err != nil {
						errCh <- err
						return
					}

					switch output.Type {
					case pb.Output_BUFFERED:
						fmt.Print(output.Data)
					case pb.Output_UNBUFFERED:
						fmt.Print(output.Data)
					case pb.Output_INSTRUCTION:
						// TODO: handle instructions.
					}
				}
			}
		}()

		go func() {
			scanner := bufio.NewScanner(os.Stdin)
			for scanner.Scan() {
				if err = conn.WriteJSON(&pb.Input{
					Data: scanner.Text(),
				}); err != nil {
					errCh <- err
					return
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

		fmt.Println("exiting")
		return nil
	})
}
