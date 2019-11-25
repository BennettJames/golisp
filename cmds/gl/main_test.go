package main

import (
	"flag"
	"fmt"
	"testing"
)

func Test_clarifyFlags(t *testing.T) {
	var (
		flags   = flag.NewFlagSet("flags", flag.PanicOnError)
		outFile = flags.String("out", "", "")
	)
	flags.Parse([]string{
		"-out",
		"file.txt",
		"test.l",
	})
	fmt.Println("@@@ out", *outFile)
	fmt.Println("@@@ values", flags.Args())
}
