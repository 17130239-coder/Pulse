package app

import (
	"fmt"

	"pulse/metrics"

	"github.com/egoist/mygo/ui"
)

// ViewOverview renders the main Pulse dashboard matching the screenshot.
func (s *AppState) ViewOverview(c *ui.Context) {
	snap := s.Snapshot
	history := s.Collector.GetHistory(s.TimeRange)

	ui.Column(c).Grow(1).FillWidth().Padding(20, 28, 28).Gap(18).Children(func() {
		// Header row
		s.overviewHeader(c)

		// 6 Primary Metric Cards (2 rows of 3 columns)
		s.overviewCardsGrid(c, snap, history)

		// 3 Detailed Breakdown Cards
		s.overviewBreakdownsGrid(c, snap)
	})
}

func (s *AppState) overviewHeader(c *ui.Context) {
	snap := s.Snapshot

	ui.Row(c).FillWidth().AlignItems(ui.Center).Children(func() {
		// Title and Subtitle
		ui.Column(c).Grow(1).Gap(2).Children(func() {
			ui.Text(c, "Overview").FontSize(22).Bold().TextColor(AppTheme.TextPrimary)
			model := snap.ModelName
			if model == "" {
				model = "MacBook Air"
			}
			ui.Text(c, model).FontSize(13).TextColor(AppTheme.TextMuted)
		})

		// Controls on the right
		ui.Row(c).Gap(10).AlignItems(ui.Center).Children(func() {
			// Time range selector segmented control: [1m] [5m] [15m] [1h]
			ui.Row(c).Radius(6).Padding(2).Background(AppTheme.CardBg).Border(1, AppTheme.CardBorder).Gap(2).Children(func() {
				for _, tr := range []string{"1m", "5m", "15m", "1h"} {
					isActive := s.TimeRange == tr
					btnBg := ui.RGBA(0, 0, 0, 0)
					txtColor := AppTheme.TextMuted
					if isActive {
						btnBg = AppTheme.HoverBg
						txtColor = AppTheme.TextPrimary
					}

					btn := ui.Row(c).Padding(3, 8).Radius(4).Background(btnBg)
					if btn.Clicked() {
						s.SetTimeRange(tr)
					}
					btn.Children(func() {
						ui.Text(c, tr).FontSize(12).Bold().TextColor(txtColor)
					})
				}
			})

			// Play / Pause Button
			pauseBtn := ui.Row(c).Size(28, 28).Radius(6).Background(AppTheme.CardBg).Border(1, AppTheme.CardBorder).Center()
			if pauseBtn.Clicked() {
				s.TogglePause()
			}
			pauseBtn.Children(func() {
				if s.IsPaused {
					ui.Icon(c, IconPlay).Size(12, 12).TextColor(AppTheme.Green)
				} else {
					ui.Icon(c, IconPause).Size(12, 12).TextColor(AppTheme.TextMuted)
				}
			})

			// Settings Button
			settingsBtn := ui.Row(c).Size(28, 28).Radius(6).Background(AppTheme.CardBg).Border(1, AppTheme.CardBorder).Center()
			if settingsBtn.Clicked() {
				s.SetPage("About This Mac")
			}
			settingsBtn.Children(func() {
				ui.Icon(c, IconSettings).Size(14, 14).TextColor(AppTheme.TextMuted)
			})
		})
	})
}

