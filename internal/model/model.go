package model

type Category string

const (
	CategoryController    Category = "controller"     // Контроллер
	CategoryModule        Category = "module"         // Модуль расширения
	CategorySensor        Category = "sensor"         // Датчик
	CategorySensorValue   Category = "sensor_value"   // Значение датчика
	CategoryPort          Category = "port"           // Порт
	CategoryRegulator     Category = "regulator"      // Регулятор
	CategoryGenericInput  Category = "generic_input"  // Универсальный вход
	CategoryRelay         Category = "relay"          // Реле
	CategoryModbus        Category = "modbus"         // MODBUS шина
	CategoryModbusGateway Category = "modbus_gateway" // Шлюз Modbus
	CategoryConditioner   Category = "conditioner"    // Кондиционер
)
