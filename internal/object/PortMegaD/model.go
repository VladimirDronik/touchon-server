// https://ab-log.ru/smart-house/ethernet/megad-2561

package PortMegaD

import (
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
	"touchon-server/internal/model"
	"touchon-server/internal/objects"
	"touchon-server/internal/scripts"
	"touchon-server/internal/store"
	"touchon-server/lib/events"
	"touchon-server/lib/events/object/port"
	"touchon-server/lib/interfaces"
	"touchon-server/lib/models"
)

func init() {
	_ = objects.Register(MakeModel)
}

func MakeModel() (objects.Object, error) {
	props := []*objects.Prop{
		{
			Code:        "number",
			Name:        "Номер порта",
			Description: "",
			Item: &models.Item{
				Type:         models.DataTypeInt,
				DefaultValue: 0,
			},
			Required: objects.True(),
			Editable: objects.False(),
			Visible:  objects.True(),
			CheckValue: func(p *objects.Prop, allProps map[string]*objects.Prop) error {
				if v, _ := p.GetIntValue(); v < 0 {
					return errors.New("Port < 0")
				}

				return nil
			},
		},
		{
			Code:        "group",
			Name:        "Группа",
			Description: "",
			Item: &models.Item{
				Type: models.DataTypeInterface,
				// Группа задается при добавлении порта в контроллере
			},
			Required: objects.True(),
			Editable: objects.False(),
			Visible:  objects.True(),
		}, {
			Code:        "type",
			Name:        "Тип",
			Description: "",
			Item: &models.Item{
				Type: models.DataTypeInterface,
				// Типы задаются при добавлении порта в контроллере
			},
			Required: objects.True(),
			Editable: objects.True(),
			Visible:  objects.True(),
		},
		{
			Code:        "mode",
			Name:        "Режим",
			Description: "",
			Item: &models.Item{
				Type: models.DataTypeInterface,
			},
			Required: objects.True(),
			Editable: objects.True(),
			Visible:  objects.True(),
		},
		{
			Code:        "parent_port",
			Name:        "Порт родительского объекта",
			Description: "",
			Item: &models.Item{
				Type:         models.DataTypeInt,
				DefaultValue: 0,
			},
			Required: objects.False(),
			Editable: objects.True(),
			Visible:  objects.True(),
			CheckValue: func(p *objects.Prop, allProps map[string]*objects.Prop) error {
				if v, _ := p.GetIntValue(); v < 0 {
					return errors.New("ParentPort < 0")
				}

				return nil
			},
		},
	}

	onPress, err := port.NewOnPress(0)
	if err != nil {
		return nil, errors.Wrap(err, "PortMegaD.MakeModel")
	}

	onRelease, err := port.NewOnRelease(0)
	if err != nil {
		return nil, errors.Wrap(err, "PortMegaD.MakeModel")
	}

	onDoubleClick, err := port.NewOnDoubleClick(0)
	if err != nil {
		return nil, errors.Wrap(err, "PortMegaD.MakeModel")
	}

	onLongPress, err := port.NewOnLongPress(0)
	if err != nil {
		return nil, errors.Wrap(err, "PortMegaD.MakeModel")
	}

	onCheck, err := port.NewOnCheck(0, "", "")
	if err != nil {
		return nil, errors.Wrap(err, "PortMegaD.MakeModel")
	}

	onChangeState, err := events.NewOnChangeState(interfaces.TargetTypeObject, 0, "", "")
	if err != nil {
		return nil, errors.Wrap(err, "PortMegaD.MakeModel")
	}

	impl, err := objects.NewObjectModelImpl(
		model.CategoryPort,
		"port_mega_d",
		true,
		"Порт",
		props,
		nil,
		[]interfaces.Event{onPress, onRelease, onDoubleClick, onLongPress, onCheck, onChangeState},
		nil,
		[]string{"port", "port_mega_d"},
	)
	if err != nil {
		return nil, errors.Wrap(err, "PortMegaD.MakeModel")
	}

	o := &PortModel{ObjectModelImpl: impl}

	check, err := objects.NewMethod("check", "Получение статуса порта", nil, o.Check)
	if err != nil {
		return nil, errors.Wrap(err, "PortMegaD.MakeModel")
	}

	on, err := objects.NewMethod("on", "Включаем порт", nil, o.On)
	if err != nil {
		return nil, errors.Wrap(err, "PortMegaD.MakeModel")
	}

	off, err := objects.NewMethod("off", "Выключаем порт", nil, o.Off)
	if err != nil {
		return nil, errors.Wrap(err, "PortMegaD.MakeModel")
	}

	toggle, err := objects.NewMethod("toggle", "Переключаем состояние порта вкл/выкл", nil, o.Toggle)
	if err != nil {
		return nil, errors.Wrap(err, "PortMegaD.MakeModel")
	}

	params := []*scripts.Param{
		{
			Code:        "value",
			Name:        "value",
			Description: "",
			Item: &models.Item{
				Type:         models.DataTypeInt,
				DefaultValue: 0,
			},
		},
		{
			Code:        "smooth",
			Name:        "smooth",
			Description: "",
			Item: &models.Item{
				Type:         models.DataTypeInt,
				DefaultValue: 0,
			},
		},
	}

	setPWM, err := objects.NewMethod("set_pwm", "Установка ШИМ для порта", params, o.SetPWM)
	if err != nil {
		return nil, errors.Wrap(err, "PortMegaD.MakeModel")
	}

	params = []*scripts.Param{
		{
			Code:        "value",
			Name:        "value",
			Description: "",
			Item: &models.Item{
				Type:         models.DataTypeInt,
				DefaultValue: 0,
			},
		},
	}

	impulse, err := objects.NewMethod("impulse", "Impulse", params, o.Impulse)
	if err != nil {
		return nil, errors.Wrap(err, "PortMegaD.MakeModel")
	}

	o.GetMethods().Add(check, on, off, toggle, setPWM, impulse)

	return o, nil
}

