# dstat-graph

generate png chart from dstat CSV file.  
data visualization using [go-chart](https://github.com/wcharczuk/go-chart)  
(inspired by [dstat2graph](https://github.com/sh2/dstat2graphs))

## Usage

output dstat csv using `dstat --output`

```
$ dstat -t --cpu --mem --disk --io --net --int --sys --tcp --output ./dstat.csv
```

load csv into `dstat-graph` (column filter with `-f` if necessary)

```
$ dstat-graph --csv ./dstat.csv -o cpu.png -f usr,sys,idl,wai
```

see more [examples](https://github.com/octu0/dstat-graph/tree/master/examples).

## Build

Build requires Go version 1.11+ installed.

```
$ go version
```

Run `make pkg` to Build and package for linux, darwin.

```
$ git clone https://github.com/octu0/dstat-graph
$ make pkg
```

## Help

```
NAME:
   dstat-graph

USAGE:
   dstat-graph [global options] command [command options] [arguments...]

VERSION:
   1.0.0

COMMANDS:
     help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --csv value, -i value     /path/to/csv dstat csv path
   --out value, -o value     output file path (if argument is empty, write to tmpfile for parse test)
   --column value, -f value  pickup columns (defaults: plot all columns)
   --chart value, -t value   chart-type 'line' or 'bar' (defaults: 'line') (default: "line")
   --width value             image width (default: 600)
   --height value            image height (default: 400)
   --debug, -d               debug mode
   --verbose, -V             verbose. more message
   --help, -h                show help
   --version, -v             print the version
```
