package app

import "github.com/egoist/mygo/ui"

// Theme holds the color palette matching Pulse's dark macOS interface.
type Theme struct {
	Bg          ui.Color
	SidebarBg   ui.Color
	CardBg      ui.Color
	CardBorder  ui.Color
	HoverBg     ui.Color
	ActivePill  ui.Color

	TextPrimary ui.Color
	TextMuted   ui.Color
	TextDim     ui.Color

	// Accents
	Blue        ui.Color
	Green       ui.Color
	Magenta     ui.Color
	Orange      ui.Color
	Cyan        ui.Color
	Yellow      ui.Color
	Purple      ui.Color
	Gray        ui.Color

	// Badges
	WarningBg   ui.Color
	WarningText ui.Color

	// Dot Matrix Unlit
	DotUnlit    ui.Color
}

var AppTheme = Theme{
	Bg:          ui.Hex("#0e0e11"),
	SidebarBg:   ui.Hex("#141418"),
	CardBg:      ui.Hex("#1a1a20"),
	CardBorder:  ui.Hex("#26262e"),
	HoverBg:     ui.Hex("#22222a"),
	ActivePill:  ui.Hex("#2563eb"),

	TextPrimary: ui.Hex("#f4f4f6"),
	TextMuted:   ui.Hex("#8e8e98"),
	TextDim:     ui.Hex("#555560"),

	Blue:        ui.Hex("#3b82f6"),
	Green:       ui.Hex("#10b981"),
	Magenta:     ui.Hex("#ec4899"),
	Orange:      ui.Hex("#f97316"),
	Cyan:        ui.Hex("#06b6d4"),
	Yellow:      ui.Hex("#eab308"),
	Purple:      ui.Hex("#a855f7"),
	Gray:        ui.Hex("#64748b"),

	WarningBg:   ui.Hex("#3a2505"),
	WarningText: ui.Hex("#f59e0b"),

	DotUnlit:    ui.Hex("#202028"),
}
