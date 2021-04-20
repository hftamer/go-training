package cli

import (
	"testing"
)

func TestGivingUnrecognizedCommandShouldFail(t *testing.T) {
	_, err := ParseArgs([]string{"not-a-command"})
	if err == nil {
		t.Error("passing in an invalid argument should fail")
	}
}

func TestValidCommandsShouldNotFail(t *testing.T) {
	validCommands := [][]string {
		{"add", "foo", "bar"},
		{"get", "foo"},
		{"update", "foo", "bar"},
		{"delete", "foo"},
	}

	for _, cmd := range validCommands {
		if _, err := ParseArgs(cmd); err != nil {
			t.Error("valid commands should not fail with: ", err)
		}
	}
}

func TestHelpFlag(t *testing.T) {
	args, err := ParseArgs([]string{"-help"})
	if err != nil {
		t.Error("help flag parsing failed: ", err)
	}

	if !args.IsHelp {
		t.Error("args.IsHelp should be true")
	}

	if args.Command != nil {
		t.Error("args.Command should be nil if args.IsHelp is true")
	}
}
