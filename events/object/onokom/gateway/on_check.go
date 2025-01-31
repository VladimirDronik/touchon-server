package gateway

import (
	"github.com/VladimirDronik/touchon-server/event"
	"github.com/VladimirDronik/touchon-server/models"
	"github.com/VladimirDronik/touchon-server/mqtt/messages"
	"github.com/pkg/errors"
)

var OpModes = map[string]string{
	"1": "Нагрев",
	"2": "Охлаждение",
	"3": "Автоматический",
	"4": "Осушение",
	"5": "Вентиляция",
}

var FanSpeed = map[string]string{
	"0": "Авто",
	"1": "Тихий режим",
	"2": "Первая скорость",
	"3": "Вторая скорость",
	"4": "Третья скорость",
	"5": "Четвертая скорость",
	"6": "Пятая скорость",
}

var HSlatsModes = map[string]string{
	"0": "Остановлено",
	"1": "Качание",
	"2": "Нижнее положение",
	"3": "Второе положение",
	"4": "Третье положение",
	"5": "Четвертое положение",
	"6": "Пятое положение",
	"7": "Шестое положение",
	"8": "Седьмое положение",
}

var VSlatsModes = map[string]string{
	"0": "Остановлено",
	"1": "Качание",
	"2": "Левое положение",
	"3": "Второе положение",
	"4": "Третье положение",
	"5": "Четвертое положение",
	"6": "Пятое положение",
	"7": "Мягкий поток",
}

var props = []*event.Prop{
	{Code: "power_status", Name: "Состояние вкл/выкл", Item: &models.Item{Type: models.DataTypeBool}},
	{Code: "display_backlight", Name: "Подсветка экрана", Item: &models.Item{Type: models.DataTypeBool}},
	{Code: "display_high_brightness", Name: "Высокая яркость экрана", Item: &models.Item{Type: models.DataTypeBool}},
	{Code: "silent_mode", Name: "Режим Тихий", Item: &models.Item{Type: models.DataTypeBool}},
	{Code: "eco_mode", Name: "Режим Эко", Item: &models.Item{Type: models.DataTypeBool}},
	{Code: "turbo_mode", Name: "Режим Турбо", Item: &models.Item{Type: models.DataTypeBool}},
	{Code: "sleep_mode", Name: "Режим Сон", Item: &models.Item{Type: models.DataTypeBool}},
	{Code: "ionization", Name: "Ионизация", Item: &models.Item{Type: models.DataTypeBool}},
	{Code: "self_cleaning", Name: "Самоочистка", Item: &models.Item{Type: models.DataTypeBool}},
	{Code: "anti_fungus", Name: "Антиплесень", Item: &models.Item{Type: models.DataTypeBool}},
	{Code: "display_enabled_always", Name: "Экран включен при отключенном устройстве", Item: &models.Item{Type: models.DataTypeBool}},
	{Code: "sounds", Name: "Звуковая индикация", Item: &models.Item{Type: models.DataTypeBool}},
	{Code: "on_duty_heating", Name: "Дежурный обогрев", Item: &models.Item{Type: models.DataTypeBool}},
	{Code: "soft_flow", Name: "Мягкий поток", Item: &models.Item{Type: models.DataTypeBool}},
	{Code: "operating_mode", Name: "Режим работы", Item: &models.Item{Type: models.DataTypeEnum, Values: OpModes}},
	{Code: "internal_temperature", Name: "Температура в помещении", Item: &models.Item{Type: models.DataTypeInt}},
	{Code: "external_temperature", Name: "Температура на улице", Item: &models.Item{Type: models.DataTypeInt}},
	{Code: "target_temperature", Name: "Целевая температура", Item: &models.Item{Type: models.DataTypeInt}},
	{Code: "fan_speed", Name: "Скорость вентилятора", Item: &models.Item{Type: models.DataTypeEnum, Values: FanSpeed}},
	{Code: "horizontal_slats_mode", Name: "Режим работы горизонтальных ламелей", Item: &models.Item{Type: models.DataTypeEnum, Values: HSlatsModes}},
	{Code: "vertical_slats_mode", Name: "Режим работы вертикальных ламелей", Item: &models.Item{Type: models.DataTypeEnum, Values: VSlatsModes}},
}

func init() {
	maker := func() (*event.Event, error) {
		e := &event.Event{
			Code:        "object.onokom.gateway.on_check",
			Name:        "on_check",
			Description: "Получение состояния устройства",
			Props:       event.NewProps(),
			TargetType:  messages.TargetTypeObject,
		}

		if err := e.Props.Add(props...); err != nil {
			return nil, errors.Wrap(err, "init.maker")
		}

		return e, nil
	}

	// Для регистрации событий надо в service/init.go добавить импорт данного _пакета_!
	if err := event.Register(maker); err != nil {
		panic(err)
	}
}

func NewOnCheckMessage(topic string, targetID int, values map[string]interface{}) (messages.Message, error) {
	e, err := event.MakeEvent("object.onokom.gateway.on_check", messages.TargetTypeObject, targetID, values)
	if err != nil {
		return nil, errors.Wrap(err, "NewOnCheckMessage")
	}

	m, err := e.ToMqttMessage(topic)
	if err != nil {
		return nil, errors.Wrap(err, "NewOnCheckMessage")
	}

	return m, nil
}
