package ModbusGW

import (
	"sort"
	"time"

	"github.com/pkg/errors"
	"touchon-server/internal/model"
	"touchon-server/internal/object/Modbus/ModbusDevice"
	"touchon-server/internal/objects"
	"touchon-server/internal/scripts"
	"touchon-server/lib/events/object/onokom/gateway"
	"touchon-server/lib/models"
)

// Шлюз напрямую не используем, только в виде кондиционера с поддержкой управления по modbus.
//func init() {
//	// Регистрируем все поддерживаемые модели Modbus-шлюзов Onokom.
//	for gwModelCode := range gateways {
//		_ = objects.Register(func() (objects.Object, error) {
//			return MakeModel(gwModelCode)
//		})
//	}
//}

type Prop struct {
	Code         string
	Name         string
	Type         models.DataType
	Values       map[string]string
	DefaultValue interface{}
}

var props = []*Prop{
	{"power_status", "Состояние вкл/выкл", models.DataTypeBool, nil, false},
	{"display_backlight", "Подсветка экрана", models.DataTypeBool, nil, false},
	{"display_high_brightness", "Высокая яркость экрана", models.DataTypeBool, nil, false},
	{"silent_mode", "Режим Тихий", models.DataTypeBool, nil, false},
	{"eco_mode", "Режим Эко", models.DataTypeBool, nil, false},
	{"turbo_mode", "Режим Турбо", models.DataTypeBool, nil, false},
	{"sleep_mode", "Режим Сон", models.DataTypeBool, nil, false},
	{"ionization", "Ионизация", models.DataTypeBool, nil, false},
	{"self_cleaning", "Самоочистка", models.DataTypeBool, nil, false},
	{"anti_fungus", "Антиплесень", models.DataTypeBool, nil, false},
	{"disable_display_on_power_off", "Отключать экран при выключении устройства", models.DataTypeBool, nil, false},
	{"sounds", "Звуковая индикация", models.DataTypeBool, nil, false},
	{"on_duty_heating", "Дежурный обогрев", models.DataTypeBool, nil, false},
	{"soft_flow", "Мягкий поток", models.DataTypeBool, nil, false},
	{"operating_mode", "Режим работы", models.DataTypeEnum, nil, nil},
	{"internal_temperature", "Температура в помещении", models.DataTypeInt, nil, 25},
	{"external_temperature", "Температура на улице", models.DataTypeInt, nil, 0},
	{"target_temperature", "Целевая температура", models.DataTypeInt, nil, 25},
	{"fan_speed", "Скорость вентилятора", models.DataTypeEnum, nil, nil},
	{"horizontal_slats_mode", "Режим работы горизонтальных ламелей", models.DataTypeEnum, nil, nil},
	{"vertical_slats_mode", "Режим работы вертикальных ламелей", models.DataTypeEnum, nil, nil},
}

var propsMap = map[string]bool{}

func init() {
	for _, item := range props {
		propsMap[item.Code] = true
	}
}

func isProp(propCode string) bool {
	return propsMap[propCode]
}

