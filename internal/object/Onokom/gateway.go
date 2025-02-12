package Onokom

import (
	"touchon-server/internal/objects"
	"touchon-server/lib/interfaces"
)

type Gateway interface {
	objects.Object
	Check(map[string]interface{}) ([]interfaces.Message, error)
	SwitchOn(map[string]interface{}) ([]interfaces.Message, error)
	SwitchOff(map[string]interface{}) ([]interfaces.Message, error)
	SwitchOnDisplayBacklight(map[string]interface{}) ([]interfaces.Message, error)
	SwitchOffDisplayBacklight(map[string]interface{}) ([]interfaces.Message, error)
	SwitchOnSilentMode(map[string]interface{}) ([]interfaces.Message, error)
	SwitchOffSilentMode(map[string]interface{}) ([]interfaces.Message, error)
	SwitchOnEcoMode(map[string]interface{}) ([]interfaces.Message, error)
	SwitchOffEcoMode(map[string]interface{}) ([]interfaces.Message, error)
	SwitchOnTurboMode(map[string]interface{}) ([]interfaces.Message, error)
	SwitchOffTurboMode(map[string]interface{}) ([]interfaces.Message, error)
	SwitchOnSleepMode(map[string]interface{}) ([]interfaces.Message, error)
	SwitchOffSleepMode(map[string]interface{}) ([]interfaces.Message, error)
	SwitchOnIonization(map[string]interface{}) ([]interfaces.Message, error)
	SwitchOffIonization(map[string]interface{}) ([]interfaces.Message, error)
	SwitchOnSelfCleaning(map[string]interface{}) ([]interfaces.Message, error)
	SwitchOffSelfCleaning(map[string]interface{}) ([]interfaces.Message, error)
	EnableSounds(map[string]interface{}) ([]interfaces.Message, error)
	DisableSounds(map[string]interface{}) ([]interfaces.Message, error)
	SetOperatingMode(args map[string]interface{}) ([]interfaces.Message, error)
	SetTargetTemperature(args map[string]interface{}) ([]interfaces.Message, error)
	SetFanSpeed(args map[string]interface{}) ([]interfaces.Message, error)
	SetHorizontalSlatsMode(args map[string]interface{}) ([]interfaces.Message, error)
	SetVerticalSlatsMode(args map[string]interface{}) ([]interfaces.Message, error)
}
