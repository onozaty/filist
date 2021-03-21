# filist

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
  -a, --abs      Absolute path
  -s, --size     Print file size
  -m, --mtime    Print modification time
  -M, --md5      Print MD5 hash
  -S, --sha1     Print SHA-1 hash
      --sha256   Print SHA-256 hash
  -h, --help     Help
```

## Install

You can download the binary from the following.

* https://github.com/onozaty/filist/releases/latest

## License

MIT

## Author

[onozaty](https://github.com/onozaty)
