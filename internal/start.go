package internal

import (
	"log"
	"time"

	"github.com/MarinX/keylogger"
)

func Start() {
	keyboards, err := GetDevices(Keyboard)
	if err != nil {
		log.Fatalf("failed to get keyboard: %v", err)
	}

	keyboard := keyboards[0]
	log.Println("Keyboards detected:", keyboards)
	log.Println("Using keyboard:", keyboard)

	k, err := keylogger.New(keyboard)
	if err != nil {
		log.Fatalf("failed to create keylogger: %v", err)
	}
	defer k.Close()

	events := k.Read()
	var isRecording, isModifierDown bool
	var currentMacro []Macro
	var macroKey uint16

	for e := range events {
		if e.Type != keylogger.EvKey {
			continue
		}

		switch {
		case e.KeyPress():
			handleKeyPress(k, e, &isRecording, &isModifierDown, &currentMacro, &macroKey)

		case e.KeyRelease():
			handleKeyRelease(e, &isModifierDown, &isRecording, &currentMacro)
		}
	}
}

func handleKeyPress(k *keylogger.KeyLogger, e keylogger.InputEvent, recording *bool, isModifierDown *bool, currentMacro *[]Macro, macroKey *uint16) {
	log.Printf("[EVENT] KEY PRESSED:  %d (%s)", e.Code, e.KeyString())

	if e.Code == CANCEL_KEY {
		*recording = false
		*currentMacro = []Macro{}
		log.Println("CANCELLED RECORDING")
		return
	}

	if *recording {
		*currentMacro = append(*currentMacro, Macro{Key: e.KeyString(), PressTime: time.Now()})
	}

	if e.Code == MODIFIER_CODE {
		*isModifierDown = true
		return
	}

	if *isModifierDown {
		handleModifierKeyPress(k, e, recording, currentMacro, macroKey)
	}
}

func handleModifierKeyPress(k *keylogger.KeyLogger, e keylogger.InputEvent, recording *bool, currentMacro *[]Macro, macroKey *uint16) {
	switch e.Code {
	case REGISTER_KEY:
		if !*recording {
			startRecording(k, e, recording, currentMacro, macroKey)
		} else {
			stopRecording(recording, currentMacro, *macroKey)
		}
	case LOAD_KEY:
		playMacro(k)
	}
}

func startRecording(k *keylogger.KeyLogger, e keylogger.InputEvent, recording *bool, currentMacro *[]Macro, macroKey *uint16) {
	*recording = true
	*currentMacro = []Macro{}
	log.Println("Press the key to bind the macro to.")
	*macroKey = waitForKeyPress(k)
	log.Println("STARTED RECORDING MACRO:", e.KeyString(), *macroKey)
}

func stopRecording(recording *bool, currentMacro *[]Macro, macroKey uint16) {
	*recording = false
	REGISTERED_MACROS[macroKey] = *currentMacro
	log.Println("SAVED MACRO TO KEY:", macroKey)
}

func playMacro(k *keylogger.KeyLogger) {
	log.Println("Press the key to play the macro from.")
	macroKey := waitForKeyPress(k)
	log.Println("[MACRO EVENT] PLAYING MACRO:", macroKey)
	if macro, ok := REGISTERED_MACROS[macroKey]; ok {
		PlaybackMacro(k, macro)
	} else {
		log.Println("No macro found for key:", macroKey)
	}
}

func waitForKeyPress(k *keylogger.KeyLogger) uint16 {
	for e := range k.Read() {
		if e.Type == keylogger.EvKey && e.KeyPress() && e.Code != MODIFIER_CODE && e.Code != REGISTER_KEY && e.Code != LOAD_KEY {
			return e.Code
		}
	}
	return 0
}

func handleKeyRelease(e keylogger.InputEvent, isModifierDown *bool, recording *bool, currentMacro *[]Macro) {
	log.Printf("[EVENT] KEY RELEASED: %d (%s)", e.Code, e.KeyString())

	if e.Code == MODIFIER_CODE {
		*isModifierDown = false
	}

	if *recording {
		updateMacroReleaseTime(currentMacro, e.KeyString())
	}
}

func updateMacroReleaseTime(currentMacro *[]Macro, key string) {
	for i := range *currentMacro {
		if (*currentMacro)[i].Key == key && (*currentMacro)[i].ReleaseTime.IsZero() {
			(*currentMacro)[i].ReleaseTime = time.Now()
			break
		}
	}
}