func MakeModel(gwModelCode string) (objects.Object, error) {
	gw, ok := gateways[gwModelCode]
	if !ok {
		return nil, errors.Wrap(errors.Errorf("unknown gw model %q", gwModelCode), "ModbusGW.MakeModel")
	}

	baseObj, err := ModbusDevice.MakeModel()
	if err != nil {
		return nil, errors.Wrap(err, "ModbusGW.MakeModel")
	}

	obj := &GatewayModel{
		modelCode: gwModelCode,
		settings:  gw,
	}
	obj.ModbusDevice = baseObj.(*ModbusDevice.ModbusDeviceImpl)

	obj.SetCategory(model.CategoryModbusGateway)
	obj.SetType("onokom/" + gwModelCode)
	obj.SetName("Modbus шлюз ONOKOM " + gw.Name)
	obj.SetTags("onokom", "modbus", "gateway", gwModelCode)

	// Создаем свойства для сохранения состояния устройства

	for _, prop := range props {
		// Если конкретная модель шлюза не поддерживает св-во, не добавляем его в список св-в объекта.
		if gw.Registers[prop.Code] == nil {
			continue
		}

		p := &objects.Prop{
			Code:        prop.Code,
			Name:        prop.Name,
			Description: "",
			Item: &models.Item{
				Type:         prop.Type,
				Values:       prop.Values,
				DefaultValue: prop.DefaultValue,
			},
			Required: objects.True(),
			Editable: objects.False(), // запрещаем редактировать
			Visible:  objects.True(),
		}

		// Задаем списки значений для св-в с типов enum
		switch p.Code {
		case "operating_mode":
			p.Values = gw.OpModes
			p.DefaultValue = getDefaultValue(gw.OpModes)
		case "fan_speed":
			p.Values = gw.FanSpeed
			p.DefaultValue = getDefaultValue(gw.FanSpeed)
		case "horizontal_slats_mode":
			p.Values = gw.HSlatsModes
			p.DefaultValue = getDefaultValue(gw.HSlatsModes)
		case "vertical_slats_mode":
			p.Values = gw.VSlatsModes
			p.DefaultValue = getDefaultValue(gw.VSlatsModes)
		}

		if err := obj.GetProps().Add(p); err != nil {
			return nil, errors.Wrap(err, "ModbusGW.MakeModel")
		}
	}

	// Добавляем свои события
	onCheck, err := gateway.NewOnCheck(0, nil)
	if err != nil {
		return nil, errors.Wrap(err, "ModbusGW.MakeModel")
	}

	onChange, err := gateway.NewOnChange(0, nil)
	if err != nil {
		return nil, errors.Wrap(err, "ModbusGW.MakeModel")
	}

	if err := obj.GetEvents().Add(onCheck, onChange); err != nil {
		return nil, errors.Wrap(err, "ModbusGW.MakeModel")
	}

	// Регистрируем методы

	type Method struct {
		DependentProp string
		Name          string
		Description   string
		Params        []*scripts.Param
		Func          objects.MethodFunc
	}

	table := []*Method{
		// Указываем DependentProp = power_status чтобы метод был зарегистрирован
		{"power_status", "check", "Получить состояние устройства", nil, obj.Check},

		{"power_status", "switch_on", "Включить", nil, obj.SwitchOn},
		{"power_status", "switch_off", "Выключить", nil, obj.SwitchOff},

		{"display_backlight", "switch_on_display_backlight", "Включить подсветку экрана", nil, obj.SwitchOnDisplayBacklight},
		{"display_backlight", "switch_off_display_backlight", "Выключить подсветку экрана", nil, obj.SwitchOffDisplayBacklight},

		{"display_high_brightness", "switch_on_display_high_brightness", "Включить высокую яркость экрана", nil, obj.SwitchOnDisplayHighBrightness},
		{"display_high_brightness", "switch_off_display_high_brightness", "Выключить высокую яркость экрана", nil, obj.SwitchOffDisplayHighBrightness},

		{"silent_mode", "switch_on_silent_mode", "Включить режим Тихий", nil, obj.SwitchOnSilentMode},
		{"silent_mode", "switch_off_silent_mode", "Выключить режим Тихий", nil, obj.SwitchOffSilentMode},

		{"eco_mode", "switch_on_eco_mode", "Включить режим Эко", nil, obj.SwitchOnEcoMode},
		{"eco_mode", "switch_off_eco_mode", "Выключить режим Эко", nil, obj.SwitchOffEcoMode},

		{"turbo_mode", "switch_on_turbo_mode", "Включить режим Турбо", nil, obj.SwitchOnTurboMode},
		{"turbo_mode", "switch_off_turbo_mode", "Выключить режим Турбо", nil, obj.SwitchOffTurboMode},

		{"sleep_mode", "switch_on_sleep_mode", "Включить режим Сон", nil, obj.SwitchOnSleepMode},
		{"sleep_mode", "switch_off_sleep_mode", "Выключить режим Сон", nil, obj.SwitchOffSleepMode},

		{"ionization", "switch_on_ionization", "Включить функцию ионизации", nil, obj.SwitchOnIonization},
		{"ionization", "switch_off_ionization", "Выключить функцию ионизации", nil, obj.SwitchOffIonization},

		{"self_cleaning", "switch_on_self_cleaning", "Включить функцию самоочистки", nil, obj.SwitchOnSelfCleaning},
		{"self_cleaning", "switch_off_self_cleaning", "Выключить функцию самоочистки", nil, obj.SwitchOffSelfCleaning},

		{"anti_fungus", "switch_on_anti_fungus", "Включить функцию антиплесень", nil, obj.SwitchOnAntiFungus},
		{"anti_fungus", "switch_off_anti_fungus", "Выключить функцию антиплесень", nil, obj.SwitchOffAntiFungus},

		{"disable_display_on_power_off", "disable_display_on_power_off", "Экран выключен при отключенном устройстве", nil, obj.DisableDisplayOnPowerOff},
		{"disable_display_on_power_off", "enable_display_on_power_off", "Экран включен при отключенном устройстве", nil, obj.EnableDisplayOnPowerOff},

		{"sounds", "enable_sounds", "Включить звуковую индикацию", nil, obj.EnableSounds},
		{"sounds", "disable_sounds", "Выключить звуковую индикацию", nil, obj.DisableSounds},

		{"on_duty_heating", "switch_on_on_duty_heating", "Включить дежурный обогрев", nil, obj.SwitchOnOnDutyHeating},
		{"on_duty_heating", "switch_off_on_duty_heating", "Выключить дежурный обогрев", nil, obj.SwitchOffOnDutyHeating},

		{"soft_flow", "switch_on_soft_flow", "Включить функцию мягкого потока", nil, obj.SwitchOnSoftFlow},
		{"soft_flow", "switch_off_soft_flow", "Выключить функцию мягкого потока", nil, obj.SwitchOffSoftFlow},

		{"operating_mode", "set_operating_mode", "Установить режим работы", []*scripts.Param{
			{
				Code:        "operating_mode",
				Name:        "Значение",
				Description: "",
				Item: &models.Item{
					Type:   models.DataTypeEnum,
					Values: gw.OpModes,
				},
			},
		}, obj.SetOperatingMode},

		{"target_temperature", "set_target_temperature", "Задать целевую температуру", []*scripts.Param{
			{
				Code:        "target_temperature",
				Name:        "Значение",
				Description: "",
				Item: &models.Item{
					Type: models.DataTypeInt,
				},
			},
		}, obj.SetTargetTemperature},

		{"fan_speed", "set_fan_speed", "Задать скорость вентилятора", []*scripts.Param{
			{
				Code:        "fan_speed",
				Name:        "Значение",
				Description: "",
				Item: &models.Item{
					Type:   models.DataTypeEnum,
					Values: gw.FanSpeed,
				},
			},
		}, obj.SetFanSpeed},

		{"horizontal_slats_mode", "set_horizontal_slats_mode", "Задать режим для горизонтальных ламелей", []*scripts.Param{
			{
				Code:        "horizontal_slats_mode",
				Name:        "Значение",
				Description: "",
				Item: &models.Item{
					Type:   models.DataTypeEnum,
					Values: gw.HSlatsModes,
				},
			},
		}, obj.SetHorizontalSlatsMode},

		{"vertical_slats_mode", "set_vertical_slats_mode", "Задать режим для вертикальных ламелей", []*scripts.Param{
			{
				Code:        "vertical_slats_mode",
				Name:        "Значение",
				Description: "",
				Item: &models.Item{
					Type:   models.DataTypeEnum,
					Values: gw.VSlatsModes,
				},
			},
		}, obj.SetVerticalSlatsMode},
	}

	for _, item := range table {
		// Если конкретная модель шлюза не поддерживает св-во, не добавляем методы для работы с ним.
		if gw.Registers[item.DependentProp] == nil {
			continue
		}

		m, err := objects.NewMethod(item.Name, item.Description, item.Params, item.Func)
		if err != nil {
			return nil, errors.Wrap(err, "ModbusGW.MakeModel")
		}

		obj.GetMethods().Add(m)
	}

	return obj, nil
}

