package ModbusGW

import (
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"touchon-server/internal/g"
	"touchon-server/internal/object/RS485"
	"touchon-server/internal/object/RS485Device"
	"touchon-server/internal/objects"
	"touchon-server/internal/store"
	gw "touchon-server/lib/events/object/onokom/gateway"
	"touchon-server/lib/helpers"
	"touchon-server/lib/interfaces"
	"touchon-server/lib/models"
)

var testError = errors.New("test error")

type f = func(*GatewayModel, map[string]interface{}) ([]interfaces.Message, error)

func init() {
	g.Logger = logrus.New()
}

// Три типа в сумме поддерживают все возможные св-ва шлюзов
var gwModelCodes = []string{"tcl_1_mb_b", "gr_1_mb_b", "dk_1_mb_b"}

func getProps(gwModelCode string) ([]*Prop, error) {
	gateWay, ok := gateways[gwModelCode]
	if !ok {
		return nil, errors.Wrap(errors.Errorf("unknown gw model %q", gwModelCode), "getProps")
	}

	r := make([]*Prop, 0, len(props))

	for _, prop := range props {
		// Пропускаем не поддерживаемые св-ва
		if gateWay.Registers[prop.Code] == nil {
			continue
		}

		// Копируем элемент
		p := *prop

		switch p.Code {
		case "operating_mode":
			p.Values = gateWay.OpModes
			p.DefaultValue = getDefaultValue(gateWay.OpModes)
		case "fan_speed":
			p.Values = gateWay.FanSpeed
			p.DefaultValue = getDefaultValue(gateWay.FanSpeed)
		case "horizontal_slats_mode":
			p.Values = gateWay.HSlatsModes
			p.DefaultValue = getDefaultValue(gateWay.HSlatsModes)
		case "vertical_slats_mode":
			p.Values = gateWay.VSlatsModes
			p.DefaultValue = getDefaultValue(gateWay.VSlatsModes)
		}

		r = append(r, &p)
	}

	return r, nil
}

func setUp(t *testing.T, gwModelCode string) (*GatewayModel, *RS485Device.MockRS485Device, *store.MockStore, *interfaces.MockMessagesService, *objects.Props) {
	obj, err := MakeModel(gwModelCode, true)
	require.NotNil(t, obj)
	require.NoError(t, err)

	deviceModel, ok := obj.(*GatewayModel)
	require.True(t, ok)

	mockRS485Device := new(RS485Device.MockRS485Device)
	deviceModel.RS485Device = mockRS485Device

	gateway, ok := obj.(*GatewayModel)
	require.True(t, ok)

	st := new(store.MockStore)
	store.I = st

	msgsService := new(interfaces.MockMessagesService)
	g.Msgs = msgsService

	pList := objects.NewProps()

	props, err := getProps(gwModelCode)
	require.NoError(t, err)

	for _, item := range props {
		p := &objects.Prop{
			Code:        item.Code,
			Name:        item.Name,
			Description: "",
			Item: &models.Item{
				Type:         item.Type,
				Values:       item.Values,
				RoundFloat:   false,
				DefaultValue: item.DefaultValue,
			},
		}
		require.NoError(t, pList.Add(p))
		require.NoError(t, p.SetValue(item.DefaultValue))
	}

	a := &objects.Prop{
		Code:        "address",
		Name:        "Адрес устройства",
		Description: "",
		Item: &models.Item{
			Type:         models.DataTypeInt,
			DefaultValue: 0,
		},
	}
	require.NoError(t, a.SetValue(a.DefaultValue))

	p := &objects.Prop{
		Code:        "update_interval",
		Name:        "Интервал опроса (10s, 1m etc)",
		Description: "Интервал опроса устройства",
		Item: &models.Item{
			Type:         models.DataTypeString,
			DefaultValue: "1s",
		},
	}
	require.NoError(t, p.SetValue(p.DefaultValue))

	require.NoError(t, pList.Add(a, p))

	return gateway, mockRS485Device, st, msgsService, pList
}

