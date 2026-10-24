package metrics

import (
	"sync"
	"time"
)

// SystemSnapshot captures all hardware and OS metrics at one instant.
type SystemSnapshot struct {
	Timestamp time.Time

	// CPU
	CPUPercent     float64 // Total CPU usage 0..100
	CPUUser        float64 // User %
	CPUSystem      float64 // System %
	CPUTemperature float64 // Celsius

	// Memory (in Bytes)
	MemTotal      uint64
	MemUsed       uint64
	MemApp        uint64
	MemWired      uint64
	MemCompressed uint64
	MemCached     uint64
	MemFree       uint64
	MemPercent    float64 // 0..100

	// GPU
	GPUModel       string
	GPUPercent     float64 // 0..100
	GPUMemoryBytes uint64
	GPUAverage     float64
	GPUPeak        float64

	// Disk
	DiskTotalBytes uint64
	DiskFreeBytes  uint64
	DiskReadRate   uint64 // Bytes/sec
	DiskWriteRate  uint64 // Bytes/sec
	DiskWritten    uint64 // Total Bytes written

	// Network
	NetDownloadRate uint64 // Bytes/sec
	NetUploadRate   uint64 // Bytes/sec
	NetReceived     uint64 // Total Bytes received
	NetSent         uint64 // Total Bytes sent

	// Battery / Energy
	BatteryPercent  int
	OnPowerAdapter  bool
	PowerRemaining  string  // "On power" or "3h 42m"
	PowerDrawWatts  float64 // Instantaneous power draw
	BatteryHealth   int     // 0..100 %

	// Machine Info
	ModelName string // "MacBook Air" / "MacBook Pro"
	ChipName  string // "Apple M1 Pro"
	Uptime    string // "Up 14d 2h"

	// Top Processes for Breakdowns
	AppMemories []AppResource
	AppEnergies []AppResource
	Processes   []ProcessItem
}

// AppResource represents an application's resource consumption.
type AppResource struct {
	Name  string
	Value string
	Raw   float64 // Bytes or Watts
	Icon  string  // App identifier for icon selection
}

// ProcessItem represents a running process.
type ProcessItem struct {
	PID     int
	Name    string
	User    string
	CPU     float64
	Mem     float64
	MemMB   uint64
	Status  string
	Command string
}

// Collector continuously samples metrics and maintains history rings.
type Collector struct {
	mu        sync.RWMutex
	current   SystemSnapshot
	history   map[string][]SystemSnapshot // "1m", "5m", "15m", "1h"
	paused    bool
	done      chan struct{}
	lastNetIn uint64
	lastNetOut uint64
	lastDiskRead uint64
	lastDiskWrite uint64
	lastSampleTime time.Time

	// Rolling stats for GPU
	gpuSamples []float64
}

// NewCollector creates and starts the metrics collection engine.
func NewCollector() *Collector {
	c := &Collector{
		history: map[string][]SystemSnapshot{
			"1m":  make([]SystemSnapshot, 0, 60),
			"5m":  make([]SystemSnapshot, 0, 60),
			"15m": make([]SystemSnapshot, 0, 60),
			"1h":  make([]SystemSnapshot, 0, 60),
		},
		done: make(chan struct{}),
	}

	// Initial synchronous sample
	c.current = c.sampleOnce()
	c.recordHistory(c.current)

	// Start background sampling goroutine (1 Hz)
	go c.run()

	return c
}

// Close stops the collector.
func (c *Collector) Close() {
	close(c.done)
}

// SetPaused toggles live updates.
func (c *Collector) SetPaused(paused bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.paused = paused
}

// IsPaused reports whether collection is paused.
func (c *Collector) IsPaused() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.paused
}

// GetCurrent returns a copy of the latest snapshot.
func (c *Collector) GetCurrent() SystemSnapshot {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.current
}

// GetHistory returns the snapshots for the given duration window.
func (c *Collector) GetHistory(window string) []SystemSnapshot {
	c.mu.RLock()
	defer c.mu.RUnlock()
	h := c.history[window]
	res := make([]SystemSnapshot, len(h))
	copy(res, h)
	return res
}

func (c *Collector) run() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-c.done:
			return
		case <-ticker.C:
			c.mu.RLock()
			paused := c.paused
			c.mu.RUnlock()

			if !paused {
				snap := c.sampleOnce()
				c.mu.Lock()
				c.current = snap
				c.recordHistory(snap)
				c.mu.Unlock()
			}
		}
	}
}

func (c *Collector) recordHistory(snap SystemSnapshot) {
	// 1m: keep last 60 1-second samples
	c.history["1m"] = appendRing(c.history["1m"], snap, 60)

	// 5m: keep last 60 samples (every 5 seconds)
	if len(c.history["1m"])%5 == 0 {
		c.history["5m"] = appendRing(c.history["5m"], snap, 60)
	}

	// 15m: keep last 60 samples (every 15 seconds)
	if len(c.history["1m"])%15 == 0 {
		c.history["15m"] = appendRing(c.history["15m"], snap, 60)
	}

	// 1h: keep last 60 samples (every 60 seconds)
	if len(c.history["1m"])%60 == 0 {
		c.history["1h"] = appendRing(c.history["1h"], snap, 60)
	}
}

func appendRing(slice []SystemSnapshot, item SystemSnapshot, maxCap int) []SystemSnapshot {
	if len(slice) >= maxCap {
		slice = slice[1:]
	}
	return append(slice, item)
}