type GatewayModel struct {
	ModbusDevice.ModbusDevice
	unitID    int
	modelCode string
	settings  *Gateway
}

func (o *GatewayModel) Start() error {
	if err := o.ModbusDevice.Start(); err != nil {
		return errors.Wrapf(err, "ModbusGW.GatewayModel.Start(%d)", o.GetID())
	}

	address, err := o.GetProps().GetIntValue("address")
	if err != nil {
		return errors.Wrapf(err, "ModbusGW.GatewayModel.Start(%d)", o.GetID())
	}
	o.unitID = address

	updateIntervalS, err := o.GetProps().GetStringValue("update_interval")
	if err != nil {
		return errors.Wrapf(err, "ModbusGW.GatewayModel.Start(%d)", o.GetID())
	}

	updateInterval, err := time.ParseDuration(updateIntervalS)
	if err != nil {
		return errors.Wrapf(err, "ModbusGW.GatewayModel.Start(%d)", o.GetID())
	}

	o.SetTimer(updateInterval, o.check)
	o.GetTimer().Start()

	return nil
}

func (o *GatewayModel) Shutdown() error {
	if err := o.ModbusDevice.Shutdown(); err != nil {
		return errors.Wrap(err, "ModbusGW.GatewayModel.Shutdown")
	}

	return nil
}

// Используется для получения первого валидного значения из списка для установки item.DefaultValue
func getDefaultValue(values map[string]string) string {
	keys := make([]string, 0, len(values))

	for k := range values {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	if len(keys) > 0 {
		return keys[0]
	}

	return ""
}