func TestDeviceModel_Check(t *testing.T) {
	for _, gwModelCode := range gwModelCodes {
		gateway, rs485Device, st, msgs, pList := setUp(t, gwModelCode)

		payload := make(map[string]interface{}, pList.Len())
		objectID := 123

		for _, item := range pList.GetAll().GetValueList() {
			if !isProp(item.Code) {
				continue
			}

			payload[item.Code] = item.DefaultValue
		}

		msg, err := gw.NewOnCheck(objectID, payload)
		require.NoError(t, err)

		t.Run(gwModelCode+"/success", func(t *testing.T) {
			rs485Device.EXPECT().GetID().Return(objectID)
			rs485Device.EXPECT().GetProps().Return(pList)
			msgs.EXPECT().Send(msg).Return(nil).Once()

			_, err = gateway.Check(nil)
			require.NoError(t, err)
			time.Sleep(500 * time.Millisecond)

			st.AssertExpectations(t)
			rs485Device.AssertExpectations(t)
			msgs.AssertExpectations(t)
		})

		t.Run(gwModelCode+"/fail", func(t *testing.T) {
			rs485Device.EXPECT().GetID().Return(objectID)
			rs485Device.EXPECT().GetProps().Return(pList)
			msgs.EXPECT().Send(msg).Return(testError).Once()

			_, err = gateway.Check(nil)
			require.ErrorIs(t, err, testError)
			time.Sleep(500 * time.Millisecond)

			st.AssertExpectations(t)
			rs485Device.AssertExpectations(t)
			msgs.AssertExpectations(t)
		})
	}
}

func TestDeviceModel_check(t *testing.T) {
	objectID := 123

	for _, gwModelCode := range gwModelCodes {
		t.Run(gwModelCode+"/success", func(t *testing.T) {
			gateway, rs485Device, st, msgs, pList := setUp(t, gwModelCode)
			objRepo := new(store.MockObjectRepository)

			enabled := 2
			rs485Device.EXPECT().GetEnabled().RunAndReturn(func() bool {
				// Отключаем повторные запуски метода check()
				enabled -= 1
				return enabled > 0
			})
			rs485Device.EXPECT().GetDefaultTries().Return(3)

			doActionResponse, expectedPayload := getDoActionResponse(t, gwModelCode, pList)

			// Обращаемся к девайсу
			rs485Device.EXPECT().DoAction(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).RunAndReturn(
				func(deviceAddr int, action RS485.Action, tries int, handler RS485.ResultHandler, priority int) error {
					handler(doActionResponse, nil)
					return nil
				})

			expectedMsg, err := gw.NewOnChange(objectID, expectedPayload)
			require.NoError(t, err)

			rs485Device.EXPECT().Start().Return(nil)
			rs485Device.EXPECT().SetTimer(mock.Anything, mock.Anything)
			rs485Device.EXPECT().GetTimer().Return(helpers.NewTimer(time.Millisecond, gateway.check))
			rs485Device.EXPECT().GetID().Return(objectID)
			// Отправляем сообщение об изменении св-ва
			msgs.EXPECT().Send(expectedMsg).Return(nil)
			rs485Device.EXPECT().GetProps().Return(pList)
			st.EXPECT().ObjectRepository().Return(objRepo)
			// Правим значение св-ва в БД
			for propCode, value := range expectedPayload {
				objRepo.EXPECT().SetProp(objectID, propCode, fmt.Sprintf("%v", value)).Return(nil)
			}

			// "Запускаем" объект
			require.NoError(t, gateway.Start())

			time.Sleep(1500 * time.Millisecond)

			rs485Device.AssertExpectations(t)
			msgs.AssertExpectations(t)
			st.AssertExpectations(t)
			objRepo.AssertExpectations(t)

			// Правим состояние св-ва в памяти
			for propCode, expectedValue := range expectedPayload {
				prop, err := pList.Get(propCode)
				require.NoError(t, err)
				require.Equal(t, expectedValue, prop.GetValue())
			}
		})

		t.Run(gwModelCode+"/fail", func(t *testing.T) {
			gateway, rs485Device, _, msgs, pList := setUp(t, gwModelCode)
			objRepo := new(store.MockObjectRepository)

			enabled := 2
			rs485Device.EXPECT().GetEnabled().RunAndReturn(func() bool {
				// Отключаем повторные запуски метода check()
				enabled -= 1
				return enabled > 0
			})
			rs485Device.EXPECT().Start().Return(nil)
			rs485Device.EXPECT().SetTimer(mock.Anything, mock.Anything)
			rs485Device.EXPECT().GetTimer().Return(helpers.NewTimer(time.Millisecond, gateway.check))
			rs485Device.EXPECT().GetProps().Return(pList)
			rs485Device.EXPECT().GetID().Return(objectID)
			rs485Device.EXPECT().GetDefaultTries().Return(3)
			// Обращаемся к девайсу
			rs485Device.EXPECT().DoAction(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).RunAndReturn(
				func(deviceAddr int, action RS485.Action, tries int, handler RS485.ResultHandler, priority int) error {
					handler(nil, testError)
					return nil
				})

			// "Запускаем" объект
			require.NoError(t, gateway.Start())

			time.Sleep(1500 * time.Millisecond)

			rs485Device.AssertExpectations(t)
			msgs.AssertNotCalled(t, "Send")
			objRepo.AssertNotCalled(t, "SetProp")

			for _, prop := range pList.GetAll().GetValueList() {
				require.Equal(t, prop.DefaultValue, prop.GetValue(), prop.Code)
			}
		})
	}
}

