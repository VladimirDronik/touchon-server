package ImpulseCounter

import (
	"github.com/pkg/errors"
	"net/http"
	"strconv"
	"strings"
	"time"
	"touchon-server/internal/g"
	"touchon-server/internal/model"
	"touchon-server/internal/objects"
	"touchon-server/internal/store"
	memStore "touchon-server/internal/store/memstore"
	"touchon-server/internal/ws"
	"touchon-server/lib/interfaces"
)

func (o *ImpulseCounter) Check(args map[string]interface{}) ([]interfaces.Message, error) {

	count, err := o.megaRelease()
	if err != nil {
		return []interfaces.Message{}, nil
	}

	//сохраняем количество импульсов в БД
	valueCount, err := o.GetProps().Get("value")
	if err != nil {
		return nil, errors.Wrap(err, "Property 'value' not found for object")
	}
	valueCount.SetValue(count)
	o.SetStatus(model.StatusAvailable)

	if err := o.Save(); err != nil {
		return nil, errors.Wrap(err, "Unable to save object")
	}

	if err := memStore.I.SaveObject(o); err != nil {
		return nil, errors.Wrap(err, "Unable to save to memory object")
	}

	msg :=
		ws.I.Send(msg)

	//генерим событие onCheck
	o.on

	msg := o.OnCheck(string(o.GetStatus()))
	return []interfaces.Message{msg}, nil

	switch newState {
	case "ON":
		o.SetStatus(model.StatusOn)
	case "OFF":
		o.SetStatus(model.StatusOff)
	}

	// заносим статус порта в БД
	go func() {
		if err := store.I.ObjectRepository().SetObjectStatus(o.GetID(), newState); err != nil {
			g.Logger.Error(errors.Wrap(err, "PortModel.Check"))
		}
	}()

	msg := o.OnChangeState(newState)
	return []interfaces.Message{msg}, nil
}

func (o *ImpulseCounter) megaRelease() (int, error) {
	addr, err := o.GetProps().GetStringValue("address")
	if err != nil {
		return 0, errors.Wrap(err, "getValues")
	}

	portObjectID, err := strconv.Atoi(addr)
	if err != nil {
		return 0, errors.Wrap(err, "getValues")
	}

	portObj, err := objects.LoadPort(portObjectID, model.ChildTypeNobody)
	if err != nil {
		return 0, errors.Wrap(err, "getValues")
	}

	var value string
	if value, err = portObj.GetPortState("get", nil, time.Duration(5)*time.Second); err != nil {
		return 0, errors.Wrap(err, "getValues")
	}

	var cnt int
	count := strings.Split(value, "/")
	if len(count) > 1 {
		cnt, err = strconv.Atoi(count[1])
		if err != nil {
			return 0, errors.Wrap(err, "getValues: bad response")
		}
	} else {
		return 0, errors.Wrap(err, "getValues: bad response")
	}

	return cnt, nil
}
