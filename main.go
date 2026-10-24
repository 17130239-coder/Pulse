package main

import (
	"log"
	"time"

	"pulse/app"
	"pulse/metrics"

	"github.com/egoist/mygo"
	"github.com/egoist/mygo/ui"
)

func main() {
	collector := metrics.NewCollector()
	defer collector.Close()

	state := app.NewAppState(collector)

	mygo.App.WhenReady(func() {
		win := mygo.NewWindow(mygo.WindowOptions{
			Title:     "Pulse",
			Width:     1120,
			Height:    760,
			MinWidth:  900,
			MinHeight: 620,
			StateKey:  "pulse-main-window",
			Content:   ui.View(state.View),
		})
		state.Win = win

		// Real-time 1 Hz UI update loop
		go func() {
			ticker := time.NewTicker(1 * time.Second)
			defer ticker.Stop()

			for range ticker.C {
				if !state.IsPaused {
					state.UpdateSnapshot()
					win.Update(func() {
						// triggers smooth frame repaint
					})
				}
			}
		}()
	})

	if err := mygo.App.Run(); err != nil {
		log.Fatal(err)
	}
}
