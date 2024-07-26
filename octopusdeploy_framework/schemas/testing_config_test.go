package schemas

import (
	"flag"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	// Parse only known flags
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	_ = flag.CommandLine.Parse(os.Args[1:]) // ignore error
	os.Exit(m.Run())
}
