package arguments

import (
	"errors"
	"os"
)

func Check() error {
	args := os.Args[1:]

	if len(args) == 0 {
		return errors.New("no arguments passed")
	}

	if args[0] == "add" || args[0] == "update" {
		if len(args) == 3 {
			return nil
		}

		return errors.New("invalid number of arguments passed")
	}

	if args[0] == "get" || args[0] == "delete" {
		if len(args) == 2 {
			return nil
		}

		return errors.New("invalid number of arguments passed")
	}

	return errors.New("incorrect command passed, should be add, delete, update or get")
}
