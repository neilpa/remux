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
	version = "undefined"

	serialFlag  = flag.Bool("s", false, "process input serially")
	versionFlag = flag.Bool("v", false, "print version and exit")
	inFlags     xflag.MultiString
)

func init() {
	flag.Var(&inFlags, "i", "input file(s), can be set multiple times")
}

func main() {
	os.Exit(realMain(os.Args[1:]))
}

func realMain(args []string) int {
	flag.CommandLine.Parse(args)
	if *versionFlag {
		fmt.Println(version)
		return 0
	}

	if flag.NArg() == 0 {
		return usageError("no filter specified")
	}
	if *serialFlag {
		return usageError("todo: serial processing not implemented yet")
	}
	if len(inFlags) > 0 {
		return usageError("todo: input files not implemented yet")
	}

	sinks := make([]sink, 0)
	for i := 0; i < flag.NArg(); i += 2 {
		arg := flag.Arg(i)
		re, err := regexp.Compile(arg)
		if err != nil {
			return fatal("invalid regex %q: %s", arg, err)
		}

		s := sink{re, nil}
		path := ""
		if i+1 < flag.NArg() {
			path = flag.Arg(i + 1)
		}
		switch path {
		case "", "-":
			s.w = os.Stdout
		default:
			f, err := os.Open(path) // TODO: Allow for appending?
			if err != nil {
				return fatal("%s", err)
			}
			defer f.Close()
			s.w = f
		}
		sinks = append(sinks, s)
	}

	scanner := bufio.NewScanner(os.Stdin)
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
	flag.PrintDefaults()
	fmt.Fprintln(os.Stderr)
}
