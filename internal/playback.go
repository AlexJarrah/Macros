package internal

import (
	"log"
	"time"

	"github.com/MarinX/keylogger"
)

func PlaybackMacro(k *keylogger.KeyLogger, macros []Macro) {
	for i, m := range macros {
		go func() {
			pressDuration := m.ReleaseTime.Sub(m.PressTime)
			log.Println("PLAYING KEY ", m.Key, " for duration ", pressDuration)

			k.Write(keylogger.KeyPress, m.Key)
			log.Println("[MACRO EVENT] KEY PRESSED:", m.Key)

			time.Sleep(pressDuration)

			k.Write(keylogger.KeyRelease, m.Key)
			log.Println("[MACRO EVENT] KEY RELEASED:", m.Key)
		}()

		isLastKeypress := len(macros)-1 == i
		if !isLastKeypress {
			nextKeypress := macros[i+1]
			waitDuration := nextKeypress.PressTime.Sub(m.PressTime)
			time.Sleep(waitDuration)
		}
	}
}