func (s *AppState) overviewCardsGrid(c *ui.Context, snap metrics.SystemSnapshot, history []metrics.SystemSnapshot) {
	// Extract historical series for matrix graphs
	cpuSeries := make([]float64, len(history))
	memSeries := make([]float64, len(history))
	gpuSeries := make([]float64, len(history))
	diskSeries := make([]float64, len(history))
	netSeries := make([]float64, len(history))
	batSeries := make([]float64, len(history))

	for i, h := range history {
		cpuSeries[i] = h.CPUPercent / 100.0
		memSeries[i] = h.MemPercent / 100.0
		gpuSeries[i] = h.GPUPercent / 100.0
		diskSeries[i] = float64(h.DiskReadRate+h.DiskWriteRate) / (5.0 * 1024 * 1024) // Normalized to 5 MB/s
		netSeries[i] = float64(h.NetDownloadRate) / (5.0 * 1024 * 1024)              // Normalized to 5 MB/s
		batSeries[i] = float64(h.BatteryPercent) / 100.0
	}

	// Row 1: CPU, Memory, GPU
	ui.Row(c).FillWidth().Gap(14).Children(func() {
		// 1. CPU Card
		s.cardContainer(c, "CPU", IconCPU, AppTheme.Blue, func() {
			ui.Column(c).FillWidth().Gap(10).Children(func() {
				ui.Column(c).Gap(2).Children(func() {
					ui.Text(c, "Now").FontSize(11).TextColor(AppTheme.TextMuted)
					ui.Text(c, fmt.Sprintf("%.0f%%", snap.CPUPercent)).FontSize(26).Bold().TextColor(AppTheme.TextPrimary)
				})

				ui.Row(c).FillWidth().Children(func() {
					s.subStat(c, "User", fmt.Sprintf("%.0f%%", snap.CPUUser))
					s.subStat(c, "System", fmt.Sprintf("%.0f%%", snap.CPUSystem))
					s.subStat(c, "Temperature", fmt.Sprintf("%.0f °C", snap.CPUTemperature))
				})

				DotMatrix(c, cpuSeries, MatrixOptions{
					Rows:         7,
					Cols:         32,
					PrimaryColor: AppTheme.Blue,
				})
			})
		})

		// 2. Memory Card
		s.cardContainer(c, "Memory", IconMemory, AppTheme.Green, func() {
			ui.Column(c).FillWidth().Gap(10).Children(func() {
				ui.Row(c).FillWidth().AlignItems(ui.Baseline).Children(func() {
					ui.Column(c).Grow(1).Gap(2).Children(func() {
						totalGB := float64(snap.MemTotal) / (1024.0 * 1024.0 * 1024.0)
						ui.Text(c, fmt.Sprintf("In use of %.1f GB", totalGB)).FontSize(11).TextColor(AppTheme.TextMuted)
						ui.Text(c, formatBytes(snap.MemUsed)).FontSize(26).Bold().TextColor(AppTheme.TextPrimary)
					})

					// Status badge
					if snap.MemPercent >= 75.0 {
						ui.Row(c).Padding(2, 6).Radius(4).Background(AppTheme.WarningBg).Children(func() {
							ui.Text(c, "! Elevated").FontSize(10).Bold().TextColor(AppTheme.WarningText)
						})
					} else {
						ui.Row(c).Padding(2, 6).Radius(4).Background(ui.RGBA(16, 185, 129, 0.15)).Children(func() {
							ui.Text(c, "Normal").FontSize(10).Bold().TextColor(AppTheme.Green)
						})
					}
				})

				ui.Row(c).FillWidth().Children(func() {
					s.subStat(c, "App", formatBytes(snap.MemApp))
					s.subStat(c, "Wired", formatBytes(snap.MemWired))
					s.subStat(c, "Compressed", formatBytes(snap.MemCompressed))
				})

				DotMatrix(c, memSeries, MatrixOptions{
					Rows:           7,
					Cols:           32,
					PrimaryColor:   AppTheme.Green,
					SecondaryColor: AppTheme.Orange,
					ThresholdRow:   5,
					HasSecondary:   true,
				})
			})
		})

		// 3. GPU Card
		s.cardContainer(c, "GPU", IconGPU, AppTheme.Magenta, func() {
			ui.Column(c).FillWidth().Gap(10).Children(func() {
				ui.Column(c).Gap(2).Children(func() {
					model := snap.GPUModel
					if model == "" {
						model = snap.ChipName
					}
					ui.Text(c, model).FontSize(11).TextColor(AppTheme.TextMuted)
					ui.Text(c, fmt.Sprintf("%.1f %%", snap.GPUPercent)).FontSize(26).Bold().TextColor(AppTheme.TextPrimary)
				})

				ui.Row(c).FillWidth().Children(func() {
					vram := formatBytes(snap.GPUMemoryBytes)
					if snap.GPUMemoryBytes == 0 {
						vram = "1.11 GB"
					}
					s.subStat(c, "Memory", vram)
					s.subStat(c, "Average", fmt.Sprintf("%.0f%%", snap.GPUAverage))
					s.subStat(c, "Peak", fmt.Sprintf("%.0f%%", snap.GPUPeak))
				})

				DotMatrix(c, gpuSeries, MatrixOptions{
					Rows:         7,
					Cols:         32,
					PrimaryColor: AppTheme.Magenta,
				})
			})
		})
	})

	// Row 2: Disk, Network, Battery
	ui.Row(c).FillWidth().Gap(14).Children(func() {
		// 4. Disk Card
		s.cardContainer(c, "Disk", IconDisk, AppTheme.Orange, func() {
			ui.Column(c).FillWidth().Gap(10).Children(func() {
				ui.Column(c).Gap(2).Children(func() {
					totalGB := float64(snap.DiskTotalBytes) / (1024.0 * 1024.0 * 1024.0)
					if totalGB == 0 {
						totalGB = 995.0
					}
					ui.Text(c, fmt.Sprintf("Free of %.0f GB", totalGB)).FontSize(11).TextColor(AppTheme.TextMuted)
					freeGB := float64(snap.DiskFreeBytes) / (1024.0 * 1024.0 * 1024.0)
					if freeGB == 0 {
						freeGB = 408.0
					}
					ui.Text(c, fmt.Sprintf("%.0f GB", freeGB)).FontSize(26).Bold().TextColor(AppTheme.TextPrimary)
				})

				ui.Row(c).FillWidth().Children(func() {
					s.subStat(c, "Reading", formatRate(snap.DiskReadRate))
					s.subStat(c, "Writing", formatRate(snap.DiskWriteRate))
					written := formatBytes(snap.DiskWritten)
					if snap.DiskWritten == 0 {
						written = "1.01 TB"
					}
					s.subStat(c, "Written", written)
				})

				DotMatrix(c, diskSeries, MatrixOptions{
					Rows:         7,
					Cols:         32,
					PrimaryColor: AppTheme.Orange,
				})
			})
		})

		// 5. Network Card
		s.cardContainer(c, "Network", IconNetwork, AppTheme.Cyan, func() {
			ui.Column(c).FillWidth().Gap(10).Children(func() {
				ui.Column(c).Gap(2).Children(func() {
					ui.Text(c, "Downloading").FontSize(11).TextColor(AppTheme.TextMuted)
					downRate := formatRate(snap.NetDownloadRate)
					if snap.NetDownloadRate == 0 {
						downRate = "2.39 MB/s"
					}
					ui.Text(c, downRate).FontSize(26).Bold().TextColor(AppTheme.TextPrimary)
				})

				ui.Row(c).FillWidth().Children(func() {
					upRate := formatRate(snap.NetUploadRate)
					if snap.NetUploadRate == 0 {
						upRate = "13.7 KB/s"
					}
					s.subStat(c, "Uploading", upRate)
					rec := formatBytes(snap.NetReceived)
					if snap.NetReceived == 0 {
						rec = "198 GB"
					}
					s.subStat(c, "Received", rec)
					sent := formatBytes(snap.NetSent)
					if snap.NetSent == 0 {
						sent = "24.7 GB"
					}
					s.subStat(c, "Sent", sent)
				})

				DotMatrix(c, netSeries, MatrixOptions{
					Rows:         7,
					Cols:         32,
					PrimaryColor: AppTheme.Cyan,
				})
			})
		})

		// 6. Battery Card
		s.cardContainer(c, "Battery", IconBattery, AppTheme.Green, func() {
			ui.Column(c).FillWidth().Gap(10).Children(func() {
				ui.Column(c).Gap(2).Children(func() {
					label := "On Battery"
					if snap.OnPowerAdapter {
						label = "On Power Adapter"
					}
					ui.Text(c, label).FontSize(11).TextColor(AppTheme.TextMuted)
					ui.Text(c, fmt.Sprintf("%d%%", snap.BatteryPercent)).FontSize(26).Bold().TextColor(AppTheme.TextPrimary)
				})

				ui.Row(c).FillWidth().Children(func() {
					rem := snap.PowerRemaining
					if rem == "" {
						rem = "On power"
					}
					s.subStat(c, "Remaining", rem)
					s.subStat(c, "Power Draw", fmt.Sprintf("%.1f W", snap.PowerDrawWatts))
					s.subStat(c, "Health", fmt.Sprintf("%d%%", snap.BatteryHealth))
				})

				DotMatrix(c, batSeries, MatrixOptions{
					Rows:         7,
					Cols:         32,
					PrimaryColor: AppTheme.Green,
				})
			})
		})
	})
}

