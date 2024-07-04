package internal

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

// Returns all devices of the specified type. An error is returned if no devices are found.
func GetDevices(type_ Device) (devices []string, err error) {
	const DEVICE_INPUT_ID_DIR = "/dev/input/by-id"
	const INTERFACE_IDENTIFIER_REGEX = `-if\d{2}`
	const KEYBOARD_IDENTIFIER = "event-kbd"
	const MOUSE_IDENTIFIER = "event-mouse"

	files, err := os.ReadDir(DEVICE_INPUT_ID_DIR)
	if err != nil {
		return nil, fmt.Errorf("error reading directory %s: %s", DEVICE_INPUT_ID_DIR, err.Error())
	}

	for _, f := range files {
		isAlternateInterface, _ := regexp.MatchString(INTERFACE_IDENTIFIER_REGEX, f.Name())
		if isAlternateInterface {
			continue
		}

		var valid bool
		switch type_ {
		case Keyboard:
			valid = strings.Contains(f.Name(), KEYBOARD_IDENTIFIER)

		case Mouse:
			valid = strings.Contains(f.Name(), MOUSE_IDENTIFIER)
		}

		if valid {
			devPath := fmt.Sprintf("%s/%s", DEVICE_INPUT_ID_DIR, f.Name())
			devices = append(devices, devPath)
		}
	}

	if len(devices) == 0 {
		return nil, fmt.Errorf("no devices found")
	}

	return devices, nil
}
