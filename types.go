package ecwid

type (
	// Custom fields

	// DateTime some like "2015-09-20 19:59:43 +0000"
	DateTime string // TODO DateTime => time.Date + Marshal|Unmarshal|String

	// ModifierType - price modifier type
	ModifierType string
)

// ModifierType types
const (
	ModifierPercent  ModifierType = "PERCENT"
	ModifierAbsolute ModifierType = "ABSOLUTE"
)
