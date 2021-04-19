package pmgr

import (
	"errors"
	"fmt"
)

type Command interface {
	Execute() error
}

type addCommand struct {
	name     string
	password string
}

func (cmd addCommand) Execute() error {
	vaultPath := GetVaultPath()

	vault, err := LoadVault(vaultPath)
	if err != nil {
		return err
	}

	err = vault.AddAccount(cmd.name, cmd.password)
	if err != nil {
		return err
	}

	return vault.Save(vaultPath)
}

type getCommand struct {
	name        string
}

func (cmd getCommand) Execute() error {
	vaultPath := GetVaultPath()

	vault, err := LoadVault(vaultPath)
	if err != nil {
		return err
	}

	pwd, err := vault.GetAccount(cmd.name)
	if err != nil {
		return err
	}

	fmt.Print(pwd)

	return nil
}

func NewCommand(args []string) (Command, error) {
	switch cmd := args[0]; cmd {
	case "add":
		if len(args) != 3 {
			return nil, errors.New("add command needs name and password")
		}

		return addCommand{
			name: args[1],
			password: args[2],
		}, nil
	case "get":
		if len(args) != 2 {
			return nil, errors.New("get command needs account name")
		}

		return getCommand{
			name: args[1],
		}, nil
	default:
		return nil, errors.New(cmd + " is not a recognized command")
	}
}





