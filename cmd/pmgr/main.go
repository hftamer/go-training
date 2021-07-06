package main

import (
	"fmt"
	"os"

	"github.com/hftamer/go-training/internal/arguments"
	"github.com/hftamer/go-training/pkg/vault"
)

func main() {
	args := os.Args[1:]
	err := arguments.Check(args)
	if err != nil {
		fmt.Println(err)
		return
	}

	secretsVault, err := vault.New()
	if err != nil {
		fmt.Println(err)
		return
	}

	if args[0] == "add" {
		err = secretsVault.Add(args[1], args[2])
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(secretsVault.Get(args[1]))

		secretsVault.SaveData()
		if err != nil {
			fmt.Println(err)
		}
	}

	if args[0] == "get" {
		str, err := secretsVault.Get(args[1])
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(str)
	}
}