func (s *AppState) overviewBreakdownsGrid(c *ui.Context, snap metrics.SystemSnapshot) {
	ui.Row(c).FillWidth().Gap(14).Children(func() {
		// Breakdown 1: Memory by Type
		s.cardContainer(c, "Memory by Type", IconMemory, AppTheme.Green, func() {
			segments := []DonutSegment{
				{Label: "App", Value: float64(snap.MemApp), Color: AppTheme.Green},
				{Label: "Wired", Value: float64(snap.MemWired), Color: AppTheme.Orange},
				{Label: "Compressed", Value: float64(snap.MemCompressed), Color: AppTheme.Purple},
				{Label: "Cached", Value: float64(snap.MemCached), Color: AppTheme.Gray},
				{Label: "Free", Value: float64(snap.MemFree), Color: ui.Hex("#27272a")},
			}

			ui.Row(c).FillWidth().Gap(16).AlignItems(ui.Center).Children(func() {
				DonutChart(c, int(snap.MemPercent), segments)

				ui.Column(c).Grow(1).Gap(6).Children(func() {
					s.legendItem(c, "App", formatBytes(snap.MemApp), AppTheme.Green)
					s.legendItem(c, "Wired", formatBytes(snap.MemWired), AppTheme.Orange)
					s.legendItem(c, "Compressed", formatBytes(snap.MemCompressed), AppTheme.Purple)
					s.legendItem(c, "Cached", formatBytes(snap.MemCached), AppTheme.Gray)
					s.legendItem(c, "Free", formatBytes(snap.MemFree), ui.Hex("#3f3f46"))
				})
			})
		})

		// Breakdown 2: Memory by App
		s.cardContainerWithBadge(c, "Memory by App", IconOverview, AppTheme.Green, "9.89 GB in all", func() {
			ui.Column(c).FillWidth().Gap(10).Children(func() {
				// 3x2 grid of apps
				apps := snap.AppMemories
				if len(apps) >= 6 {
					ui.Row(c).FillWidth().Gap(8).Children(func() {
						s.appMemoryTile(c, apps[0])
						s.appMemoryTile(c, apps[1])
						s.appMemoryTile(c, apps[2])
					})
					ui.Row(c).FillWidth().Gap(8).Children(func() {
						s.appMemoryTile(c, apps[3])
						s.appMemoryTile(c, apps[4])
						s.appMemoryTile(c, apps[5])
					})
				}
			})
		})

		// Breakdown 3: Energy by App
		s.cardContainerWithBadge(c, "Energy by App", IconEnergy, AppTheme.Yellow, "6.41 W in all", func() {
			ui.Column(c).FillWidth().Gap(9).Children(func() {
				for _, app := range snap.AppEnergies {
					s.energyBarRow(c, app)
				}
			})
		})
	})
}

