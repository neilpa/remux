# remux

[![CI](https://github.com/neilpa/remux/workflows/CI/badge.svg)](https://github.com/neilpa/remux/actions/)

A command line regex multiplexer for filtering and streaming lines to multiple targets.

```sh
tail -f log.txt | remux ERROR errors.txt WARN warnings.txt INFO verbose.txt
```

## Install

```sh
go get -u neilpa.me/remux
```

## Usage

```
remux [options] <regex0> [<file0> [<regex1> <file1> ...]]
```

By default, source data is streamed from STDIN, use `-i` to specify a file instead. The `-i` argument can be set more than once to read from multiple input sources in parallel. Use `-s` to instead read serially in the order specified.

An empty string or `-` can be provided to `-i` to explicitly use STDIN. Similarly, an empty string or `-` for an output file to use STDOUT. Additionally, omitting the last file also implies STDOUT as the target.

## Alternatives

* [Process substitution](https://unix.stackexchange.com/a/43536)
* [Explicit file descriptors](https://unix.stackexchange.com/a/43536)
* [Named pipes](https://unix.stackexchange.com/a/43536)
* [`pee`](https://linux.die.net/man/1/pee) command in [`moreutils`](https://packages.debian.org/en/sid/moreutils)
* [`MULTIOS`](http://zsh.sourceforge.net/Doc/Release/Redirection.html#Multios) option from `zsh`

## License

MIT
