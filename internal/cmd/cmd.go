package cmd

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"path/filepath"
)

func Run(ctx context.Context, cmd Commander, args []string) {
	if err := cmd.Run(ctx, args...); err != nil {
		fmt.Fprintf(os.Stderr, "%s: %v\n", cmd.Name(), err)
		os.Exit(2)
	}
}

type Commander interface {
	Name() string
	Run(ctx context.Context, args ...string) error
}

func New(name string) *Main {
	return &Main{
		name: name,
		network: "unix",
		address: filepath.Join(os.TempDir(), fmt.Sprintf("%s.sock", name)),
	}
}

type Main struct {
	name string
	network string
	address string
}

func (m *Main) Name() string {
	return m.name
}

func (m *Main) Run(ctx context.Context, args ...string) error {
	sigctx, stop := signal.NotifyContext(ctx, os.Interrupt)
	defer stop()

	listener, err := net.Listen(m.network, m.address)
	if err != nil {
		return err
	}
	defer func() {
		listener.Close()
		os.Remove(m.address)
	}()
	log.Printf("Server is running on %s\n", listener.Addr())

	go func (ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				conn, err := listener.Accept()
				if err != nil {
					if !errors.Is(err, net.ErrClosed) {
						log.Printf("Accept failed: %v\n", err)
					}
					break
				}
				defer conn.Close()
				log.Println("Accepted connection")

				go func (ctx context.Context, conn net.Conn) {
					const buffer = 1024
					for {
						select {
						case <-ctx.Done():
							return
						default:
							buf := make([]byte, buffer)
							if _, err := conn.Read(buf); err != nil {
								log.Printf("Read failed: %v\n", err)
								break
							}

							if _, err := conn.Write(buf); err != nil {
								log.Printf("Write failed: %v\n", err)
							}
						}
					}
				}(ctx, conn)
			}
		}
	}(sigctx)

	<-sigctx.Done()
	log.Println("Server is terminated")

	return nil
}