func getDoActionResponse(t *testing.T, gwModelCode string, pList *objects.Props) (*checkDoActionResult, map[string]interface{}) {
	expectedPayload := map[string]interface{}{}

	item := gateways[gwModelCode]
	require.NotNil(t, item)

	// Получаем максимальные адреса регистров для текущего шлюза
	regs := item.Registers

	// Выделяем память сколько нужно
	doActionResponse := &checkDoActionResult{
		Coils:    make([]bool, item.getMaxCoilAddress()),
		Holdings: make([]uint16, item.getMaxHoldAddress()),
	}

	n := 1010
	for _, p := range pList.GetAll().GetValueList() {
		if !isProp(p.Code) {
			continue
		}

		reg := regs[p.Code]
		require.NotNil(t, reg)

		switch p.Type {
		case models.DataTypeBool:
			// Заполняем логические поля инвертированными значениями по умолчанию
			v := !p.DefaultValue.(bool)
			doActionResponse.Coils[reg.Address-1] = v

			// Заполняем значение для проверки
			expectedPayload[p.Code] = v

		case models.DataTypeInt:
			doActionResponse.Holdings[reg.Address-1] = uint16(n)

			// Заполняем значение для проверки
			expectedPayload[p.Code] = n / 100

			n += 1010

		case models.DataTypeEnum:
			s := getNotDefaultValue(p.Values, p.DefaultValue.(string))
			i, err := strconv.Atoi(s)
			require.NoError(t, err)

			doActionResponse.Holdings[reg.Address-1] = uint16(i)

			// Заполняем значение для проверки
			expectedPayload[p.Code] = s
		}
	}

	return doActionResponse, expectedPayload
}

