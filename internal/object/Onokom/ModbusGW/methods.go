package ModbusGW

import (
	"fmt"
	"strconv"

	"github.com/pkg/errors"
	"github.com/simonvetter/modbus"
	"touchon-server/internal/g"
	"touchon-server/internal/object/Modbus"
	"touchon-server/internal/store"
	"touchon-server/lib/events/object/onokom/gateway"
	"touchon-server/lib/interfaces"
)

var doActionErr = errors.Errorf("ModbusDeviceImpl.DoAction returned bad value")

// Check возвращает значения всех свойств устройства
func (o *GatewayModel) Check(map[string]interface{}) ([]interfaces.Message, error) {
	payload := make(map[string]interface{}, o.GetProps().Len())

	for _, prop := range o.GetProps().GetAll().GetValueList() {
		if !isProp(prop.Code) {
			continue
		}

		payload[prop.Code] = prop.GetValue()
	}

	msg, err := gateway.NewOnCheck(o.GetID(), payload)
	if err != nil {
		return nil, errors.Wrap(err, "Check")
	}

	if err := g.Msgs.Send(msg); err != nil {
		return nil, errors.Wrap(err, "Check")
	}

	return nil, nil
}

type checkDoActionResult struct {
	Coils    []bool
	Holdings []uint16
}

// check опрашивает устройство, обновляет данные в БД и памяти, оповещает об изменениях
func (o *GatewayModel) check() {
	if !o.GetEnabled() {
		return
	}

	g.Logger.Debugf("ModbusGW.GatewayModel(%d): check()", o.GetID())

	maxCoilAddr := o.settings.getMaxCoilAddress()
	maxHoldAddr := o.settings.getMaxHoldAddress()

	action := func(client Modbus.Client) (interface{}, error) {
		r := &checkDoActionResult{}
		var err error

		// Вычитываем все логические значения разом
		r.Coils, err = client.ReadCoils(0x0001, maxCoilAddr)
		if err != nil {
			return nil, err
		}

		// Вычитываем все holding регистры разом
		r.Holdings, err = client.ReadRegisters(0x0001, maxHoldAddr, modbus.HOLDING_REGISTER)
		if err != nil {
			return nil, err
		}

		return r, nil
	}

	resultHandler := func(r interface{}, err error) {
		// Перезапускаем таймер только после полного завершения обновления св-в
		defer o.GetTimer().Reset()

		if err != nil {
			g.Logger.Error(errors.Wrap(err, "ModbusGW.GatewayModel.check"))
			return
		}

		res, ok := r.(*checkDoActionResult)
		if !ok {
			g.Logger.Error(errors.Wrap(doActionErr, "ModbusGW.GatewayModel.check"))
			return
		}

		switch {
		case len(res.Coils) != int(maxCoilAddr) || len(res.Holdings) != int(maxHoldAddr):
			g.Logger.Error(errors.Wrap(doActionErr, "ModbusGW.GatewayModel.check"))
			return
		}

		payload := make(map[string]interface{}, len(props))

		// Обрабатываем логические св-ва
		for propCode, reg := range o.settings.Registers {
			if reg == nil || reg.Type != Coil {
				continue
			}

			status := res.Coils[reg.Address-1]

			prop, err := o.GetProps().Get(propCode)
			if err != nil {
				g.Logger.Error(errors.Wrap(err, "ModbusGW.GatewayModel.check"))
				return
			}

			actualValue, err := prop.GetBoolValue()
			if err != nil {
				g.Logger.Error(errors.Wrap(err, "ModbusGW.GatewayModel.check"))
				return
			}

			// Если значение не изменилось, пропускаем
			if status == actualValue {
				continue
			}

			// Обновляем значение св-ва в объекте
			if err := prop.SetValue(status); err != nil {
				g.Logger.Error(errors.Wrap(err, "ModbusGW.GatewayModel.check"))
				return
			}

			// Обновляем значение св-ва в базе
			if err := store.I.ObjectRepository().SetProp(o.GetID(), propCode, strconv.FormatBool(status)); err != nil {
				g.Logger.Error(errors.Wrap(err, "ModbusGW.GatewayModel.check"))
				return
			}

			// Добавляем в список измененных полей
			payload[propCode] = status
		}

		// Обрабатываем holding регистры
		for propCode, reg := range o.settings.Registers {
			if reg == nil || reg.Type != Hold {
				continue
			}

			value := int(res.Holdings[reg.Address-1])

			if propCode == "internal_temperature" || propCode == "external_temperature" || propCode == "target_temperature" {
				value /= 100
			}

			prop, err := o.GetProps().Get(propCode)
			if err != nil {
				g.Logger.Error(errors.Wrap(err, "ModbusGW.GatewayModel.check"))
				return
			}

			// Если значение не изменилось, пропускаем
			if value == prop.GetValue() {
				continue
			}

			// Обновляем значение св-ва в объекте
			if err := prop.SetValue(value); err != nil {
				g.Logger.Error(errors.Wrap(err, "ModbusGW.GatewayModel.check"))

				// Если значение не валидное, пропускаем его и переходим к следующему св-ву.
				// Значение может быть не валидным (например, равным 0), если кондер не подключен.
				// При этом шлюз может возвращать валидные значения из других регистров,
				// сохраненные в памяти (это догадка). Например, hr_1_mb_b возвращает 0 для
				// vertical_slats_mode - это не валидное значение. При этом для operating_mode
				// возвращает 3, а для target_temperature - 2400 - оба валидные значения.
				continue
			}

			// Обновляем значение св-ва в базе
			if err := store.I.ObjectRepository().SetProp(o.GetID(), propCode, fmt.Sprintf("%v", value)); err != nil {
				g.Logger.Error(errors.Wrap(err, "ModbusGW.GatewayModel.check"))
				return
			}

			// Добавляем в список измененных полей
			payload[propCode] = prop.GetValue() // Возвращаем значение с нужным типом (для enum типов будет строка)
		}

		// Отправляем сообщение с измененными полями
		msg, err := gateway.NewOnChange(o.GetID(), payload)
		if err != nil {
			g.Logger.Error(errors.Wrap(err, "ModbusGW.GatewayModel.check"))
			return
		}

		// Отправляем сообщение об изменении св-ва объекта
		if err := g.Msgs.Send(msg); err != nil {
			g.Logger.Error(errors.Wrap(err, "ModbusGW.GatewayModel.check"))
			return
		}
	}

	if err := o.ModbusDevice.DoAction(o.unitID, action, o.GetDefaultTries(), resultHandler, Modbus.QueueMinPriority); err != nil {
		g.Logger.Error(errors.Wrap(err, "ModbusGW.GatewayModel.check"))
		return
	}
}

