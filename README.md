# go-build-cleaner

This removes following that go generated. Often, `go run` doesn't clean temporary directory which go generated.

* Remove temporary directory.
* Remove Windows Fireall Rules.

## Usage

```
Usage of go-build-cleaner.exe:
  -dryrun
    	dryrun
  -verbose
    	verbose
```

If you run this on Windows, you must run this as administrator permitted.

## Installation

```
$ go get github.com/mattn/go-build-cleaner
```

## License

MIT

## Author

Yasuhiro Matsumoto (a.k.a. mattn)
