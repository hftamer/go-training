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

	return vault.Save()
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

	pwd, err := vault.GetAccountPassword(cmd.name)
	if err != nil {
		return err
	}

	fmt.Print(pwd)

	return nil
}

type updateCommand struct {
	name        string
	newPassword string
}

func (cmd updateCommand) Execute() error {
	vaultPath := GetVaultPath()

	vault, err := LoadVault(vaultPath)
	if err != nil {
		return err
	}

	err = vault.UpdateAccount(cmd.name, cmd.newPassword)
	if err != nil {
		return err
	}

	return vault.Save()
}

type deleteCommand struct {
	name string
}

func (cmd deleteCommand) Execute() error {
    vaultPath := GetVaultPath()

    vault, err := LoadVault(vaultPath)
    if err != nil {
    	return err
	}

	err = vault.DeleteAccount(cmd.name)
	if err != nil {
		return err
	}

	return vault.Save()
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
	case "update":
		if len(args) != 3 {
			return nil, errors.New("update command needs name and new password")
		}

		return updateCommand{
			name: args[1],
			newPassword: args[2],
		}, nil
	case "delete":
		if len(args) != 2 {
			return nil, errors.New("delete command needs account name")
		}

		return deleteCommand{
			name: args[1],
		}, nil
	default:
		return nil, errors.New(cmd + " is not a recognized command")
	}
}





