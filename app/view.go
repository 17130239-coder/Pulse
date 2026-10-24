package app

import "github.com/egoist/mygo/ui"

// View is the top-level root view function for the MyGo window.
func (s *AppState) View(c *ui.Context) {
	// Root horizontal layout: Sidebar on the left, Main content area on the right
	ui.Row(c).Fill().Background(AppTheme.Bg).AlignItems(ui.Stretch).Children(func() {
		// Left Navigation Sidebar
		s.ViewSidebar(c)

		// Main Content Area with Scroll
		ui.Column(c).Grow(1).FillHeight().Children(func() {
			ui.Scroll(c).Grow(1).FillWidth().Children(func() {
				s.Router.View(c, func(r *ui.Route) {
					switch r.Path() {
					case "/overview":
						r.Title("Overview — Pulse")
						s.ViewOverview(c)
					case "/processes":
						r.Title("Processes — Pulse")
						s.ViewProcesses(c)
					case "/cpu":
						r.Title("CPU — Pulse")
						s.ViewDetails(c, "CPU")
					case "/memory":
						r.Title("Memory — Pulse")
						s.ViewDetails(c, "Memory")
					case "/gpu":
						r.Title("GPU — Pulse")
						s.ViewDetails(c, "GPU")
					case "/disk":
						r.Title("Disk — Pulse")
						s.ViewDetails(c, "Disk")
					case "/network":
						r.Title("Network — Pulse")
						s.ViewDetails(c, "Network")
					case "/energy":
						r.Title("Energy — Pulse")
						s.ViewDetails(c, "Energy")
					case "/about":
						r.Title("About This Mac — Pulse")
						s.ViewDetails(c, "About This Mac")
					default:
						r.Title("Overview — Pulse")
						s.ViewOverview(c)
					}
				})
			})
		})
	})
}