// Состояние (on/off)

func (o *GatewayModel) SwitchOn(map[string]interface{}) ([]interfaces.Message, error) {
	return o.setCoilWrapper("power_status", true, "SwitchOn")
}

func (o *GatewayModel) SwitchOff(map[string]interface{}) ([]interfaces.Message, error) {
	return o.setCoilWrapper("power_status", false, "SwitchOff")
}

// Подсветка экрана (on/off)

func (o *GatewayModel) SwitchOnDisplayBacklight(map[string]interface{}) ([]interfaces.Message, error) {
	return o.setCoilWrapper("display_backlight", true, "SwitchOnDisplayBacklight")
}

func (o *GatewayModel) SwitchOffDisplayBacklight(map[string]interface{}) ([]interfaces.Message, error) {
	return o.setCoilWrapper("display_backlight", false, "SwitchOffDisplayBacklight")
}

// Яркость экрана

func (o *GatewayModel) SwitchOnDisplayHighBrightness(map[string]interface{}) ([]interfaces.Message, error) {
	return o.setCoilWrapper("display_high_brightness", true, "SwitchOnDisplayHighBrightness")
}

func (o *GatewayModel) SwitchOffDisplayHighBrightness(map[string]interface{}) ([]interfaces.Message, error) {
	return o.setCoilWrapper("display_high_brightness", false, "SwitchOffDisplayHighBrightness")
}

// Режим Тихий (on/off)

func (o *GatewayModel) SwitchOnSilentMode(map[string]interface{}) ([]interfaces.Message, error) {
	return o.setCoilWrapper("silent_mode", true, "SwitchOnSilentMode")
}

func (o *GatewayModel) SwitchOffSilentMode(map[string]interface{}) ([]interfaces.Message, error) {
	return o.setCoilWrapper("silent_mode", false, "SwitchOffSilentMode")
}

// Режим Эко (on/off)

