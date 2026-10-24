package app

import (
	"strings"
	"sync"

	"pulse/metrics"

	"github.com/egoist/mygo"
	"github.com/egoist/mygo/ui"
)

// AppState coordinates the entire Pulse desktop interface.
type AppState struct {
	mu sync.RWMutex

	Win        *mygo.Window
	Collector  *metrics.Collector
	Snapshot   metrics.SystemSnapshot
	Router     *ui.Router

	// User Controls
	ActivePage  string // "Overview", "Processes", "CPU", "Memory", "GPU", "Disk", "Network", "Energy", "About This Mac"
	TimeRange   string // "1m", "5m", "15m", "1h"
	IsPaused    bool
	SearchQuery string

	// Process Table Sort
	SortField string // "name", "pid", "cpu", "mem"
	SortDesc  bool
}

// NewAppState creates and initializes the application state.
func NewAppState(collector *metrics.Collector) *AppState {
	s := &AppState{
		Collector:   collector,
		Snapshot:    collector.GetCurrent(),
		Router:      ui.NewRouter("/overview"),
		ActivePage:  "Overview",
		TimeRange:   "1m",
		IsPaused:    false,
		SortField:   "cpu",
		SortDesc:    true,
	}
	return s
}

// UpdateSnapshot refreshes the snapshot from collector.
func (s *AppState) UpdateSnapshot() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Snapshot = s.Collector.GetCurrent()
}

// SetTimeRange changes the active history window.
func (s *AppState) SetTimeRange(tr string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.TimeRange = tr
}

// TogglePause flips between paused and live sampling.
func (s *AppState) TogglePause() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.IsPaused = !s.IsPaused
	s.Collector.SetPaused(s.IsPaused)
}

// SetPage navigates to another page.
func (s *AppState) SetPage(page string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.ActivePage = page
	path := "/overview"
	switch page {
	case "Processes":
		path = "/processes"
	case "CPU":
		path = "/cpu"
	case "Memory":
		path = "/memory"
	case "GPU":
		path = "/gpu"
	case "Disk":
		path = "/disk"
	case "Network":
		path = "/network"
	case "Energy":
		path = "/energy"
	case "About This Mac":
		path = "/about"
	}
	s.Router.Push(path)
}

// GetFilteredProcesses returns processes matching the search query.
func (s *AppState) GetFilteredProcesses() []metrics.ProcessItem {
	s.mu.RLock()
	defer s.mu.RUnlock()

	list := s.Snapshot.Processes
	if s.SearchQuery == "" {
		return list
	}

	q := strings.ToLower(s.SearchQuery)
	var res []metrics.ProcessItem
	for _, p := range list {
		if strings.Contains(strings.ToLower(p.Name), q) || strings.Contains(strings.ToLower(p.Command), q) {
			res = append(res, p)
		}
	}
	return res
}
