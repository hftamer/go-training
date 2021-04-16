package main

import (
	"github.com/hftamer/go-training/pkg/cli"
	"os"
)

func main() {
	cli.Run(os.Args[1:])
}