func (o *GatewayModel) SwitchOnEcoMode(map[string]interface{}) ([]interfaces.Message, error) {
	return o.setCoilWrapper("eco_mode", true, "SwitchOnEcoMode")
}

func (o *GatewayModel) SwitchOffEcoMode(map[string]interface{}) ([]interfaces.Message, error) {
	return o.setCoilWrapper("eco_mode", false, "SwitchOffEcoMode")
}

// Режим Турбо (on/off)

func (o *GatewayModel) SwitchOnTurboMode(map[string]interface{}) ([]interfaces.Message, error) {
	return o.setCoilWrapper("turbo_mode", true, "SwitchOnTurboMode")
}

func (o *GatewayModel) SwitchOffTurboMode(map[string]interface{}) ([]interfaces.Message, error) {
	return o.setCoilWrapper("turbo_mode", false, "SwitchOffTurboMode")
}

// Режим Сон (on/off)

func (o *GatewayModel) SwitchOnSleepMode(map[string]interface{}) ([]interfaces.Message, error) {
	return o.setCoilWrapper("sleep_mode", true, "SwitchOnSleepMode")
}

func (o *GatewayModel) SwitchOffSleepMode(map[string]interface{}) ([]interfaces.Message, error) {
	return o.setCoilWrapper("sleep_mode", false, "SwitchOffSleepMode")
}

// Функция ионизации (on/off)

func (o *GatewayModel) SwitchOnIonization(map[string]interface{}) ([]interfaces.Message, error) {
	return o.setCoilWrapper("ionization", true, "SwitchOnIonization")
}

func (o *GatewayModel) SwitchOffIonization(map[string]interface{}) ([]interfaces.Message, error) {
	return o.setCoilWrapper("ionization", false, "SwitchOffIonization")
}

// Функция самоочистки (on/off)

func (o *GatewayModel) SwitchOnSelfCleaning(map[string]interface{}) ([]interfaces.Message, error) {
	return o.setCoilWrapper("self_cleaning", true, "SwitchOnSelfCleaning")
}

func (o *GatewayModel) SwitchOffSelfCleaning(map[string]interface{}) ([]interfaces.Message, error) {
	return o.setCoilWrapper("self_cleaning", false, "SwitchOffSelfCleaning")
}

// Функция антиплесень

func (o *GatewayModel) SwitchOnAntiFungus(map[string]interface{}) ([]interfaces.Message, error) {
	return o.setCoilWrapper("anti_fungus", true, "SwitchOnAntiFungus")
}

func (o *GatewayModel) SwitchOffAntiFungus(map[string]interface{}) ([]interfaces.Message, error) {
	return o.setCoilWrapper("anti_fungus", false, "SwitchOffAntiFungus")
}

// Отключение экрана при отключенном кондиционере

func (o *GatewayModel) DisableDisplayOnPowerOff(map[string]interface{}) ([]interfaces.Message, error) {
	// По умолчанию при отключении кондиционера экран отображает красный символ питания,
	// при активации этой настройки при отключенном кондиционере экран будет отключаться.
	// 0 - Отключена 1 - Включена

	return o.setCoilWrapper("disable_display_on_power_off", true, "DisableDisplayOnPowerOff")
}

func (o *GatewayModel) EnableDisplayOnPowerOff(map[string]interface{}) ([]interfaces.Message, error) {
	return o.setCoilWrapper("disable_display_on_power_off", false, "EnableDisplayOnPowerOff")
}

// Звуковая индикация (on/off)

func (o *GatewayModel) EnableSounds(map[string]interface{}) ([]interfaces.Message, error) {
	return o.setCoilWrapper("sounds", true, "EnableSounds")
}

func (o *GatewayModel) DisableSounds(map[string]interface{}) ([]interfaces.Message, error) {
	return o.setCoilWrapper("sounds", false, "DisableSounds")
}

// Дежурный обогрев

func (o *GatewayModel) SwitchOnOnDutyHeating(map[string]interface{}) ([]interfaces.Message, error) {
	return o.setCoilWrapper("on_duty_heating", true, "SwitchOnOnDutyHeating")
}

func (o *GatewayModel) SwitchOffOnDutyHeating(map[string]interface{}) ([]interfaces.Message, error) {
	return o.setCoilWrapper("on_duty_heating", false, "SwitchOffOnDutyHeating")
}

// Мягкий поток

