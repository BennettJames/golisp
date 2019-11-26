package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/bennettjames/go-compiler-experiments/golisp2"
)

func main() {
	ctx, cancel := RootContext()
	defer cancel()
	var _ = ctx

	var (
		flags    = flag.NewFlagSet("flags", flag.PanicOnError)
		showVals = flags.Bool("show-vals", false,
			"Shows all evaluated values; rather than just printed ones")
	)
	flags.Parse(os.Args[1:])
	files := flags.Args()

	if len(files) != 1 {
		// note (bs): let's see if this can trigger an interpreter
		fmt.Fprint(os.Stderr, "gl requires a file argument to execute")
		return
	}

	if err := execFile(ctx, files[0], *showVals); err != nil {
		log.Fatal(err)
	}
}

func execFile(ctx context.Context, file string, showVals bool) error {
	f, err := os.Open(file)
	if err != nil {
		return fmt.Errorf("Could not read file '%s': %w", file, err)
	}

	// note (bs): consider folding these up into a utility method. It seems
	// reasonable to have a "prepare file" function.
	ts := golisp2.NewTokenScanner(
		golisp2.NewRuneScanner(file, f),
	)
	exprs, exprsErr := golisp2.ParseTokens(ts)
	if exprsErr != nil {
		return fmt.Errorf("Parse error in '%s': %w", file, exprsErr)
	}
	baseCtx := golisp2.BuiltinContext()
	execCtx := baseCtx.SubContext(nil)

	for _, e := range exprs {
		if val, err := e.Eval(execCtx); err != nil {
			return fmt.Errorf("Execution error in '%s': %w", file, err)
		} else if _, isNil := val.(*golisp2.NilValue); !isNil && showVals {
			fmt.Println(val.InspectStr())
		}
	}

	return nil
}
