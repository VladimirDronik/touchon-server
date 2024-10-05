package touchon_server

// Регистрируем события
import (
	_ "github.com/VladimirDronik/touchon-server/events"
	_ "github.com/VladimirDronik/touchon-server/events/item"
	_ "github.com/VladimirDronik/touchon-server/events/object/port"
	_ "github.com/VladimirDronik/touchon-server/events/object/regulator"
	_ "github.com/VladimirDronik/touchon-server/events/object/sensor"
	_ "github.com/VladimirDronik/touchon-server/events/script"
)
