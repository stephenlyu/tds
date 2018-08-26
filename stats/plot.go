package stats

import (
	"github.com/gonum/plot/plotter"
	"github.com/gonum/plot"
	"github.com/gonum/plot/vg"
	"fmt"
	"github.com/stephenlyu/tds/util"
	"image/color"
)

var plotColors = []color.Color {
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

func Plot(title string, titles []string, plotData map[string][]float64, pdfFile string) error {
	util.Assert(len(titles) == len(plotData), "")
	convert := func (values []float64) plotter.XYs {
		pts := make(plotter.XYs, len(values))
		for i := range pts {
			pts[i].X = float64(i)
			pts[i].Y = values[i]

		}
		return pts
	}

	p, err := plot.New()
	if err != nil {
		return err
	}

	p.Title.Text = title
	p.X.Label.Text = ""
	p.Y.Label.Text = "Profit"
	p.Y.Tick.Marker = round2Ticks{}

	for i, title := range titles {
		values := plotData[title]
		points := convert(values)
		line, err := plotter.NewLine(points)
		line.LineStyle.Color = plotColors[i % len(plotColors)]
		line.LineStyle.Width = 1

		if err != nil {
			return err
		}
		p.Add(line)
		p.Legend.Add(title, line)

	}

	return p.Save(297*vg.Millimeter, 210*vg.Millimeter, pdfFile)
}

type round2Ticks struct{}
// Ticks computes the default tick marks, but inserts commas
// into the labels for the major tick marks.
func (round2Ticks) Ticks(min, max float64) []plot.Tick {
	tks := plot.DefaultTicks{}.Ticks(min, max)
	for i, t := range tks {
		if t.Label == "" { // Skip minor ticks, they are fine.
			continue
		}
		tks[i].Label = fmt.Sprintf("%.02f", t.Value)
	}
	return tks
}
