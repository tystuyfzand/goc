goc
===

A simple wrapper to run your Go builds with, supporting multiple threads for multiple builds at once.

Usage
-----

Instead of `go build`, you use `goc`. The `GOOS` and `GOARCH` env variables can be set to multiple, comma separated items.

Example source build command:

`go build -o goc`

Resulting command:
`GOOS=windows,linux,drawin GOARCH=386,amd64,arm64 goc -o goc`

This will generate binaries for Windows, Linux, and MacOS, with whichever architectures are supported on each OS.

If you wish to change how goc outputs binaries (for example, output to folders for each OS) you can use the `GOBINARY` environment variable:

`GOBINARY="{os}/{name}_{os}_{arch}"`

Credits
-------

[across](https://github.com/LordRusk/across) was the inspiration for this, however I wanted to do it a bit differently.

Caveats
-------

Builds requiring native code/gcc are not supported, there are other programs for this.