package internal

const (
	MODIFIER_CODE uint16 = 100 // Right ALT
	REGISTER_KEY  uint16 = 4   // 3
	LOAD_KEY      uint16 = 3   // 2
	CANCEL_KEY    uint16 = 2   // 1
)

const (
	Keyboard Device = iota
	Mouse
)

// Map to store recorded macros
var REGISTERED_MACROS = make(map[uint16][]Macro)
