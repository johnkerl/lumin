# lumin

`lumin` is a simple command-line program which highlights matches to a
specified pattern (string or regex) in the specified files. This is like `grep`
with `--color`, except that `lumin` shows all lines, not just matching lines.

This uses ANSI 256-color escape sequences which work on Linux/Unix systems,
BSD-like systems, MacOS, etc., but typically not Windows.

To build the `lumin` executable:

- Install Go
- `go build`

Matching a string:

![screenshot1](./pix/screenshot1.png)

Matching a regular expression:

![screenshot2](./pix/screenshot2.png)