type PortModel struct {
	*objects.ObjectModelImpl

	method *method

	parentPortObjectID int // Номер родительского порта для порта расширения
	portNumber         int
	contrAddr          string

	//Number       string `json:"number_port"` //Номер порта
	//ControllerIP string `json:"controller_ip"` //IP контроллера
	//Name          string
	//Type          string // тип порта
	//ControllerID  int    //ид контроллера, на котором находится порт
	//RelatedObject int    // ИД привязанного объекта
	//EventName     string `json:"event_name"`       //Название события, которое сработало
	//Error         string `json:"error"`            //Описание текущей ошибки, если она есть
	//Params        string `json:"portParamsStruct"` //Параметры, которые пришли вместе с событием
}

// GetPortState получение состояние порта по его номеру
func (o *PortModel) GetPortState(command string, params map[string]string, timeout time.Duration) (string, error) {
	port := o.GetPortNumber()

	if o.GetParentPortObjectID() != 0 {
		parentPortNumber, err := store.I.ObjectRepository().GetProp(o.GetParentPortObjectID(), "number")
		if err != nil {
			return "", errors.Wrap(err, "GetPortState")
		}

		params["ext"] = strconv.Itoa(port)
		port, err = strconv.Atoi(parentPortNumber)
		if err != nil {
			return "", errors.Wrap(err, "GetPortState")
		}
	}

	code, content, err := o.sendCommand(o.GetContrAddr(), port, true, command, params, timeout)
	if err != nil {
		return "", errors.Wrap(err, "GetPortState")
	}

	// Если код вернулся ошибочный
	if code > 299 {
		return "", errors.Wrap(errors.Errorf("error code %d", code), "GetPortState")
	}

	return strings.ToLower(string(content)), nil
}
