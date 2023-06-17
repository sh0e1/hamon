package cmd

import (
	"context"
	"flag"
	"fmt"
	"os"
)

func Run(ctx context.Context, cmd Commander, args []string) {
	s := flag.NewFlagSet(cmd.Name(), flag.ExitOnError)
	if err := cmd.Run(ctx, args...); err != nil {
		fmt.Fprintf(s.Output(), "%s: %v\n", cmd.Name(), err)
		os.Exit(2)
	}
}

type Commander interface {
	Name() string
	Run(ctx context.Context, args ...string) error
}

func New(nama string) *Main {
	return &Main{
		name: nama,
	}
}

type Main struct {
	name string

	Commander
}

func (m *Main) Name() string {
	return m.name
}

func (m *Main) Run(ctx context.Context, args ...string) error {
	fmt.Println(m.name)
	return nil
}
