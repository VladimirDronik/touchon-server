package Onokom

import (
	"touchon-server/internal/objects"
	"touchon-server/lib/mqtt/messages"
)

type Gateway interface {
	objects.Object
	Check(map[string]interface{}) ([]messages.Message, error)
	SwitchOn(map[string]interface{}) ([]messages.Message, error)
	SwitchOff(map[string]interface{}) ([]messages.Message, error)
	SwitchOnDisplayBacklight(map[string]interface{}) ([]messages.Message, error)
	SwitchOffDisplayBacklight(map[string]interface{}) ([]messages.Message, error)
	SwitchOnSilentMode(map[string]interface{}) ([]messages.Message, error)
	SwitchOffSilentMode(map[string]interface{}) ([]messages.Message, error)
	SwitchOnEcoMode(map[string]interface{}) ([]messages.Message, error)
	SwitchOffEcoMode(map[string]interface{}) ([]messages.Message, error)
	SwitchOnTurboMode(map[string]interface{}) ([]messages.Message, error)
	SwitchOffTurboMode(map[string]interface{}) ([]messages.Message, error)
	SwitchOnSleepMode(map[string]interface{}) ([]messages.Message, error)
	SwitchOffSleepMode(map[string]interface{}) ([]messages.Message, error)
	SwitchOnIonization(map[string]interface{}) ([]messages.Message, error)
	SwitchOffIonization(map[string]interface{}) ([]messages.Message, error)
	SwitchOnSelfCleaning(map[string]interface{}) ([]messages.Message, error)
	SwitchOffSelfCleaning(map[string]interface{}) ([]messages.Message, error)
	EnableSounds(map[string]interface{}) ([]messages.Message, error)
	DisableSounds(map[string]interface{}) ([]messages.Message, error)
	SetOperatingMode(args map[string]interface{}) ([]messages.Message, error)
	SetTargetTemperature(args map[string]interface{}) ([]messages.Message, error)
	SetFanSpeed(args map[string]interface{}) ([]messages.Message, error)
	SetHorizontalSlatsMode(args map[string]interface{}) ([]messages.Message, error)
	SetVerticalSlatsMode(args map[string]interface{}) ([]messages.Message, error)
}
