package ImpulseCounter

import (
	"github.com/pkg/errors"
	"strconv"
	"strings"
	"time"
	"touchon-server/internal/g"
	helpersObj "touchon-server/internal/helpers"
	"touchon-server/internal/model"
	"touchon-server/internal/objects"
	"touchon-server/lib/events/object/impulse_counter"
	"touchon-server/lib/events/object/onokom/gateway"
	"touchon-server/lib/interfaces"
)

func (o *ImpulseCounter) Check(params map[string]interface{}) ([]interfaces.Message, error) {
	count, err := o.megaRelease()
	if err != nil {
		return []interfaces.Message{}, nil
	}

	err = o.saveImpulses(count)
	if err != nil {
		return []interfaces.Message{},
			errors.Wrap(err, "ModelImpulseCounter.Check: save impulse data to DB failed")
	}

	//генерим событие onCheck
	impulseCntMsg, err := impulse_counter.NewOnCheck(o.GetID(), count)
	if err != nil {
		return nil, errors.Wrap(err, "ModelImpulseCounter.Check")
	}

	return []interfaces.Message{impulseCntMsg}, nil
}

func (o *ImpulseCounter) check() {
	defer o.GetTimer().Reset()
	_, err := o.Check(nil)
	if err != nil {
		g.Logger.Error(errors.Wrap(err, "ImpulseCounter.check"))
		return
	}

	payload := make(map[string]interface{}, o.GetProps().Len())

	payload["current_count"], err = o.GetProps().GetIntValue("value")
	if err != nil {
		g.Logger.Error(errors.Wrap(err, "ImpulseCounter.check"))
		return
	}

	// Отправляем сообщение с измененными полями
	msg, err := gateway.NewOnChange(o.GetID(), payload)
	if err != nil {
		g.Logger.Error(errors.Wrap(err, "ImpulseCounter.check"))
		return
	}

	// Отправляем сообщение об изменении св-ва объекта
	if err := g.Msgs.Send(msg); err != nil {
		g.Logger.Error(errors.Wrap(err, "ImpulseCounter.check"))
		return
	}

	return
}

func (o *ImpulseCounter) megaRelease() (int, error) {
	portObj, err := o.getPort()
	if err != nil {
		return 0, errors.Wrap(err, "getValues: get port for counter is fault")
	}

	var value string
	if value, err = portObj.GetPortState("get", nil, time.Duration(5)*time.Second); err != nil {
		//TODO: добавить отправку статуса объекта в сокет и изменение в БД
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

func (o *ImpulseCounter) getPort() (interfaces.Port, error) {
	addr, err := o.GetProps().GetStringValue("address")
	if err != nil {
		return nil, errors.Wrap(err, "getValues")
	}

	portObjectID, err := strconv.Atoi(addr)
	if err != nil {
		return nil, errors.Wrap(err, "getValues")
	}

	portObj, err := objects.LoadPort(portObjectID, false)
	if err != nil {
		return nil, errors.Wrap(err, "getValues")
	}

	return portObj, nil
}

// сохраняем количество импульсов в БД
func (o *ImpulseCounter) saveImpulses(count int) error {
	valueCount, err := o.GetProps().Get("value")
	if err != nil {
		return errors.Wrap(err, "Property 'value' not found for object")
	}
	valueCount.SetValue(count)

	err = helpersObj.SaveAndSendStatus(o, model.StatusAvailable)

	return err
}

func (o *ImpulseCounter) resetTo(val int) error {
	portObj, err := o.getPort()
	if err != nil {
		return errors.Wrap(err, "getValues")
	}

	params := make(map[string]string)
	params["cnt"] = strconv.Itoa(val)
	code, err := portObj.SetPortParams(params)
	if err != nil || code < 299 {
		return errors.Wrap(err, "resetTo: reset counter fault")
	}

	return nil
}
