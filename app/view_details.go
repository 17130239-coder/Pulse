package app

import (
	"fmt"
	"runtime"

	"pulse/metrics"

	"github.com/egoist/mygo/ui"
)

// ViewDetails renders detailed metrics for a specific hardware subsystem.
func (s *AppState) ViewDetails(c *ui.Context, category string) {
	snap := s.Snapshot

	ui.Column(c).Grow(1).FillWidth().Padding(20, 28, 28).Gap(16).Children(func() {
		// Header
		ui.Row(c).FillWidth().AlignItems(ui.Center).Children(func() {
			ui.Column(c).Grow(1).Gap(2).Children(func() {
				ui.Text(c, category).FontSize(22).Bold().TextColor(AppTheme.TextPrimary)
				ui.Text(c, fmt.Sprintf("Detailed telemetry for %s", category)).FontSize(13).TextColor(AppTheme.TextMuted)
			})

			backBtn := ui.Row(c).Padding(5, 12).Radius(6).Background(AppTheme.CardBg).Border(1, AppTheme.CardBorder).AlignItems(ui.Center)
			if backBtn.Clicked() {
				s.SetPage("Overview")
			}
			backBtn.Children(func() {
				ui.Text(c, "← Back to Overview").FontSize(12).TextColor(AppTheme.TextPrimary)
			})
		})

		switch category {
		case "CPU":
			s.viewCPUDetails(c, snap)
		case "Memory":
			s.viewMemoryDetails(c, snap)
		case "GPU":
			s.viewGPUDetails(c, snap)
		case "Disk":
			s.viewDiskDetails(c, snap)
		case "Network":
			s.viewNetworkDetails(c, snap)
		case "Energy":
			s.viewEnergyDetails(c, snap)
		case "About This Mac":
			s.viewAboutDetails(c, snap)
		}
	})
}

func (s *AppState) viewCPUDetails(c *ui.Context, snap metrics.SystemSnapshot) {
	history := s.Collector.GetHistory(s.TimeRange)
	cpuSeries := make([]float64, len(history))
	for i, h := range history {
		cpuSeries[i] = h.CPUPercent / 100.0
	}

	ui.Column(c).FillWidth().Gap(14).Children(func() {
		s.cardContainer(c, "CPU Load History", IconCPU, AppTheme.Blue, func() {
			ui.Column(c).FillWidth().Gap(12).Children(func() {
				DotMatrix(c, cpuSeries, MatrixOptions{
					Rows:         10,
					Cols:         48,
					PrimaryColor: AppTheme.Blue,
				})
			})
		})

		ui.Row(c).FillWidth().Gap(14).Children(func() {
			s.detailMetricCard(c, "Current Utilization", fmt.Sprintf("%.1f%%", snap.CPUPercent), AppTheme.Blue)
			s.detailMetricCard(c, "User Space", fmt.Sprintf("%.1f%%", snap.CPUUser), AppTheme.Cyan)
			s.detailMetricCard(c, "System Kernel", fmt.Sprintf("%.1f%%", snap.CPUSystem), AppTheme.Orange)
			s.detailMetricCard(c, "Die Temperature", fmt.Sprintf("%.1f °C", snap.CPUTemperature), AppTheme.Yellow)
		})
	})
}

func (s *AppState) viewMemoryDetails(c *ui.Context, snap metrics.SystemSnapshot) {
	ui.Column(c).FillWidth().Gap(14).Children(func() {
		ui.Row(c).FillWidth().Gap(14).Children(func() {
			s.detailMetricCard(c, "Total Installed", formatBytes(snap.MemTotal), AppTheme.Green)
			s.detailMetricCard(c, "Used Memory", formatBytes(snap.MemUsed), AppTheme.Blue)
			s.detailMetricCard(c, "App Memory", formatBytes(snap.MemApp), AppTheme.Green)
			s.detailMetricCard(c, "Wired Memory", formatBytes(snap.MemWired), AppTheme.Orange)
		})

		ui.Row(c).FillWidth().Gap(14).Children(func() {
			s.detailMetricCard(c, "Compressed Memory", formatBytes(snap.MemCompressed), AppTheme.Purple)
			s.detailMetricCard(c, "Cached Files", formatBytes(snap.MemCached), AppTheme.Gray)
			s.detailMetricCard(c, "Free Pages", formatBytes(snap.MemFree), AppTheme.TextMuted)
			s.detailMetricCard(c, "Memory Pressure", fmt.Sprintf("%.1f%%", snap.MemPercent), AppTheme.Yellow)
		})
	})
}

