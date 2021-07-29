# filist

[![GitHub license](https://img.shields.io/github/license/onozaty/filist)](https://github.com/onozaty/filist/blob/main/LICENSE)
[![Test](https://github.com/onozaty/filist/actions/workflows/test.yaml/badge.svg)](https://github.com/onozaty/filist/actions/workflows/test.yaml)

Displays a list of files.

## Usage

```
$ filist -s -M .
a.txt   24  3d3a42d900823afcfdfeb6de338bcec1
b/1.txt 81  ae23e0b40e773ac132f477f661e89b86
b/2.txt 163 494ba81d0d828ff9a244da627b5ece47
```

The arguments are as follows.

```
Usage: filist [options] directory ...
options
  -r, --rel      Print relative path (If neither 'rel' nor 'abs' is specified, 'rel' will be printed first column.)
  -a, --abs      Print absolute path
  -s, --size     Print file size
  -m, --mtime    Print modification time
  -M, --md5      Print MD5 hash
  -S, --sha1     Print SHA-1 hash
      --sha256   Print SHA-256 hash
  -h, --help     Help
```

Prints in the order the options are specified.

```
$ filist -M -r .
3d3a42d900823afcfdfeb6de338bcec1  a.txt
ae23e0b40e773ac132f477f661e89b86  b/1.txt
494ba81d0d828ff9a244da627b5ece47  b/2.txt
```

## Install

You can download the binary from the following.

* https://github.com/onozaty/filist/releases/latest

## License

MIT

## Author

[onozaty](https://github.com/onozaty)
