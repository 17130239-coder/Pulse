package app

import (
	"fmt"
	"math"

	"github.com/egoist/mygo/ui"
)

// DonutSegment represents a slice in the donut chart.
type DonutSegment struct {
	Label string
	Value float64
	Color ui.Color
}

// DonutChart renders a circular ring breakdown chart.
func DonutChart(c *ui.Context, centerPercent int, segments []DonutSegment) ui.Element {
	var total float64
	for _, s := range segments {
		total += s.Value
	}
	if total <= 0 {
		total = 1.0
	}

	el := ui.Row(c).Size(110, 110).Center()
	el.Draw(func(p *ui.Painter, r ui.Rect) {
		cx := r.X + r.W/2.0
		cy := r.Y + r.H/2.0
		radius := (math.Min(float64(r.W), float64(r.H)) - 16.0) / 2.0
		strokeWidth := float32(11.0)

		// First draw background ring
		bgPath := (&ui.Path{}).Circle(cx, cy, float32(radius))
		p.StrokePath(bgPath, strokeWidth, AppTheme.CardBorder)

		// Draw each segment arc
		var currentFraction float32 = 0.0
		for _, s := range segments {
			fraction := float32(s.Value / total)
			if fraction <= 0.001 {
				continue
			}
			start := currentFraction
			end := currentFraction + fraction
			if end > 1.0 {
				end = 1.0
			}

			arcPath := (&ui.Path{}).Circle(cx, cy, float32(radius))
			arcPath.Trim(start, end)
			p.StrokePath(arcPath, strokeWidth, s.Color)

			currentFraction = end
		}
	})

	// Center text
	el.Children(func() {
		ui.Column(c).Center().Children(func() {
			ui.Text(c, fmt.Sprintf("%d%%", centerPercent)).FontSize(17).Bold().TextColor(AppTheme.TextPrimary)
			ui.Text(c, "in use").FontSize(10).TextColor(AppTheme.TextMuted)
		})
	})

	return el
}