func (s *AppState) viewGPUDetails(c *ui.Context, snap metrics.SystemSnapshot) {
	ui.Column(c).FillWidth().Gap(14).Children(func() {
		ui.Row(c).FillWidth().Gap(14).Children(func() {
			s.detailMetricCard(c, "GPU Model", snap.GPUModel, AppTheme.Magenta)
			s.detailMetricCard(c, "Utilization", fmt.Sprintf("%.1f%%", snap.GPUPercent), AppTheme.Magenta)
			s.detailMetricCard(c, "Allocated VRAM", formatBytes(snap.GPUMemoryBytes), AppTheme.Cyan)
			s.detailMetricCard(c, "Peak Utilization", fmt.Sprintf("%.1f%%", snap.GPUPeak), AppTheme.Yellow)
		})
	})
}

func (s *AppState) viewDiskDetails(c *ui.Context, snap metrics.SystemSnapshot) {
	ui.Column(c).FillWidth().Gap(14).Children(func() {
		ui.Row(c).FillWidth().Gap(14).Children(func() {
			s.detailMetricCard(c, "Total Storage", formatBytes(snap.DiskTotalBytes), AppTheme.Orange)
			s.detailMetricCard(c, "Free Storage", formatBytes(snap.DiskFreeBytes), AppTheme.Green)
			s.detailMetricCard(c, "Read Speed", formatRate(snap.DiskReadRate), AppTheme.Cyan)
			s.detailMetricCard(c, "Write Speed", formatRate(snap.DiskWriteRate), AppTheme.Orange)
		})
	})
}

func (s *AppState) viewNetworkDetails(c *ui.Context, snap metrics.SystemSnapshot) {
	ui.Column(c).FillWidth().Gap(14).Children(func() {
		ui.Row(c).FillWidth().Gap(14).Children(func() {
			s.detailMetricCard(c, "Download Rate", formatRate(snap.NetDownloadRate), AppTheme.Cyan)
			s.detailMetricCard(c, "Upload Rate", formatRate(snap.NetUploadRate), AppTheme.Blue)
			s.detailMetricCard(c, "Total Received", formatBytes(snap.NetReceived), AppTheme.Cyan)
			s.detailMetricCard(c, "Total Sent", formatBytes(snap.NetSent), AppTheme.Blue)
		})
	})
}

func (s *AppState) viewEnergyDetails(c *ui.Context, snap metrics.SystemSnapshot) {
	ui.Column(c).FillWidth().Gap(14).Children(func() {
		ui.Row(c).FillWidth().Gap(14).Children(func() {
			s.detailMetricCard(c, "Charge Level", fmt.Sprintf("%d%%", snap.BatteryPercent), AppTheme.Green)
			s.detailMetricCard(c, "Power State", snap.PowerRemaining, AppTheme.Yellow)
			s.detailMetricCard(c, "Power Draw", fmt.Sprintf("%.1f W", snap.PowerDrawWatts), AppTheme.Orange)
			s.detailMetricCard(c, "Battery Health", fmt.Sprintf("%d%%", snap.BatteryHealth), AppTheme.Green)
		})
	})
}

func (s *AppState) viewAboutDetails(c *ui.Context, snap metrics.SystemSnapshot) {
	ui.Column(c).FillWidth().Gap(14).Children(func() {
		s.cardContainer(c, "About This Mac", IconAbout, AppTheme.Blue, func() {
			ui.Column(c).FillWidth().Gap(10).Children(func() {
				s.aboutRow(c, "Model", snap.ModelName)
				s.aboutRow(c, "Processor", snap.ChipName)
				s.aboutRow(c, "Memory", fmt.Sprintf("%d GB Unified Memory", snap.MemTotal/(1024*1024*1024)))
				s.aboutRow(c, "Architecture", runtime.GOARCH)
				s.aboutRow(c, "OS Platform", "macOS Darwin")
				s.aboutRow(c, "System Uptime", snap.Uptime)
				s.aboutRow(c, "Framework", "MyGo (Pure Go, Metal GPU Rendering)")
			})
		})
	})
}

func (s *AppState) detailMetricCard(c *ui.Context, label, value string, accent ui.Color) {
	ui.Column(c).Grow(1).Padding(16).Radius(10).Background(AppTheme.CardBg).Border(1, AppTheme.CardBorder).Gap(4).Children(func() {
		ui.Text(c, label).FontSize(11).TextColor(AppTheme.TextMuted)
		ui.Text(c, value).FontSize(20).Bold().TextColor(accent)
	})
}

func (s *AppState) aboutRow(c *ui.Context, label, value string) {
	ui.Row(c).FillWidth().Padding(6, 0).BorderWidth(0, 0, 1, 0).BorderColor(ui.RGBA(255, 255, 255, 0.03)).Children(func() {
		ui.Text(c, label).FontSize(13).Bold().TextColor(AppTheme.TextMuted).Width(150)
		ui.Text(c, value).FontSize(13).TextColor(AppTheme.TextPrimary).Grow(1)
	})
}
