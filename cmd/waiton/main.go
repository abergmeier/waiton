package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/abergmeier/waiton/internal/execute"
)

var (
	printHelp   = flag.Bool("help", false, "Print usage")
	expectedOut = flag.String("output", "", "Compare to command output")
	verbose     = flag.Bool("verbose", false, "Print verbose information")
	interval    = flag.Int("interval", 1, "Seconds until executing again")
)

func init() {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), `Usage of %s: -output <out> [flags] <command>
Available flags:
`, os.Args[0])
		flag.PrintDefaults()
	}
}

func main() {

	// Do not use flag for parsing since we need more narrow
	// handling of arguments
	args := parsePrefixFlags()

	for {
		collectedOutput := execute.MustExecuteInSh(args, *verbose)

		if string(collectedOutput) == *expectedOut {
			os.Exit(0)
		}

		time.Sleep(time.Duration(*interval) * time.Second)
	}
}

func parsePrefixFlags() []string {
	args := os.Args
	if len(os.Args) > 1 {
		args = os.Args[1:]
	}

	known := true

	for known && len(args) != 0 {
		switch args[0] {
		case "-h":
			fallthrough
		case "--help":
			*printHelp = true
			flag.Usage()
			os.Exit(0)
		case "-output":
			fallthrough
		case "--output":
			fallthrough
		case "-out":
			*expectedOut = args[1]
			args = args[2:]
		case "-verbose":
			fallthrough
		case "--verbose":
			fallthrough
		case "-v":
			*verbose = true
			args = args[1:]
		case "-interval":
			fallthrough
		case "--interval":
			i, err := strconv.Atoi(args[1])
			if err != nil {
				panic(err)
			}
			*interval = i
		default:
			known = false
		}
	}

	if *expectedOut == "" {
		flag.Usage()
		fmt.Fprintln(os.Stderr, "output flag missing")
		os.Exit(1)
	}

	if len(args) == 0 {
		flag.Usage()
		fmt.Fprintln(os.Stderr, "command missing")
		os.Exit(1)
	}

	return args
}
