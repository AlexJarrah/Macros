package internal

import "time"

// Macro structure to store key events with timing
type Macro struct {
	Key         string
	PressTime   time.Time
	ReleaseTime time.Time
}

type Device uint8
