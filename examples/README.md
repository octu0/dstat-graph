# example

To run example, you need to install dstat-graph in `$GOPATH`

```
go get github.com/octu0/dstat-graph
```

Next, generate a graph using csv in the example directory

```
$ cd examples
$ go run $GOPATH/src/github.com/octu0/dstat-graph/cmd/main.go --csv ./dstat.csv -o cpu.png -f usr,sys,idl,wai
```

You can check the generated image using the viewer app.

![cpu](https://github.com/octu0/dstat-graph/blob/master/examples/cpu.png?raw=true)  

useful as [dstat2graph.sh](https://github.com/octu0/dstat-graph/blob/master/examples/dstat2graph.sh)