func (o *GatewayModel) SwitchOnSoftFlow(map[string]interface{}) ([]interfaces.Message, error) {
	return o.setCoilWrapper("soft_flow", true, "SwitchOnSoftFlow")
}

func (o *GatewayModel) SwitchOffSoftFlow(map[string]interface{}) ([]interfaces.Message, error) {
	return o.setCoilWrapper("soft_flow", false, "SwitchOffSoftFlow")
}

// Режим работы (Нагрев,Охлаждение,Автоматический,Осушение,Вентиляция)

func (o *GatewayModel) SetOperatingMode(args map[string]interface{}) ([]interfaces.Message, error) {
	return o.setHoldingEnumValue("operating_mode", args["operating_mode"], "SetOperatingMode")
}

// Целевая t° (есть/нет)

func (o *GatewayModel) SetTargetTemperature(args map[string]interface{}) ([]interfaces.Message, error) { // target_temperature
	propCode := "target_temperature"

	v, ok := args[propCode]
	if !ok {
		return nil, errors.Wrap(errors.New("value not found"), "SetTargetTemperature")
	}

	var t int
	switch v := v.(type) {
	case int:
		t = v
	case float64:
		t = int(v)
	default:
		return nil, errors.Wrap(errors.Errorf("value is not int/float64: %T", v), "SetTargetTemperature")
	}

	action := func(client Modbus.Client) (interface{}, error) {
		// Выставляем состояние устройства
		return nil, client.WriteRegister(0x0005, uint16(t*100))
	}

	resultHandler := func(r interface{}, err error) {
		if err != nil {
			g.Logger.Error(errors.Wrap(err, "SetTargetTemperature"))
			return
		}

		msg, err := gateway.NewOnChange(o.GetID(), map[string]interface{}{propCode: t})
		if err != nil {
			g.Logger.Error(errors.Wrap(err, "SetTargetTemperature"))
			return
		}

		// Отправляем сообщение об изменении св-ва объекта
		if err := g.Msgs.Send(msg); err != nil {
			g.Logger.Error(errors.Wrap(err, "SetTargetTemperature"))
			return
		}

		// Обновляем значение св-ва в объекте
		if err := o.GetProps().Set(propCode, t); err != nil {
			g.Logger.Error(errors.Wrap(err, "SetTargetTemperature"))
			return
		}

		// Обновляем значение св-ва в базе
		if err := store.I.ObjectRepository().SetProp(o.GetID(), propCode, strconv.Itoa(t)); err != nil {
			g.Logger.Error(errors.Wrap(err, "SetTargetTemperature"))
			return
		}
	}

	if err := o.ModbusDevice.DoAction(o.unitID, action, o.GetDefaultTries(), resultHandler, Modbus.QueueMinPriority); err != nil {
		return nil, errors.Wrap(err, "SetTargetTemperature")
	}

	return nil, nil
}

// Скорость вентилятора (Авто,Первая скорость,Вторая скорость,Третья скорость)

func (o *GatewayModel) SetFanSpeed(args map[string]interface{}) ([]interfaces.Message, error) {
	return o.setHoldingEnumValue("fan_speed", args["fan_speed"], "SetFanSpeed")
}

// Горизонтальные ламели (Остановлено,Качание,1 Положение (нижнее),2 Положение,3 Положение,4 Положение,5 Положение,6 Положение,7 Положение)

func (o *GatewayModel) SetHorizontalSlatsMode(args map[string]interface{}) ([]interfaces.Message, error) {
	return o.setHoldingEnumValue("horizontal_slats_mode", args["horizontal_slats_mode"], "SetHorizontalSlatsMode")
}

// Вертикальные ламели (Остановлено,Качание,1 Положение (левое),2 Положение,3 Положение,4 Положение,5 Положение)

func (o *GatewayModel) SetVerticalSlatsMode(args map[string]interface{}) ([]interfaces.Message, error) {
	return o.setHoldingEnumValue("vertical_slats_mode", args["vertical_slats_mode"], "SetVerticalSlatsMode")
}

