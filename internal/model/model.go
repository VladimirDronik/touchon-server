package model

type Category = string

const (
	CategoryServer       Category = "server"        // Сервер
	CategoryController   Category = "controller"    // Контроллер
	CategoryModule       Category = "module"        // Модуль расширения
	CategorySensor       Category = "sensor"        // Датчик
	CategorySensorValue  Category = "sensor_value"  // Значение датчика
	CategoryPort         Category = "port"          // Порт
	CategoryRegulator    Category = "regulator"     // Регулятор
	CategoryGenericInput Category = "generic_input" // Универсальный вход
	CategoryRelay        Category = "relay"         // Реле
	CategoryRS485        Category = "rs485"         // Шина RS485
	CategoryModbus       Category = "modbus"        // Устройства MODBUS
	CategoryConditioner  Category = "conditioner"   // Кондиционер
)
