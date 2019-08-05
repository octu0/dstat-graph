package main

import (
  "log"
  "os"
  "strings"
  "io/ioutil"

  "github.com/comail/colog"
  "gopkg.in/urfave/cli.v1"

  "github.com/octu0/dstat-graph"
)

var (
  Commands = make([]cli.Command, 0)
)
func AddCommand(cmd cli.Command){
  Commands = append(Commands, cmd)
}

func action(c *cli.Context) error {
  filename := c.String("csv")
  if filename == "" {
    log.Printf("error: requires dstat csv path(--help print usage)")
    return nil
  }
  f, err := os.Open(filename)
  if err != nil {
    log.Printf("error: failed to open file: %s", filename)
    return err
  }
  defer f.Close()

  outfile := c.String("out")
  if outfile == "" {
    tmpfile, err := ioutil.TempFile("", "chart-*.png")
    if err != nil {
      log.Printf("error: failed to create tempfile: %s", err.Error())
      return err
    }
    defer os.Remove(tmpfile.Name())
    outfile = tmpfile.Name()
  }
  out, err := os.OpenFile(outfile, os.O_RDWR|os.O_CREATE, 0655)
  if err != nil {
    log.Printf("error: failed to open tempfile: %s", err.Error())
    return err
  }
  defer out.Close()

  filter   := c.String("column")
  columns  := make([]string, 0)
  if "" != filter {
    for _, v := range strings.Split(filter, ",") {
      v = strings.TrimSpace(v)
      if "" != v {
        columns = append(columns, v)
      }
    }
  }

  chartType := c.String("chart")
  if chartType == "" {
    chartType = "line"
  }

  var chart graph.Chart
  switch chartType {
  case "line":
    chart = graph.NewLineChart(columns, c.Int("width"), c.Int("height"))
  default:
    log.Printf("error: unknown chartype: %s", chartType)
    return nil
  }

  if err := chart.Read(f); err != nil {
    log.Printf("error: csv read error: %s", err.Error())
    return err
  }
  if err := chart.RenderToFile(out); err != nil {
    log.Printf("error: render to file error: %s", err.Error())
    return err
  }
  log.Printf("info: write to file: %s", outfile)
  return nil
}

func main(){
  colog.SetDefaultLevel(colog.LDebug)
  colog.SetMinLevel(colog.LInfo)

  colog.SetFormatter(&colog.StdFormatter{
    Flag: log.Ldate | log.Ltime | log.Lshortfile,
  })
  colog.Register()

  app         := cli.NewApp()
  app.Version  = graph.Version
  app.Name     = graph.AppName
  app.Author   = ""
  app.Email    = ""
  app.Usage    = ""
  app.Action   = action
  app.Commands = Commands
  app.Flags    = []cli.Flag{
    cli.StringFlag{
      Name: "csv, i",
      Usage: "/path/to/csv dstat csv path",
      Value: "",
    },
    cli.StringFlag{
      Name: "out, o",
      Usage: "output file path (if argument is empty, write to tmpfile for parse test)",
      Value: "",
    },
    cli.StringFlag{
      Name: "column, f",
      Usage: "pickup columns (defaults: plot all columns)",
      Value: "",
    },
    cli.StringFlag{
      Name: "chart, t",
      Usage: "chart-type 'line' or 'bar' (defaults: 'line')",
      Value: "line",
    },
    cli.IntFlag{
      Name: "width",
      Usage: "image width",
      Value: 600,
    },
    cli.IntFlag{
      Name: "height",
      Usage: "image height",
      Value: 400,
    },
    cli.BoolFlag{
      Name: "debug, d",
      Usage: "debug mode",
    },
    cli.BoolFlag{
      Name: "verbose, V",
      Usage: "verbose. more message",
    },
  }
  if err := app.Run(os.Args); err != nil {
    log.Printf("error: %s", err.Error())
    cli.OsExiter(1)
  }
}
