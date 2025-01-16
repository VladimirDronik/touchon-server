package gateway

import (
	"github.com/VladimirDronik/touchon-server/event"
	"github.com/VladimirDronik/touchon-server/models"
	"github.com/VladimirDronik/touchon-server/mqtt/messages"
	"github.com/pkg/errors"
)

func init() {
	maker := func() (*event.Event, error) {
		e := &event.Event{
			Code:        "object.onokom.gateway.on_check",
			Name:        "on_check",
			Description: "Получение состояния устройства",
			Props:       event.NewProps(),
			TargetType:  messages.TargetTypeObject,
		}

		props := []*event.Prop{
			{Code: "power_status", Name: "Состояние вкл/выкл", Item: &models.Item{Type: models.DataTypeBool}},
			{Code: "display_backlight", Name: "Подсветка экрана", Item: &models.Item{Type: models.DataTypeBool}},
			{Code: "display_brightness", Name: "Яркость экрана", Item: &models.Item{Type: models.DataTypeInt}},
			{Code: "silent_mode", Name: "Режим Тихий", Item: &models.Item{Type: models.DataTypeBool}},
			{Code: "eco_mode", Name: "Режим Эко", Item: &models.Item{Type: models.DataTypeBool}},
			{Code: "turbo_mode", Name: "Режим Турбо", Item: &models.Item{Type: models.DataTypeBool}},
			{Code: "sleep_mode", Name: "Режим Сон", Item: &models.Item{Type: models.DataTypeBool}},
			{Code: "ionization", Name: "Ионизация", Item: &models.Item{Type: models.DataTypeBool}},
			{Code: "self_cleaning", Name: "Самоочистка", Item: &models.Item{Type: models.DataTypeBool}},
			{Code: "anti_fungus", Name: "Антиплесень", Item: &models.Item{Type: models.DataTypeBool}},
			{Code: "disable_display_on_power_off", Name: "Отключение экрана при отключенном устройстве", Item: &models.Item{Type: models.DataTypeBool}},
			{Code: "sounds", Name: "Звуковая индикация", Item: &models.Item{Type: models.DataTypeBool}},
			{Code: "on_duty_heating", Name: "Дежурный обогрев", Item: &models.Item{Type: models.DataTypeBool}},
			{Code: "soft_flow", Name: "Мягкий поток", Item: &models.Item{Type: models.DataTypeBool}},
			{Code: "operating_mode", Name: "Режим работы", Item: &models.Item{Type: models.DataTypeEnum, Values: map[string]string{}}},
			{Code: "internal_temperature", Name: "Температура в помещении", Item: &models.Item{Type: models.DataTypeInt}},
			{Code: "external_temperature", Name: "Температура на улице", Item: &models.Item{Type: models.DataTypeInt}},
			{Code: "target_temperature", Name: "Целевая температура", Item: &models.Item{Type: models.DataTypeInt}},
			{Code: "fan_speed", Name: "Скорость вентилятора", Item: &models.Item{Type: models.DataTypeEnum, Values: map[string]string{}}},
			{Code: "horizontal_slats_mode", Name: "Режим работы горизонтальных ламелей", Item: &models.Item{Type: models.DataTypeEnum, Values: map[string]string{}}},
			{Code: "vertical_slats_mode", Name: "Режим работы вертикальных ламелей", Item: &models.Item{Type: models.DataTypeEnum, Values: map[string]string{}}},
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
	e, err := event.MakeEvent("object.sensor.on_check", messages.TargetTypeObject, targetID, values)
	if err != nil {
		return nil, errors.Wrap(err, "NewOnCheckMessage")
	}

	m, err := e.ToMqttMessage(topic)
	if err != nil {
		return nil, errors.Wrap(err, "NewOnCheckMessage")
	}

	return m, nil
}
