# Invoice generator

Generate a basic invoice via CLI. 

Logic: If not last week of a month, generates invoice for prev month.

See example_generated.pdf for example output.

## Usage

Fill out example_config.yaml and name it as config.yaml

Example run:

```
go mod tidy
go build

./invoice -h
Usage of ./invoice:
  -days int
        Number of days worked
  -outdir string
        Output directory for invoice files (defaults to current directory) (default ".")
  -po string
        Purchase Order reference number

./invoice -days 4 -po 2007
```