package graph

import (
	"fmt"
	"io"
	"time"

	"github.com/wcharczuk/go-chart"
)

type LineChart struct {
	filterColumns []string
	width         int
	height        int
}

func NewLineChart(columns []string, width, height int) *LineChart {
	return &LineChart{
		filterColumns: columns,
		width:         width,
		height:        height,
	}
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

func (c *LineChart) Read(r io.Reader) (chart.Chart, error) {
	columns, rows, err := Parse(r)
	if err != nil {
		return chart.Chart{}, err
	}
	if len(rows) < 10 {
		return chart.Chart{}, fmt.Errorf("more than 10 rows of data needed to render chart")
	}

	targetColumns := c.filterColumns
	if len(c.filterColumns) < 1 {
		targetColumns = columns
	}
	columnExists := make(map[string]struct{})
	for _, k := range columns {
		columnExists[k] = struct{}{}
	}

	filterCols := make([]string, 0, len(columnExists))
	for _, k := range targetColumns {
		if _, ok := columnExists[k]; ok {
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
				Show:        true,
				StrokeColor: color,
				FillColor:   color.WithAlpha(15),
			},
			XValues: xvalues,
			YValues: yvalues,
		})
	}

	return chart.Chart{
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
	}, nil
}