// TestDeviceModel_BoolProps проверяет работу функций, включающих и выключающих логические св-ва объекта.
func TestDeviceModel_BoolProps(t *testing.T) {
	var boolProps = map[string][]f{
		"power_status":                 {(*GatewayModel).SwitchOn, (*GatewayModel).SwitchOff},
		"display_backlight":            {(*GatewayModel).SwitchOnDisplayBacklight, (*GatewayModel).SwitchOffDisplayBacklight},
		"display_high_brightness":      {(*GatewayModel).SwitchOnDisplayHighBrightness, (*GatewayModel).SwitchOffDisplayHighBrightness},
		"silent_mode":                  {(*GatewayModel).SwitchOnSilentMode, (*GatewayModel).SwitchOffSilentMode},
		"eco_mode":                     {(*GatewayModel).SwitchOnEcoMode, (*GatewayModel).SwitchOffEcoMode},
		"turbo_mode":                   {(*GatewayModel).SwitchOnTurboMode, (*GatewayModel).SwitchOffTurboMode},
		"sleep_mode":                   {(*GatewayModel).SwitchOnSleepMode, (*GatewayModel).SwitchOffSleepMode},
		"ionization":                   {(*GatewayModel).SwitchOnIonization, (*GatewayModel).SwitchOffIonization},
		"self_cleaning":                {(*GatewayModel).SwitchOnSelfCleaning, (*GatewayModel).SwitchOffSelfCleaning},
		"anti_fungus":                  {(*GatewayModel).SwitchOnAntiFungus, (*GatewayModel).SwitchOffAntiFungus},
		"disable_display_on_power_off": {(*GatewayModel).DisableDisplayOnPowerOff, (*GatewayModel).EnableDisplayOnPowerOff},
		"sounds":                       {(*GatewayModel).EnableSounds, (*GatewayModel).DisableSounds},
		"on_duty_heating":              {(*GatewayModel).SwitchOnOnDutyHeating, (*GatewayModel).SwitchOffOnDutyHeating},
		"soft_flow":                    {(*GatewayModel).SwitchOnSoftFlow, (*GatewayModel).SwitchOffSoftFlow},
	}

	objectID := 123

	for _, gwModelCode := range gwModelCodes {
		for propCode, funcs := range boolProps {
			if gateways[gwModelCode].Registers[propCode] == nil {
				t.Logf("%s не поддерживает свойство %s, пропускаем его методы", gwModelCode, propCode)
				continue
			}

			t.Run(gwModelCode+"/"+propCode+"/on/success", func(t *testing.T) {
				// Невозможно выполнить тест в параллельном режиме,
				// в setUp() строка "g.Msgs = msgs" в каждом тесте перезатирает mock объект.
				// TODO избавиться от все глобальных объектов (нужна хорошая реализация)
				//t.Parallel()

				gateway, rs485Device, st, msgs, pList := setUp(t, gwModelCode)
				objRepo := new(store.MockObjectRepository)
				require.NoError(t, pList.Set(propCode, false))

				rs485Device.EXPECT().GetDefaultTries().Return(3)
				// Обращаемся к девайсу
				rs485Device.EXPECT().DoAction(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).RunAndReturn(
					func(deviceAddr int, action RS485.Action, tries int, handler RS485.ResultHandler, priority int) error {
						handler(nil, nil)
						return nil
					})

				rs485Device.EXPECT().GetID().Return(objectID)
				// Отправляем сообщение об изменении св-ва
				msgs.EXPECT().Send(mock.Anything).Return(nil)
				rs485Device.EXPECT().GetProps().Return(pList)
				st.EXPECT().ObjectRepository().Return(objRepo)
				// Правим значение св-ва в БД
				objRepo.EXPECT().SetProp(objectID, propCode, "true").Return(nil)

				_, err := funcs[0](gateway, nil)
				require.NoError(t, err)
				time.Sleep(500 * time.Millisecond)

				rs485Device.AssertExpectations(t)
				msgs.AssertExpectations(t)
				st.AssertExpectations(t)
				objRepo.AssertExpectations(t)

				// Правим состояние св-ва в памяти
				powerStatus, err := pList.GetBoolValue(propCode)
				require.NoError(t, err)
				require.True(t, powerStatus)
			})

			t.Run(gwModelCode+"/"+propCode+"/on/fail", func(t *testing.T) {
				// Невозможно выполнить тест в параллельном режиме,
				// в setUp() строка "g.Msgs = msgs" в каждом тесте перезатирает mock объект.
				// TODO избавиться от все глобальных объектов (нужна хорошая реализация)
				//t.Parallel()

				gateway, rs485Device, _, msgs, pList := setUp(t, gwModelCode)
				objRepo := new(store.MockObjectRepository)
				require.NoError(t, pList.Set(propCode, false))

				rs485Device.EXPECT().GetDefaultTries().Return(3)
				// Обращаемся к девайсу
				rs485Device.EXPECT().DoAction(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).RunAndReturn(
					func(deviceAddr int, action RS485.Action, tries int, handler RS485.ResultHandler, priority int) error {
						handler(nil, testError)
						return nil
					})

				_, err := funcs[0](gateway, nil)
				require.NoError(t, err)
				time.Sleep(500 * time.Millisecond)

				rs485Device.AssertExpectations(t)
				msgs.AssertNotCalled(t, "Send")
				objRepo.AssertNotCalled(t, "SetProp")

				powerStatus, err := pList.GetBoolValue(propCode)
				require.NoError(t, err)
				require.False(t, powerStatus)
			})

			t.Run(gwModelCode+"/"+propCode+"/off/success", func(t *testing.T) {
				// Невозможно выполнить тест в параллельном режиме,
				// в setUp() строка "g.Msgs = msgs" в каждом тесте перезатирает mock объект.
				// TODO избавиться от все глобальных объектов (нужна хорошая реализация)
				//t.Parallel()

				gateway, rs485Device, st, msgs, pList := setUp(t, gwModelCode)
				objRepo := new(store.MockObjectRepository)
				require.NoError(t, pList.Set(propCode, true))

				rs485Device.EXPECT().GetDefaultTries().Return(3)
				// Обращаемся к девайсу
				rs485Device.EXPECT().DoAction(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).RunAndReturn(
					func(deviceAddr int, action RS485.Action, tries int, handler RS485.ResultHandler, priority int) error {
						handler(nil, nil)
						return nil
					})

				rs485Device.EXPECT().GetID().Return(objectID)
				// Отправляем сообщение об изменении св-ва
				msgs.EXPECT().Send(mock.Anything).Return(nil)
				rs485Device.EXPECT().GetProps().Return(pList)
				st.EXPECT().ObjectRepository().Return(objRepo)
				// Правим значение св-ва в БД
				objRepo.EXPECT().SetProp(objectID, propCode, "false").Return(nil)

				_, err := funcs[1](gateway, nil)
				require.NoError(t, err)
				time.Sleep(500 * time.Millisecond)

				rs485Device.AssertExpectations(t)
				msgs.AssertExpectations(t)
				st.AssertExpectations(t)
				objRepo.AssertExpectations(t)

				// Правим состояние св-ва в памяти
				powerStatus, err := pList.GetBoolValue(propCode)
				require.NoError(t, err)
				require.False(t, powerStatus)
			})

			t.Run(gwModelCode+"/"+propCode+"/off/fail", func(t *testing.T) {
				// Невозможно выполнить тест в параллельном режиме,
				// в setUp() строка "g.Msgs = msgs" в каждом тесте перезатирает mock объект.
				// TODO избавиться от все глобальных объектов (нужна хорошая реализация)
				//t.Parallel()

				gateway, rs485Device, _, msgs, pList := setUp(t, gwModelCode)
				objRepo := new(store.MockObjectRepository)
				require.NoError(t, pList.Set(propCode, true))

				rs485Device.EXPECT().GetDefaultTries().Return(3)
				// Обращаемся к девайсу
				rs485Device.EXPECT().DoAction(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).RunAndReturn(
					func(deviceAddr int, action RS485.Action, tries int, handler RS485.ResultHandler, priority int) error {
						handler(nil, testError)
						return nil
					})

				_, err := funcs[1](gateway, nil)
				require.NoError(t, err)
				time.Sleep(500 * time.Millisecond)

				rs485Device.AssertExpectations(t)
				msgs.AssertNotCalled(t, "Send")
				objRepo.AssertNotCalled(t, "SetProp")

				powerStatus, err := pList.GetBoolValue(propCode)
				require.NoError(t, err)
				require.True(t, powerStatus)
			})
		}
	}
}

