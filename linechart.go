package graph

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/wcharczuk/go-chart"
)

type LineChart struct {
	filterColumns []string
	graph         chart.Chart
	width         int
	height        int
}

func NewLineChart(columns []string, width, height int) *LineChart {
	c := new(LineChart)
	c.filterColumns = columns
	c.width = width
	c.height = height
	return c
}

func (c *LineChart) interval(values []time.Time, interval int) []chart.GridLine {
	t := make([]chart.GridLine, 0)
	freq := len(values) / interval
	for i := 0; i < len(values); i += freq {
		t = append(t, chart.GridLine{
			Value: float64(values[i].UnixNano()),
		})
	}
	return t
}
func (c *LineChart) Read(r io.Reader) error {
	columns, rows, err := Parse(r)
	if err != nil {
		return err
	}
	targetColumns := c.filterColumns
	if len(c.filterColumns) < 1 {
		targetColumns = columns
	}
	columnExists := make(map[string]bool)
	for _, k := range columns {
		columnExists[k] = true
	}

	filterCols := make([]string, 0)
	for _, k := range targetColumns {
		if columnExists[k] {
			filterCols = append(filterCols, k)
		}
	}

	xvalues := make([]time.Time, 0, len(rows))
	for _, row := range rows {
		xvalues = append(xvalues, row.Time)
	}
	series := make([]chart.Series, 0, len(filterCols))
	for i, col := range filterCols {
		yvalues := make([]float64, 0, len(rows))
		for _, row := range rows {
			if value, ok := row.Values[col]; ok == true {
				yvalues = append(yvalues, value)
			}
		}
		idx := i % len(chart.DefaultAlternateColors)
		color := chart.DefaultAlternateColors[idx]
		series = append(series, chart.TimeSeries{
			Name: col,
			Style: chart.Style{
				StrokeColor: color,
				FillColor:   color.WithAlpha(15),
			},
			XValues: xvalues,
			YValues: yvalues,
		})
	}
	c.graph = chart.Chart{
		Width:  c.width,
		Height: c.height,
		Background: chart.Style{
			Padding: chart.Box{
				Top: 50,
			},
		},
		YAxis: chart.YAxis{
			ValueFormatter: func(v interface{}) string {
				return fmt.Sprintf("%3.2f", v.(float64))
			},
		},
		XAxis: chart.XAxis{
			ValueFormatter: func(v interface{}) string {
				format := "01-02 03:04:05"
				if t, isTime := v.(time.Time); isTime {
					return t.Format(format)
				}
				if t, isFloat := v.(float64); isFloat {
					return time.Unix(0, int64(t)).Format(format)
				}
				return fmt.Sprintf("<unknown_axis>: %#v", v)
			},
			GridMajorStyle: chart.Style{
				StrokeColor: chart.ColorAlternateGray,
				StrokeWidth: 2.0,
			},
			GridLines: c.interval(xvalues, 10),
		},
		Series: series,
	}
	return nil
}
func (c *LineChart) RenderToFile(f *os.File) error {
	c.graph.Elements = []chart.Renderable{chart.LegendThin(&c.graph)}
	return c.graph.Render(chart.PNG, f)
}
