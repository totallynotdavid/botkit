// Command envfile writes .env.example or the generated text to standard output
// when requested, from the settings the bots read.
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/totallynotdavid/botkit/internal/envfile"
)

const (
	exampleMode = 0o644
	usageExit   = 2
)

func main() {
	stdout := flag.Bool("stdout", false, "write the generated example to standard output")

	flag.Parse()

	if flag.NArg() > 0 {
		fmt.Fprintln(os.Stderr, "usage: envfile [--stdout]")
		os.Exit(usageExit)
	}

	text, err := envfile.Example(envfile.Bots())
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if *stdout {
		_, err = os.Stdout.WriteString(text)
	} else {
		err = os.WriteFile(".env.example", []byte(text), exampleMode)
	}

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
