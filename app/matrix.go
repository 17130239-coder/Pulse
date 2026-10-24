package app

import (
	"math"

	"github.com/egoist/mygo/ui"
)

// MatrixOptions configures the appearance of the Dot Matrix LED visualizer.
type MatrixOptions struct {
	Rows           int
	Cols           int
	PrimaryColor   ui.Color
	SecondaryColor ui.Color // Optional threshold or dual-series color
	ThresholdRow   int      // Rows >= ThresholdRow use SecondaryColor
	HasSecondary   bool
}

// DotMatrix renders a hardware dot-matrix LED activity graph.
func DotMatrix(c *ui.Context, values []float64, opt MatrixOptions) ui.Element {
	if opt.Rows <= 0 {
		opt.Rows = 7
	}
	if opt.Cols <= 0 {
		opt.Cols = 32
	}

	el := ui.Row(c).Height(46).FillWidth()
	el.Draw(func(p *ui.Painter, r ui.Rect) {
		cols := opt.Cols
		rows := opt.Rows

		// Calculate spacing and dot size
		colWidth := r.W / float32(cols)
		rowHeight := r.H / float32(rows)

		dotDiameter := float32(math.Min(float64(colWidth)*0.55, float64(rowHeight)*0.55))
		if dotDiameter < 2.5 {
			dotDiameter = 2.5
		}
		if dotDiameter > 5.5 {
			dotDiameter = 5.5
		}
		radius := dotDiameter / 2.0

		// Extract last 'cols' data points from values slice
		data := make([]float64, cols)
		valLen := len(values)
		for i := 0; i < cols; i++ {
			srcIdx := valLen - cols + i
			if srcIdx >= 0 && srcIdx < valLen {
				data[i] = values[srcIdx]
			} else {
				data[i] = 0.0
			}
		}

		// Draw the LED grid
		for col := 0; col < cols; col++ {
			fraction := data[col]
			if fraction < 0 {
				fraction = 0
			}
			if fraction > 1.0 {
				fraction = 1.0
			}

			litCount := int(math.Round(fraction * float64(rows)))
			if fraction > 0.02 && litCount == 0 {
				litCount = 1
			}

			cx := r.X + (float32(col)+0.5)*colWidth

			for row := 0; row < rows; row++ {
				// row 0 is bottom, row rows-1 is top
				yIdx := rows - 1 - row
				cy := r.Y + (float32(yIdx)+0.5)*rowHeight

				dotRect := ui.Rect{
					X: cx - radius,
					Y: cy - radius,
					W: dotDiameter,
					H: dotDiameter,
				}

				if row < litCount {
					// Lit dot
					dotColor := opt.PrimaryColor
					if opt.HasSecondary && row >= opt.ThresholdRow {
						dotColor = opt.SecondaryColor
					}
					p.Fill(dotRect, dotColor, radius)
				} else {
					// Unlit dot
					p.Fill(dotRect, AppTheme.DotUnlit, radius)
				}
			}
		}
	})
	return el
}
