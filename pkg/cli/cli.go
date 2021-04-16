package cli

import (
	"errors"
	"flag"
	"fmt"
	"github.com/hftamer/go-training/pkg/pmgr"
	"os"
)

type Args struct {
	IsHelp   bool
	Command  string
	Name     string
	Password string
}

func ParseArgs(cmdLine []string) (Args, error) {
	if len(cmdLine) < 1 {
		return Args{}, errors.New("no command provided")
	}

	args := Args{
		Command: cmdLine[0],
	}
	switch args.Command {
	case "add":
		if len(cmdLine) != 3 {
			return Args{}, errors.New("add command needs name and password")
		}
		args.Name = cmdLine[1]
		args.Password = cmdLine[2]
		break
	case "get":
		if len(cmdLine) != 2 {
			return Args{}, errors.New("get command needs account name")
		}
		args.Name = cmdLine[1]
		break
	default:
		return Args{}, errors.New(cmdLine[0] + " is not a recognized command")
	}

	isHelpPtr := flag.Bool("help", false, "display help")
	flag.Parse()

	args.IsHelp = *isHelpPtr

	return args, nil
}

func PrintHelp(exitCode int) {
	msg :=
`CLI password manager: pmgr <add|get|update|delete>

	add <accountName> <accountPassword>		adds an entry for the given name and password if <accountName> doesn't exist
	update <accountName> <newPassword>		updates an entry's password to the new password
	get <accountName>                 		gets the given account's password
	delete <accountName>              		deletes the given account name
`
	fmt.Print(msg)

	os.Exit(exitCode)
}

func Run(cmdLine []string) int {
	args, err := ParseArgs(cmdLine)
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		return 1
	}

	if args.IsHelp {
		PrintHelp(0)
	}

	if err := executeCommand(args); err != nil {
		fmt.Fprint(os.Stderr, err)
		return 1
	}

	return 0
}

func GetVaultPath() string {
	return "pmgr-vault.json"
}

func executeCommand(args Args) error {
	vaultPath := GetVaultPath()

	vault, err := pmgr.LoadVault(vaultPath)
	if err != nil {
		return err
	}

	// this declares a lambda function that will save the vault no matter what
	defer func() {
		if err := vault.Save(vaultPath); err != nil {
			fmt.Fprint(os.Stderr, err)
		}
	}()

	switch args.Command {
	case "add":
		return vault.AddAccount(args.Name, args.Password)
	case "get":
		if pwd, e := vault.GetAccount(args.Name); e == nil {
			fmt.Print(pwd)
			return nil
		} else {
			return e
		}
	default:
		panic("this shouldn't happen if ParseArgs() is implemented correctly")
	}
}

