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
	Command  pmgr.Command
}

func ParseArgs(cmdLine []string) (Args, error) {
	if len(cmdLine) < 1 {
		return Args{}, errors.New("no command provided")
	}

	isHelpPtr := flag.Bool("help", false, "display help")
	flag.Parse()

	// we want to parse out any flags first
	// if *isHelpPtr is true, that means the -help flag was passed
	// and we should not attempt to parse a command out of the command line
	if *isHelpPtr {
		return Args {
			IsHelp: true,
			Command: nil,
		}, nil
	}

	cmd, err := pmgr.NewCommand(cmdLine)
	if err != nil {
		return Args{}, err
	}

	return Args {
		IsHelp: *isHelpPtr,
		Command: cmd,
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

	if err := args.Command.Execute(); err != nil {
		fmt.Fprint(os.Stderr, err)
		return 1
	}

	return 0
}