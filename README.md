# grpcdoc

![build workflow](https://github.com/ScaleableWebWorks/grpcdoc/actions/workflows/build.yml/badge.svg)

> Create comprehensive html documentation from your grpc/protobuf files.

`grpcdoc` is a command line tool written in golang which takes grpc/protobuf files and generates a comprehensive html documentation for it.

## Installation

OS X & Linux:

```sh
# TODO
```

Windows:

```sh
# TODO
```

## Usage example

__Create a documentation for a single file__

```sh
grpcdoc -out doc.html ./path/to/your/file.proto
```

__Create a documentation for multiple files__

```sh
grpcdoc -out doc.html ./path/to/your/file1.proto ./path/to/your/file2.proto
```

__Read all .proto files from a directory__

```sh
grpcdoc -out doc.html ./path/to/your/protos
```

__Create a documentation for multiple files and include a custom css file__

```sh
grpcdoc -out doc.html -style custom.css ./path/to/your/file1.proto
```

_For more examples and usage, please refer to the [Wiki][wiki]._

## Development setup

To start development you need to have [golang](https://go.dev/dl/) installed on your system.

```sh
git clone https://github.com/ScaleableWebWorks/grpcdoc.git
cd grpcdoc

go get
go build
```

## Release History

* 0.0.1
    * Initial version Work in progress

## Meta

Marco Rico – [@mricog](https://twitter.com/mricog) – marco@scaleablewebworks.com

Distributed under the MIT license. See ``LICENSE`` for more information.

[https://github.com/ScaleableWebWorks](https://github.com/ScaleableWebWorks)

## Contributing

1. Fork it (<https://github.com/ScaleableWebWorks/grpcdoc/fork>)
2. Create your feature branch (`git checkout -b feature/fooBar`)
3. Commit your changes (`git commit -am 'Add some fooBar'`)
4. Push to the branch (`git push origin feature/fooBar`)
5. Create a new Pull Request

<!-- Markdown link & img dfn's -->
[wiki]: https://github.com/ScaleableWebWorks/grpcdoc/wiki
