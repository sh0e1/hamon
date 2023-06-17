package main

import (
	"context"
	"os"

	"github.com/sh0e1/hamon/internal/cmd"
)

func main() {
	cmd.Run(context.Background(), cmd.New("hamon"), os.Args[1:])
}
