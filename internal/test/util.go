package test

import (
	"os"
	"testing"
)

func SkipCI(t *testing.T, reason string) {
	if os.Getenv("Skip_Legacy_Tests") == "" {
		t.Skip(reason)
	}
}
