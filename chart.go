package graph

import (
	"io"
	"os"
	"time"

	"github.com/wcharczuk/go-chart"
)

type Chart interface {
	Read(io.Reader) (chart.Chart, error)
}

type DstatRecord map[string]float64
type DstatCSVRow struct {
	Time   time.Time
	Values DstatRecord
}

func RenderToFile(graph chart.Chart, f *os.File) error {
	graph.Elements = []chart.Renderable{chart.LegendThin(&graph)}
	return graph.Render(chart.PNG, f)
}
