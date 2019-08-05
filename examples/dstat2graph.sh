#!/bin/bash
go run $GOPATH/src/github.com/octu0/dstat-graph/cmd/main.go --csv ./dstat.csv -o cpu.png -f usr,sys,idl,wai
go run $GOPATH/src/github.com/octu0/dstat-graph/cmd/main.go --csv ./dstat.csv -o traf.png -f recv,send
go run $GOPATH/src/github.com/octu0/dstat-graph/cmd/main.go --csv ./dstat.csv -o intr.png -f int,csw
go run $GOPATH/src/github.com/octu0/dstat-graph/cmd/main.go --csv ./dstat.csv -o net.png -f lis,act,syn,tim,clo
