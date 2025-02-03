// https://ab-log.ru/smart-house/ethernet/megad-2561

package MegaD

import (
	"fmt"
	"net"
	"sort"
	"strconv"

	"github.com/pkg/errors"
	"touchon-server/internal/model"
	"touchon-server/internal/object/PortMegaD"
	"touchon-server/internal/objects"
	"touchon-server/internal/store"
	"touchon-server/lib/models"
)

func init() {
	_ = objects.Register(MakeModel)
}

func MakeModel() (objects.Object, error) {
	props := []*objects.Prop{
		{
			Code:        "id",
			Name:        "Идентификатор контроллера",
			Description: "",
			Item: &models.Item{
				Type:         models.DataTypeString,
				DefaultValue: "",
			},
			Required: objects.NewRequired(true),
			Editable: objects.NewCondition(),
			Visible:  objects.NewCondition(),
		},
		{
			Code:        "password",
			Name:        "Пароль",
			Description: "",
			Item: &models.Item{
				Type:         models.DataTypeString,
				DefaultValue: "sec",
			},
			Required: objects.NewRequired(true),
			Editable: objects.NewCondition(),
			Visible:  objects.NewCondition(),
		},
		{
			Code:        "address",
			Name:        "IP адрес",
			Description: "",
			Item: &models.Item{
				Type:         models.DataTypeString,
				DefaultValue: "127.0.0.1",
			},
			Required: objects.NewRequired(true),
			Editable: objects.NewCondition(),
			Visible:  objects.NewCondition(),
			CheckValue: func(p *objects.Prop, allProps map[string]*objects.Prop) error {
				if v, _ := p.GetStringValue(); net.ParseIP(v) == nil {
					return errors.New("IP address is bad")
				}

				return nil
			},
		},
		{
			Code:        "protocol",
			Name:        "Протокол",
			Description: "",
			Item: &models.Item{
				Type:         models.DataTypeEnum,
				Values:       map[string]string{"http": "http", "mqtt": "mqtt"},
				DefaultValue: "http",
			},
			Required: objects.NewRequired(true),
			Editable: objects.NewCondition(),
			Visible:  objects.NewCondition(),
		},
		{
			Code:        "module",
			Name:        "Модуль",
			Description: "",
			Item: &models.Item{
				Type:         models.DataTypeBool,
				DefaultValue: false,
			},
			Required: objects.NewRequired(false),
			Editable: objects.NewCondition().AccessLevel(model.AccessLevelDenied),
			Visible:  objects.NewCondition().AccessLevel(model.AccessLevelDenied),
		},
	}

	// Создаем сортированный список портов
	// Чтобы в базу порты записывались в нужном порядке
	ports := make([]int, 0, Ports.Len())
	for _, port := range Ports.GetKeyValueList() {
		ports = append(ports, port.Key)
	}
	sort.Ints(ports)

	children := make([]objects.Object, 0, 45)
	for _, portNumber := range ports {
		// Создаем основу объекта-порта
		port, err := PortMegaD.MakeModel()
		if err != nil {
			return nil, err
		}

		// Задаем имя порта
		port.SetName("Порт " + strconv.Itoa(portNumber))

		// Задаем номер порта
		if err := port.GetProps().Set("number", portNumber); err != nil {
			return nil, err
		}

		p, err := Ports.Get(portNumber)
		if err != nil {
			return nil, err
		}

		// Задаем принадлежность к группе портов

		g, err := port.GetProps().Get("group")
		if err != nil {
			return nil, err
		}

		g.Type = models.DataTypeEnum
		g.Values = make(map[string]string, len(PortGroupList))
		for _, item := range PortGroupList {
			g.Values[item.Code] = item.Name
		}

		if err := g.SetValue(p.Group); err != nil {
			return nil, err
		}

		// Задаем список типов, поддерживаемых портом

		t, err := port.GetProps().Get("type")
		if err != nil {
			return nil, err
		}

		t.Type = models.DataTypeEnum
		t.Values = make(map[string]string, p.Types.Len())
		for _, item := range p.Types.GetValueList() {
			t.Values[item.Code] = item.Name
		}

		// Задаем тип по умолчанию.
		if p.Types.Len() > 0 {
			t.DefaultValue = p.Types.GetValueList()[0].Code
		}

		// Задаем список режимов, поддерживаемых портом

		m, err := port.GetProps().Get("mode")
		if err != nil {
			return nil, err
		}

		m.Type = models.DataTypeEnum
		m.Values = make(map[string]string, 25)
		for _, portType := range p.Types.GetValueList() {
			for _, mode := range portType.Modes.GetValueList() {
				m.Values[mode.Code] = fmt.Sprintf("%s / %s", portType.Name, mode.Name)
			}
		}

		// Задаем режим по умолчанию.
		if p.Types.Len() > 0 && p.Types.GetValueList()[0].Modes.Len() > 0 {
			m.DefaultValue = p.Types.GetValueList()[0].Modes.GetValueList()[0].Code
		}

		m.CheckValue = func(prop *objects.Prop, allProps map[string]*objects.Prop) error {
			portNumber, err := allProps["number"].GetIntValue()
			if err != nil {
				return err
			}

			typeCode, err := allProps["type"].GetStringValue()
			if err != nil {
				return err
			}

			modeCode, err := prop.GetStringValue()
			if err != nil {
				return err
			}

			port, err := Ports.Get(portNumber)
			if err != nil {
				return err
			}

			portType, err := port.Types.Get(typeCode)
			if err != nil {
				return err
			}

			if _, err := portType.Modes.Get(modeCode); err != nil {
				return err
			}

			return nil
		}

		children = append(children, port)
	}

	impl, err := objects.NewObjectModelImpl(
		model.CategoryController,
		"mega_d",
		false,
		"MegaD",
		props,
		children,
		[]string{"object.controller.on_load", "object.controller.on_unavailable"},
		nil,
		[]string{"controller", "mega_d"},
	)
	if err != nil {
		return nil, errors.Wrap(err, "MegaD.MakeModel")
	}

	return &ObjectModel{ObjectModelImpl: impl}, nil
}

// Implementation of Object interface

type ObjectModel struct {
	*objects.ObjectModelImpl
}

func (o *ObjectModel) DeleteChildren() error {
	for _, child := range o.GetChildren().GetAll() {
		// При удалении контроллера дочерние объекты не обрабатываем
		//if err := child.DeleteChildren(); err != nil {
		//	return errors.Wrap(err, "DeleteChildren")
		//}

		// Удаляем только порты контроллера
		if !(child.GetCategory() == model.CategoryPort && child.GetType() == "port_mega_d") {
			continue
		}

		if err := store.I.ObjectRepository().DelObject(child.GetID()); err != nil {
			return errors.Wrap(err, "DeleteChildren")
		}
	}

	return nil
}