func (o *GatewayModel) setCoilWrapper(propCode string, status bool, errFuncName string) ([]interfaces.Message, error) {
	reg := o.settings.Registers[propCode]
	if reg == nil {
		return nil, errors.Wrap(errors.Wrap(errors.Errorf("шлюз %s не поддерживает свойство %s", o.modelCode, propCode), "setCoilWrapper"), errFuncName)
	}

	action := func(client Modbus.Client) (interface{}, error) {
		// Выставляем состояние устройства
		return nil, client.WriteCoil(reg.Address, status)
	}

	resultHandler := func(r interface{}, err error) {
		if err != nil {
			g.Logger.Error(errors.Wrap(errors.Wrap(err, "setCoilWrapper"), errFuncName))
			return
		}

		msg, err := gateway.NewOnChange(o.GetID(), map[string]interface{}{propCode: status})
		if err != nil {
			g.Logger.Error(errors.Wrap(errors.Wrap(err, "setCoilWrapper"), errFuncName))
			return
		}

		// Отправляем сообщение об изменении св-ва объекта
		if err := g.Msgs.Send(msg); err != nil {
			g.Logger.Error(errors.Wrap(errors.Wrap(err, "setCoilWrapper"), errFuncName))
			return
		}

		// Обновляем значение св-ва в объекте
		if err := o.GetProps().Set(propCode, status); err != nil {
			g.Logger.Error(errors.Wrap(errors.Wrap(err, "setCoilWrapper"), errFuncName))
			return
		}

		// Обновляем значение св-ва в базе
		if err := store.I.ObjectRepository().SetProp(o.GetID(), propCode, strconv.FormatBool(status)); err != nil {
			g.Logger.Error(errors.Wrap(errors.Wrap(err, "setCoilWrapper"), errFuncName))
			return
		}
	}

	if err := o.ModbusDevice.DoAction(o.unitID, action, o.GetDefaultTries(), resultHandler, Modbus.QueueMinPriority); err != nil {
		return nil, errors.Wrap(errors.Wrap(err, "setCoilWrapper"), errFuncName)
	}

	return nil, nil
}

func (o *GatewayModel) setHoldingEnumValue(propCode string, value interface{}, errFuncName string) ([]interfaces.Message, error) {
	reg := o.settings.Registers[propCode]
	if reg == nil {
		return nil, errors.Wrap(errors.Wrap(errors.Errorf("шлюз %s не поддерживает свойство %s", o.modelCode, propCode), "setHoldingEnumValue"), errFuncName)
	}

	if value == nil {
		return nil, errors.Wrap(errors.Wrap(errors.New("value not found"), "setHoldingEnumValue"), errFuncName)
	}

	v, ok := value.(string)
	if !ok {
		return nil, errors.Wrap(errors.Wrap(errors.Errorf("value is not string: %T", value), "setHoldingEnumValue"), errFuncName)
	}

	i, err := strconv.Atoi(v)
	if err != nil {
		return nil, errors.Wrap(errors.Wrap(err, "setHoldingEnumValue"), errFuncName)
	}

	action := func(client Modbus.Client) (interface{}, error) {
		return nil, client.WriteRegister(reg.Address, uint16(i))
	}

	resultHandler := func(r interface{}, err error) {
		if err != nil {
			g.Logger.Error(errors.Wrap(errors.Wrap(err, "setHoldingEnumValue"), errFuncName))
			return
		}

		msg, err := gateway.NewOnChange(o.GetID(), map[string]interface{}{propCode: v})
		if err != nil {
			g.Logger.Error(errors.Wrap(errors.Wrap(err, "setHoldingEnumValue"), errFuncName))
			return
		}

		// Отправляем сообщение об изменении св-ва объекта
		if err := g.Msgs.Send(msg); err != nil {
			g.Logger.Error(errors.Wrap(errors.Wrap(err, "setHoldingEnumValue"), errFuncName))
			return
		}

		// Обновляем значение св-ва в объекте
		if err := o.GetProps().Set(propCode, v); err != nil {
			g.Logger.Error(errors.Wrap(errors.Wrap(err, "setHoldingEnumValue"), errFuncName))
			return
		}

		// Обновляем значение св-ва в базе
		if err := store.I.ObjectRepository().SetProp(o.GetID(), propCode, v); err != nil {
			g.Logger.Error(errors.Wrap(errors.Wrap(err, "setHoldingEnumValue"), errFuncName))
			return
		}
	}

	if err := o.ModbusDevice.DoAction(o.unitID, action, o.GetDefaultTries(), resultHandler, Modbus.QueueMinPriority); err != nil {
		return nil, errors.Wrap(errors.Wrap(err, "setHoldingEnumValue"), errFuncName)
	}

	return nil, nil
}
