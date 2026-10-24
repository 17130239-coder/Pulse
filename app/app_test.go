package app

import (
	"testing"

	"pulse/metrics"
)

func TestAppState(t *testing.T) {
	collector := metrics.NewCollector()
	defer collector.Close()

	state := NewAppState(collector)

	if state.ActivePage != "Overview" {
		t.Errorf("expected active page Overview, got %s", state.ActivePage)
	}

	state.SetTimeRange("5m")
	if state.TimeRange != "5m" {
		t.Errorf("expected time range 5m, got %s", state.TimeRange)
	}

	state.TogglePause()
	if !state.IsPaused {
		t.Errorf("expected IsPaused to be true")
	}

	state.SetPage("Processes")
	if state.ActivePage != "Processes" {
		t.Errorf("expected active page Processes, got %s", state.ActivePage)
	}

	state.SearchQuery = "test"
	filtered := state.GetFilteredProcesses()
	t.Logf("found %d processes matching query", len(filtered))
}
