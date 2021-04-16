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

	switch cmdLine[0] {
	case "add":
		if len(cmdLine) != 3 {
			return Args{}, errors.New("add command needs name and password")
		}
		break
	default:
		return Args{}, errors.New(cmdLine[0] + " is not a recognized command")
	}

	isHelpPtr := flag.Bool("help", false, "display help")

	flag.Parse()

	return Args {
		IsHelp:   *isHelpPtr,
		Command:  cmdLine[0],
		Name: 	  cmdLine[1],
		Password: cmdLine[2],
	}, nil
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

	switch args.Command {
	case "add":
		err = vault.AddAccount(args.Name, args.Password)
		break
	default:
		panic("this shouldn't happen if ParseArgs() is implemented correctly")
	}

	if err != nil {
		return err
	}

	return vault.Save(vaultPath)
}

