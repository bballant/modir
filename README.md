# Go Utilities

This repository contains a collection of Go utilities, including the Soccer Sub Schedule Application and other useful tools. Each utility is located in its own directory under the `cmd` folder.

## Utilities

### Soccer Sub Schedule Application (Cheetah)

The Soccer Sub Schedule Application, also known as `cheetah`, is a tool for generating images of soccer field formations with player positions and substitution information. This utility reads CSV files containing player positions and creates a visual representation of the changes in the lineup over time.

- Source: `./cmd/cheetah/cheetah.go`
- Example input CSV file: `./cmd/cheetah/wildcats.csv`
- Sample output image: `./cmd/cheetah/soccer_fields_sample.png`

For detailed instructions on how to use the Soccer Sub Schedule Application, please refer to the [Cheetah README](./cmd/cheetah/README.md).

### Golden Ragwort

`goldenragwort` is a utility that helps you create simple images of flowers without using any external libraries.

- Source: `./cmd/goldenragwort/goldenrawort.go`
- Sample output image: `./cmd/goldenragwort/simple_flower_no_lib.png`

### Lion

`lion` is another utility in this repository.

- Source: `./cmd/lion/lion.go`

## Installation

To install the utilities, navigate to the specific utility directory under `cmd` and run `go install`. For example, to install the Soccer Sub Schedule Application, run the following commands:

```bash
cd cmd/cheetah
go install
```

This will create an executable in your `$GOPATH/bin` directory, which you can run from anywhere if `$GOPATH/bin` is in your system's `PATH`.

## License

This repository is licensed under the [MIT License](./LICENSE).