func TestDeviceModel_SetTargetTemperature(t *testing.T) {
	objectID := 123
	propCode := "target_temperature"
	targetTemp := 30

	t.Run("success", func(t *testing.T) {
		gateway, rs485Device, st, msgs, pList := setUp(t, gwModelCodes[0])
		objRepo := new(store.MockObjectRepository)
		require.NoError(t, pList.Set(propCode, 0))

		rs485Device.EXPECT().GetDefaultTries().Return(3)
		// Обращаемся к девайсу
		rs485Device.EXPECT().DoAction(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).RunAndReturn(
			func(deviceAddr int, action RS485.Action, tries int, handler RS485.ResultHandler, priority int) error {
				handler(nil, nil)
				return nil
			})

		rs485Device.EXPECT().GetID().Return(objectID)
		// Отправляем сообщение об изменении св-ва
		msgs.EXPECT().Send(mock.Anything).Return(nil)
		rs485Device.EXPECT().GetProps().Return(pList)
		st.EXPECT().ObjectRepository().Return(objRepo)
		// Правим значение св-ва в БД
		objRepo.EXPECT().SetProp(objectID, propCode, strconv.Itoa(targetTemp)).Return(nil)

		_, err := gateway.SetTargetTemperature(map[string]interface{}{propCode: targetTemp})
		require.NoError(t, err)
		time.Sleep(500 * time.Millisecond)

		rs485Device.AssertExpectations(t)
		msgs.AssertExpectations(t)
		st.AssertExpectations(t)
		objRepo.AssertExpectations(t)

		// Правим состояние св-ва в памяти
		value, err := pList.GetIntValue(propCode)
		require.NoError(t, err)
		require.Equal(t, targetTemp, value)
	})

	t.Run("fail", func(t *testing.T) {
		gateway, rs485Device, _, msgs, pList := setUp(t, gwModelCodes[0])
		objRepo := new(store.MockObjectRepository)
		require.NoError(t, pList.Set(propCode, 0))

		rs485Device.EXPECT().GetDefaultTries().Return(3)
		// Обращаемся к девайсу
		rs485Device.EXPECT().DoAction(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).RunAndReturn(
			func(deviceAddr int, action RS485.Action, tries int, handler RS485.ResultHandler, priority int) error {
				handler(nil, testError)
				return nil
			})

		_, err := gateway.SetTargetTemperature(map[string]interface{}{propCode: targetTemp})
		require.NoError(t, err)
		time.Sleep(500 * time.Millisecond)

		rs485Device.AssertExpectations(t)
		msgs.AssertNotCalled(t, "Send")
		objRepo.AssertNotCalled(t, "SetProp")

		value, err := pList.GetIntValue(propCode)
		require.NoError(t, err)
		require.NotEqual(t, targetTemp, value)
	})
}

