package main

import (
	"fmt"

	"github.com/hftamer/go-training/internal/arguments"
)

func main() {
	err := arguments.Check()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("valid arguments passed")
}
