package service

// Регистрируем события
import (
	_ "touchon-server/lib/events"
	_ "touchon-server/lib/events/item"
	_ "touchon-server/lib/events/object/controller"
	_ "touchon-server/lib/events/object/generic_input"
	_ "touchon-server/lib/events/object/onokom/gateway"
	_ "touchon-server/lib/events/object/port"
	_ "touchon-server/lib/events/object/regulator"
	_ "touchon-server/lib/events/object/relay"
	_ "touchon-server/lib/events/object/sensor"
	_ "touchon-server/lib/events/object/wiren_board/wb_mrm2_mini"
	_ "touchon-server/lib/events/script"
)
