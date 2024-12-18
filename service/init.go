package service

// Регистрируем события
import (
	_ "github.com/VladimirDronik/touchon-server/events"
	_ "github.com/VladimirDronik/touchon-server/events/item"
	_ "github.com/VladimirDronik/touchon-server/events/object/controller"
	_ "github.com/VladimirDronik/touchon-server/events/object/generic_input"
	_ "github.com/VladimirDronik/touchon-server/events/object/modbus/wb_mrm2_mini"
	_ "github.com/VladimirDronik/touchon-server/events/object/port"
	_ "github.com/VladimirDronik/touchon-server/events/object/regulator"
	_ "github.com/VladimirDronik/touchon-server/events/object/relay"
	_ "github.com/VladimirDronik/touchon-server/events/object/sensor"
	_ "github.com/VladimirDronik/touchon-server/events/script"
)
