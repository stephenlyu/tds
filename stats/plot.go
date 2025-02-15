package stats

import (
	"fmt"
	"image/color"
	"math"

	"github.com/stephenlyu/tds/util"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
)

var plotColors = []color.Color{
	color.RGBA{255, 0, 0, 255},
	color.RGBA{0, 255, 0, 255},
	color.RGBA{0, 0, 255, 255},
	color.RGBA{255, 0, 255, 255},
	color.RGBA{0, 255, 255, 255},
	color.RGBA{192, 0, 0, 255},
	color.RGBA{0, 192, 0, 255},
	color.RGBA{0, 0, 192, 255},
	color.RGBA{0, 192, 192, 255},
	color.RGBA{192, 192, 0, 255},
	color.RGBA{192, 0, 192, 255},
	color.RGBA{128, 128, 0, 255},
	color.RGBA{0, 0, 0, 255},
}

func Plot(title string, titles []string, xTicks []string, plotData map[string][]float64, pdfFile string, useLogScale bool) error {
	util.Assert(len(titles) == len(plotData), "")
	convert := func(values []float64) plotter.XYs {
		pts := make(plotter.XYs, len(values))
		for i := range pts {
			pts[i].X = float64(i)
			if useLogScale {
				pts[i].Y = math.Log(values[i])
			} else {
				pts[i].Y = values[i]
			}
		}
		return pts
	}

	p := plot.New()

	p.Title.Text = title
	p.X.Label.Text = ""
	p.Y.Label.Text = ""
	p.Y.Tick.Marker = round2Ticks{useLogScale: useLogScale}

	if xTicks != nil {
		p.X.Tick.Marker = stringTicks{ticks: xTicks}
		p.X.Tick.Label.Rotation = 0.7
		p.X.Tick.Label.YAlign = draw.YCenter
		p.X.Tick.Label.XAlign = draw.XRight
	}

	for i, title := range titles {
		values := plotData[title]
		if xTicks != nil {
			util.Assert(len(values) == len(xTicks), "")
		}
		points := convert(values)
		line, err := plotter.NewLine(points)
		line.LineStyle.Color = plotColors[i%len(plotColors)]
		line.LineStyle.Width = 1

		if err != nil {
			return err
		}
		p.Add(line)
		p.Legend.Add(title, line)

	}

	return p.Save(297*vg.Millimeter, 210*vg.Millimeter, pdfFile)
}

type round2Ticks struct {
	useLogScale bool
}

// Ticks computes the default tick marks, but inserts commas
// into the labels for the major tick marks.
func (this round2Ticks) Ticks(min, max float64) []plot.Tick {
	tks := plot.DefaultTicks{}.Ticks(min, max)
	for i, t := range tks {
		if t.Label == "" { // Skip minor ticks, they are fine.
			continue
		}
		var value float64
		if this.useLogScale {
			value = math.Exp(t.Value)
		} else {
			value = t.Value
		}
		tks[i].Label = fmt.Sprintf("%.02f", value)
	}
	return tks
}

type stringTicks struct {
	ticks []string
}

// Ticks computes the default tick marks, but inserts commas
// into the labels for the major tick marks.
func (this stringTicks) Ticks(min, max float64) []plot.Tick {
	tks := plot.DefaultTicks{}.Ticks(min, max)
	for i, t := range tks {
		if t.Label == "" { // Skip minor ticks, they are fine.
			continue
		}
		v := int(util.Round(t.Value, 0))
		tks[i].Label = this.ticks[v]
	}
	return tks
}
