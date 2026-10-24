package metrics

import (
	"testing"
	"time"
)

func TestCollector(t *testing.T) {
	c := NewCollector()
	defer c.Close()

	time.Sleep(100 * time.Millisecond)
	snap := c.GetCurrent()

	t.Logf("Model: %s, Chip: %s, Uptime: %s", snap.ModelName, snap.ChipName, snap.Uptime)
	t.Logf("CPU: %.1f%% (User: %.1f%%, Sys: %.1f%%), Temp: %.1f °C", snap.CPUPercent, snap.CPUUser, snap.CPUSystem, snap.CPUTemperature)
	t.Logf("Memory: Total: %d GB, Used: %d GB (%.1f%%)", snap.MemTotal/(1024*1024*1024), snap.MemUsed/(1024*1024*1024), snap.MemPercent)
	t.Logf("GPU: %s (%.1f%%)", snap.GPUModel, snap.GPUPercent)
	t.Logf("Battery: %d%%, OnPower: %v, Health: %d%%, PowerDraw: %.1f W", snap.BatteryPercent, snap.OnPowerAdapter, snap.BatteryHealth, snap.PowerDrawWatts)

	if snap.MemTotal == 0 {
		t.Errorf("expected total memory > 0")
	}
}
