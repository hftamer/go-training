package pmgr

import (
	"testing"
)

func TestGetQuote(t *testing.T) {
	quote := GetQuote()
	if quote == "" {
		t.Fatal("could not get quote")
	}
}
