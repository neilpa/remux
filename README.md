# remux

A command line regex multiplexer for filtering and streaming data to multiple targets.

## Usage

```sh
remux [options] <regex0> <file0> [<regex1> <file1> ...]
```

Omitting the last file implies STDOUT as the target, as do an explicit empty string or `-` file.

## Alternatives

* [Process substitution](https://unix.stackexchange.com/a/43536)
* [Explicit file descriptors](https://unix.stackexchange.com/a/43536)
* [Named pipes](https://unix.stackexchange.com/a/43536)
* [`pee`](https://linux.die.net/man/1/pee) command in [`moreutils`](https://packages.debian.org/en/sid/moreutils)
* [`MULTIOS`](http://zsh.sourceforge.net/Doc/Release/Redirection.html#Multios) option from `zsh`

## License

MIT