func TestDeviceModel_EnumProps(t *testing.T) {
	type enumTest struct {
		F      f
		Values map[string]string
	}

	objectID := 123

	for _, gwModelCode := range gwModelCodes {
		gateWay := gateways[gwModelCode]

		var enumProps = map[string]*enumTest{
			"operating_mode":        {(*GatewayModel).SetOperatingMode, gateWay.OpModes},
			"fan_speed":             {(*GatewayModel).SetFanSpeed, gateWay.FanSpeed},
			"horizontal_slats_mode": {(*GatewayModel).SetHorizontalSlatsMode, gateWay.HSlatsModes},
			"vertical_slats_mode":   {(*GatewayModel).SetVerticalSlatsMode, gateWay.VSlatsModes},
		}

		for propCode, test := range enumProps {
			if gateways[gwModelCode].Registers[propCode] == nil {
				t.Logf("%s не поддерживает свойство %s, пропускаем его методы", gwModelCode, propCode)
				continue
			}

			for v, descr := range test.Values {
				t.Run(fmt.Sprintf("%s/%s/success (%s, %s)", gwModelCode, propCode, v, descr), func(t *testing.T) {
					gateway, rs485Device, st, msgs, pList := setUp(t, gwModelCode)
					objRepo := new(store.MockObjectRepository)

					rs485Device.EXPECT().GetDefaultTries().Return(3)
					// Обращаемся к девайсу
					rs485Device.EXPECT().DoAction(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).RunAndReturn(
						func(deviceAddr int, action RS485.Action, tries int, handler RS485.ResultHandler, priority int) error {
							handler(nil, nil)
							return nil
						})

					rs485Device.EXPECT().GetID().Return(objectID)
					// Отправляем сообщение об изменении св-ва
					msgs.EXPECT().Send(mock.Anything).Return(nil)
					rs485Device.EXPECT().GetProps().Return(pList)
					st.EXPECT().ObjectRepository().Return(objRepo)
					// Правим значение св-ва в БД
					objRepo.EXPECT().SetProp(objectID, propCode, v).Return(nil)

					_, err := test.F(gateway, map[string]interface{}{propCode: v})
					require.NoError(t, err)
					time.Sleep(500 * time.Millisecond)

					rs485Device.AssertExpectations(t)
					msgs.AssertExpectations(t)
					st.AssertExpectations(t)
					objRepo.AssertExpectations(t)

					// Правим состояние св-ва в памяти
					value, err := pList.GetEnumValue(propCode)
					require.NoError(t, err)
					require.Equal(t, v, value)
				})
			}

			t.Run(fmt.Sprintf("%s/%s/fail", gwModelCode, propCode), func(t *testing.T) {
				gateway, rs485Device, _, msgs, pList := setUp(t, gwModelCode)
				objRepo := new(store.MockObjectRepository)
				v := "123"

				rs485Device.EXPECT().GetDefaultTries().Return(3)
				// Обращаемся к девайсу
				rs485Device.EXPECT().DoAction(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).RunAndReturn(
					func(deviceAddr int, action RS485.Action, tries int, handler RS485.ResultHandler, priority int) error {
						handler(nil, testError)
						return nil
					})

				_, err := test.F(gateway, map[string]interface{}{propCode: v})
				require.NoError(t, err)
				time.Sleep(500 * time.Millisecond)

				rs485Device.AssertExpectations(t)
				msgs.AssertNotCalled(t, "Send")
				objRepo.AssertNotCalled(t, "SetProp")

				value, err := pList.GetEnumValue(propCode)
				require.NoError(t, err)
				require.NotEqual(t, v, value)
			})
		}
	}
}

func getNotDefaultValue(values map[string]string, defaultValue string) string {
	for k := range values {
		if k != defaultValue {
			return k
		}
	}

	return ""
}