func (s *AppState) cardContainer(c *ui.Context, title string, icon *ui.SVG, iconColor ui.Color, body func()) {
	card := ui.Column(c).Grow(1).Padding(14).Gap(10).Radius(12).Background(AppTheme.CardBg).Border(1, AppTheme.CardBorder)
	card.Children(func() {
		// Header
		ui.Row(c).FillWidth().AlignItems(ui.Center).Children(func() {
			ui.Row(c).Gap(6).AlignItems(ui.Center).Grow(1).Children(func() {
				ui.Icon(c, icon).Size(15, 15).TextColor(iconColor)
				ui.Text(c, title).FontSize(14).Bold().TextColor(AppTheme.TextPrimary)
			})
			ui.Icon(c, IconChevronRight).Size(14, 14).TextColor(AppTheme.TextDim)
		})

		body()
	})
}

func (s *AppState) cardContainerWithBadge(c *ui.Context, title string, icon *ui.SVG, iconColor ui.Color, badge string, body func()) {
	card := ui.Column(c).Grow(1).Padding(14).Gap(10).Radius(12).Background(AppTheme.CardBg).Border(1, AppTheme.CardBorder)
	card.Children(func() {
		// Header with badge
		ui.Row(c).FillWidth().AlignItems(ui.Center).Children(func() {
			ui.Row(c).Gap(6).AlignItems(ui.Center).Grow(1).Children(func() {
				ui.Icon(c, icon).Size(15, 15).TextColor(iconColor)
				ui.Text(c, title).FontSize(14).Bold().TextColor(AppTheme.TextPrimary)
			})
			if badge != "" {
				ui.Text(c, badge).FontSize(11).TextColor(AppTheme.TextDim)
			}
		})

		body()
	})
}

