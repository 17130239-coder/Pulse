package app

import (
	"fmt"
	"os/exec"
	"strconv"

	"github.com/egoist/mygo/ui"
)

// ViewProcesses renders the activity/task manager view.
func (s *AppState) ViewProcesses(c *ui.Context) {
	procs := s.GetFilteredProcesses()

	ui.Column(c).Grow(1).FillWidth().Padding(20, 28, 28).Gap(16).Children(func() {
		// Header row with search
		ui.Row(c).FillWidth().AlignItems(ui.Center).Children(func() {
			ui.Column(c).Grow(1).Gap(2).Children(func() {
				ui.Text(c, "Processes").FontSize(22).Bold().TextColor(AppTheme.TextPrimary)
				ui.Text(c, fmt.Sprintf("%d active processes", len(procs))).FontSize(13).TextColor(AppTheme.TextMuted)
			})

			// Search bar
			ui.Row(c).Width(240).Padding(5, 10).Radius(8).Background(AppTheme.CardBg).Border(1, AppTheme.CardBorder).AlignItems(ui.Center).Gap(8).Children(func() {
				ui.Icon(c, IconSearch).Size(14, 14).TextColor(AppTheme.TextDim)
				ui.TextInput(c, &s.SearchQuery).Placeholder("Search processes...").Grow(1).FontSize(12)
			})
		})

		// Table Container
		ui.Column(c).Grow(1).FillWidth().Radius(10).Background(AppTheme.CardBg).Border(1, AppTheme.CardBorder).Children(func() {
			// Table Header
			ui.Row(c).FillWidth().Padding(10, 14).Background(ui.RGBA(255, 255, 255, 0.02)).BorderWidth(0, 0, 1, 0).BorderColor(AppTheme.CardBorder).Children(func() {
				ui.Text(c, "Process Name").FontSize(11).Bold().TextColor(AppTheme.TextMuted).Width(200)
				ui.Text(c, "PID").FontSize(11).Bold().TextColor(AppTheme.TextMuted).Width(70)
				ui.Text(c, "% CPU").FontSize(11).Bold().TextColor(AppTheme.TextMuted).Width(70)
				ui.Text(c, "Memory").FontSize(11).Bold().TextColor(AppTheme.TextMuted).Width(90)
				ui.Text(c, "Command Path").FontSize(11).Bold().TextColor(AppTheme.TextMuted).Grow(1)
				ui.Text(c, "Action").FontSize(11).Bold().TextColor(AppTheme.TextMuted).Width(70).TextAlign(ui.End)
			})

			// Table Rows
			ui.Scroll(c).Grow(1).FillWidth().Children(func() {
				ui.Column(c).FillWidth().Children(func() {
					for _, p := range procs {
						proc := p
						ui.Row(c).FillWidth().Padding(8, 14).AlignItems(ui.Center).BorderWidth(0, 0, 1, 0).BorderColor(ui.RGBA(255, 255, 255, 0.03)).Children(func() {
							// Name with icon
							ui.Row(c).Width(200).AlignItems(ui.Center).Gap(8).Children(func() {
								ui.Icon(c, IconProcesses).Size(13, 13).TextColor(AppTheme.Blue)
								ui.Text(c, proc.Name).FontSize(12).Bold().TextColor(AppTheme.TextPrimary).SingleLine()
							})

							ui.Text(c, fmt.Sprintf("%d", proc.PID)).FontSize(12).TextColor(AppTheme.TextDim).Width(70)
							ui.Text(c, fmt.Sprintf("%.1f%%", proc.CPU)).FontSize(12).Bold().TextColor(AppTheme.TextPrimary).Width(70)
							ui.Text(c, fmt.Sprintf("%d MB", proc.MemMB)).FontSize(12).TextColor(AppTheme.TextMuted).Width(90)
							ui.Text(c, proc.Command).FontSize(11).TextColor(AppTheme.TextDim).Grow(1).SingleLine()

							// End process button
							killBtn := ui.Row(c).Padding(2, 6).Radius(4).Background(ui.RGBA(239, 68, 68, 0.12)).AlignItems(ui.Center)
							if killBtn.Clicked() {
								exec.Command("kill", "-9", strconv.Itoa(proc.PID)).Run()
								s.UpdateSnapshot()
							}
							killBtn.Children(func() {
								ui.Text(c, "End").FontSize(11).Bold().TextColor(ui.Hex("#ef4444"))
							})
						})
					}
				})
			})
		})
	})
}
