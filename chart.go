package graph

import(
  "io"
  "os"
  "time"
)

type Chart interface {
  Read(io.Reader) error
  RenderToFile(*os.File) error
}

type DstatRecord map[string]float64
type DstatCSVRow struct {
  Time      time.Time
  Values    DstatRecord
}