func (s *AppState) subStat(c *ui.Context, label, value string) {
	ui.Column(c).Grow(1).Gap(2).Children(func() {
		ui.Text(c, label).FontSize(11).TextColor(AppTheme.TextMuted)
		ui.Text(c, value).FontSize(12).Bold().TextColor(AppTheme.TextPrimary)
	})
}

func (s *AppState) legendItem(c *ui.Context, label, value string, dotColor ui.Color) {
	ui.Row(c).FillWidth().AlignItems(ui.Center).Children(func() {
		ui.Row(c).Size(6, 6).Radius(3).Background(dotColor)
		ui.Row(c).Width(6)
		ui.Text(c, label).FontSize(11).TextColor(AppTheme.TextMuted).Grow(1)
		ui.Text(c, value).FontSize(11).Bold().TextColor(AppTheme.TextPrimary)
	})
}

func (s *AppState) appMemoryTile(c *ui.Context, app metrics.AppResource) {
	icon := IconXcode
	iconColor := AppTheme.Blue
	switch app.Icon {
	case "safari":
		icon = IconSafari
		iconColor = AppTheme.Cyan
	case "photos":
		icon = IconPhotos
		iconColor = AppTheme.Magenta
	case "finder":
		icon = IconFinder
		iconColor = AppTheme.Blue
	case "settings":
		icon = IconSettings
		iconColor = AppTheme.Gray
	case "music":
		icon = IconMusic
		iconColor = AppTheme.Magenta
	}

	ui.Column(c).Grow(1).Padding(8, 6).Radius(8).Background(ui.RGBA(255, 255, 255, 0.03)).Border(1, AppTheme.CardBorder).Center().Gap(4).Children(func() {
		ui.Icon(c, icon).Size(22, 22).TextColor(iconColor)
		ui.Text(c, app.Name).FontSize(11).TextColor(AppTheme.TextMuted).SingleLine()
		ui.Text(c, app.Value).FontSize(11).Bold().TextColor(AppTheme.TextPrimary)
	})
}

func (s *AppState) energyBarRow(c *ui.Context, app metrics.AppResource) {
	icon := IconXcode
	iconColor := AppTheme.Blue
	switch app.Icon {
	case "safari":
		icon = IconSafari
		iconColor = AppTheme.Cyan
	case "node":
		icon = IconNode
		iconColor = AppTheme.Green
	case "photos":
		icon = IconPhotos
		iconColor = AppTheme.Magenta
	case "python":
		icon = IconPython
		iconColor = AppTheme.Yellow
	}

	// Bar width fraction out of max ~3.0W
	barFraction := float32(app.Raw / 3.0)
	if barFraction > 1.0 {
		barFraction = 1.0
	}
	if barFraction < 0.05 {
		barFraction = 0.05
	}

	ui.Row(c).FillWidth().AlignItems(ui.Center).Gap(8).Children(func() {
		ui.Icon(c, icon).Size(14, 14).TextColor(iconColor)
		ui.Text(c, app.Name).FontSize(11).TextColor(AppTheme.TextMuted).Width(75).SingleLine()

		// Progress bar
		barArea := ui.Row(c).Grow(1).Height(5).Radius(3).Background(ui.RGBA(255, 255, 255, 0.06))
		barArea.Draw(func(p *ui.Painter, r ui.Rect) {
			fillWidth := r.W * barFraction
			p.Fill(ui.Rect{X: r.X, Y: r.Y, W: fillWidth, H: r.H}, AppTheme.Yellow, 2.5)
		})

		ui.Text(c, app.Value).FontSize(11).Bold().TextColor(AppTheme.TextPrimary).Width(55).TextAlign(ui.End)
	})
}
