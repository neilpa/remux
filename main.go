// Remux is a command line regex multiplexer
package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"regexp"

	xflag "neilpa.me/go-x/flag"
)

var (
	version = "v1.0.0-dev"

	flags *flag.FlagSet
)

func main() {
	os.Exit(realMain(os.Args[1:], os.Stdout))
}

func realMain(args []string, stdout io.Writer) int {
	flags = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	var inFlags xflag.MultiString
	flags.Var(&inFlags, "i", "input file(s), can be set multiple times")
	verFlag := flags.Bool("v", false, "print version and exit")
	flags.Parse(args)

	if *verFlag {
		fmt.Fprintln(stdout, version)
		return 0
	}
	if flags.NArg() == 0 {
		return usageError("no filter specified")
	}
	if len(inFlags) == 0 {
		inFlags = append(inFlags, "-")
	}

	stdin := false
	readers := make([]io.Reader, len(inFlags))
	for i, path := range inFlags {
		if path == "" || path == "-" {
			stdin = true
			readers[i] = os.Stdin
		} else {
			f, err := os.Open(path)
			if err != nil {
				return fatal(err.Error())
			}
			defer f.Close()
			readers[i] = f
		}
	}
	if stdin && len(inFlags) > 1 {
		warn("STDIN specified mulitple times and/or with other files")
	}

	sinks := make([]sink, 0)
	for i := 0; i < flags.NArg(); i += 2 {
		arg := flags.Arg(i)
		re, err := regexp.Compile(arg)
		if err != nil {
			return fatal("invalid regex %q: %s", arg, err)
		}

		s := sink{re, nil}
		path := ""
		if i+1 < flags.NArg() {
			path = flags.Arg(i + 1)
		}
		switch path {
		case "", "-":
			s.w = stdout
		default:
			f, err := os.Create(path) // TODO: Allow for appending?
			if err != nil {
				return fatal("%s", err)
			}
			defer f.Close()
			s.w = f
		}
		sinks = append(sinks, s)
	}

	scanner := bufio.NewScanner(io.MultiReader(readers...))
	for scanner.Scan() {
		line := scanner.Text()
		for _, s := range sinks {
			if s.re.MatchString(line) {
				_, err := io.WriteString(s.w, line+"\n")
				if err != nil {
					return fatal("write: ", err)
				}
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return fatal("scanner: %s", err)
	}
	return 0
}

type sink struct {
	re *regexp.Regexp
	w  io.Writer
}

func usageError(msg string) int {
	fmt.Fprintln(os.Stderr, msg)
	printUsage()
	return 2
}

func warn(format string, args ...interface{}) {
	format = os.Args[0] + ": warn: " + format + "\n"
	fmt.Fprintf(os.Stderr, format, args...)
}

func fatal(format string, args ...interface{}) int {
	format = os.Args[0] + ": " + format + "\n"
	fmt.Fprintf(os.Stderr, format, args...)
	return 1
}

func printUsage() {
	fmt.Fprintf(os.Stderr, `Usage: %s [options] <regex0> [<file0> [<regex1> <file1> ...]]

  Remux is a regex multiplexer filtering input to multiple target files.

Options:

`, os.Args[0])
	flags.PrintDefaults()
	fmt.Fprintln(os.Stderr)
}
