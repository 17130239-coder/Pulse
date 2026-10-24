package app

import (
	"fmt"

	"github.com/egoist/mygo/ui"
)

// ViewSidebar renders the left navigation sidebar.
func (s *AppState) ViewSidebar(c *ui.Context) {
	snap := s.Snapshot

	ui.Column(c).Width(215).FillHeight().Padding(14, 12).Gap(4).Background(AppTheme.SidebarBg).BorderWidth(0, 1, 0, 0).BorderColor(AppTheme.CardBorder).Children(func() {
		// Window top spacing
		ui.Row(c).Height(8)

		// 1. Primary Pages
		s.sidebarNavButton(c, "Overview", IconOverview, "", s.ActivePage == "Overview")
		procCount := fmt.Sprintf("%d", len(snap.Processes))
		if len(snap.Processes) == 0 {
			procCount = "24"
		}
		s.sidebarNavButton(c, "Processes", IconProcesses, procCount, s.ActivePage == "Processes")

		ui.Row(c).Height(14)

		// 2. Hardware Section
		ui.Text(c, "Hardware").FontSize(11).Bold().TextColor(AppTheme.TextDim).Padding(4, 8, 4)

		cpuVal := fmt.Sprintf("%.0f%%", snap.CPUPercent)
		s.sidebarStatButton(c, "CPU", IconCPU, cpuVal, AppTheme.Blue, s.ActivePage == "CPU")

		memVal := fmt.Sprintf("%.0f%%", snap.MemPercent)
		s.sidebarStatButton(c, "Memory", IconMemory, memVal, AppTheme.Green, s.ActivePage == "Memory")

		gpuVal := fmt.Sprintf("%.1f%%", snap.GPUPercent)
		s.sidebarStatButton(c, "GPU", IconGPU, gpuVal, AppTheme.Magenta, s.ActivePage == "GPU")

		diskVal := formatRate(snap.DiskReadRate + snap.DiskWriteRate)
		s.sidebarStatButton(c, "Disk", IconDisk, diskVal, AppTheme.Orange, s.ActivePage == "Disk")

		netVal := formatRate(snap.NetDownloadRate)
		s.sidebarStatButton(c, "Network", IconNetwork, netVal, AppTheme.Cyan, s.ActivePage == "Network")

		energyVal := fmt.Sprintf("%d%%", snap.BatteryPercent)
		s.sidebarStatButton(c, "Energy", IconEnergy, energyVal, AppTheme.Yellow, s.ActivePage == "Energy")

		ui.Row(c).Height(14)

		// 3. System Section
		ui.Text(c, "System").FontSize(11).Bold().TextColor(AppTheme.TextDim).Padding(4, 8, 4)
		s.sidebarNavButton(c, "About This Mac", IconAbout, "", s.ActivePage == "About This Mac")

		// Spacer to push Mac model info to bottom
		ui.Column(c).Grow(1)

		// 4. Bottom Device Info Card
		ui.Column(c).FillWidth().Padding(8, 10).Gap(3).Radius(8).Background(AppTheme.CardBg).Border(1, AppTheme.CardBorder).Children(func() {
			model := snap.ModelName
			if model == "" {
				model = "MacBook Air"
			}
			ui.Text(c, model).FontSize(13).Bold().TextColor(AppTheme.TextPrimary)

			chip := snap.ChipName
			if chip == "" {
				chip = "Apple Silicon"
			}
			up := snap.Uptime
			if up == "" {
				up = "Up 14d 2h"
			}
			ui.Text(c, fmt.Sprintf("%s · %s", chip, up)).FontSize(11).TextColor(AppTheme.TextMuted)
		})
	})
}

func (s *AppState) sidebarNavButton(c *ui.Context, label string, icon *ui.SVG, badge string, active bool) {
	bg := ui.RGBA(0, 0, 0, 0)
	textColor := AppTheme.TextMuted
	if active {
		bg = AppTheme.ActivePill
		textColor = AppTheme.TextPrimary
	}

	btn := ui.Row(c).FillWidth().Padding(7, 10).Radius(7).Background(bg).Gap(10).AlignItems(ui.Center)
	if btn.Clicked() {
		s.SetPage(label)
	}

	btn.Children(func() {
		ui.Icon(c, icon).Size(16, 16).TextColor(textColor)
		ui.Text(c, label).FontSize(13).TextColor(textColor).Grow(1)

		if badge != "" {
			ui.Row(c).Padding(2, 6).Radius(10).Background(ui.RGBA(255, 255, 255, 0.12)).Children(func() {
				ui.Text(c, badge).FontSize(11).Bold().TextColor(AppTheme.TextPrimary)
			})
		}
	})
}

func (s *AppState) sidebarStatButton(c *ui.Context, label string, icon *ui.SVG, value string, iconColor ui.Color, active bool) {
	bg := ui.RGBA(0, 0, 0, 0)
	textColor := AppTheme.TextMuted
	if active {
		bg = AppTheme.HoverBg
		textColor = AppTheme.TextPrimary
	}

	btn := ui.Row(c).FillWidth().Padding(7, 10).Radius(7).Background(bg).Gap(10).AlignItems(ui.Center)
	if btn.Clicked() {
		s.SetPage(label)
	}

	btn.Children(func() {
		ui.Icon(c, icon).Size(16, 16).TextColor(iconColor)
		ui.Text(c, label).FontSize(13).TextColor(textColor).Grow(1)
		ui.Text(c, value).FontSize(12).TextColor(AppTheme.TextDim)
	})
}

func formatRate(bytesPerSec uint64) string {
	if bytesPerSec >= 1024*1024 {
		return fmt.Sprintf("%.2f MB/s", float64(bytesPerSec)/(1024.0*1024.0))
	} else if bytesPerSec >= 1024 {
		return fmt.Sprintf("%.0f KB/s", float64(bytesPerSec)/1024.0)
	}
	return fmt.Sprintf("%d B/s", bytesPerSec)
}

func formatBytes(b uint64) string {
	gb := float64(b) / (1024.0 * 1024.0 * 1024.0)
	if gb >= 1000.0 {
		return fmt.Sprintf("%.2f TB", gb/1024.0)
	}
	if gb >= 1.0 {
		return fmt.Sprintf("%.1f GB", gb)
	}
	mb := float64(b) / (1024.0 * 1024.0)
	return fmt.Sprintf("%.0f MB", mb)
}
