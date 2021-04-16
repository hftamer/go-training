package cli

import (
	"flag"
	"fmt"
	"os"
)

type Args struct {
	IsHelp bool
}

func ParseArgs() (Args, error) {
	isHelpPtr := flag.Bool("help", false, "display help")

	flag.Parse()

	return Args {
		IsHelp: *isHelpPtr,
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

func Run() int {
	args, err := ParseArgs()
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		return 1
	}

	if args.IsHelp {
		PrintHelp(0)
	}

	return 0
}

