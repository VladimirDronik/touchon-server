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
	"touchon-server/internal/ws"
	"touchon-server/lib/events/object/impulse_counter"
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
	msg, err := o.Check(nil)
	if err != nil {
		g.Logger.Error(errors.Wrap(err, "ImpulseCounter.check"))
		return
	}

	// Отправляем сообщение об изменении св-ва объекта
	for _, m := range msg {
		if err := g.Msgs.Send(m); err != nil {
			g.Logger.Error(errors.Wrap(err, "ImpulseCounter.check"))
			return
		}
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
	portObjectID, err := o.GetProps().GetIntValue("address")
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
	currentProp, err := o.GetProps().Get("current")
	if err != nil {
		return errors.Wrap(err, "Property 'current' not found for object")
	}
	totalProp, err := o.GetProps().Get("total")
	if err != nil {
		return errors.Wrap(err, "Property 'total' not found for object")
	}

	current, err := currentProp.GetIntValue()
	if err != nil {
		return errors.Wrap(err, "Property 'current' error getIntValue")
	}

	total, err := totalProp.GetFloatValue()
	if err != nil {
		return errors.Wrap(err, "Property 'total' error getIntValue")
	}

	multiplier, err := o.GetProps().GetFloatValue("multiplier")
	if err != nil {
		return errors.Wrap(err, "Property 'multiplier' error getFloatValue")
	}

	//если кол-во снятых импульсов меньше хранимых, значит счетчик сбросили из вне
	d := count - current
	if d < 0 {
		d = count
	}

	totalValue := total * (multiplier + float32(d))
	totalProp.SetValue(total * (multiplier + float32(d)))

	if err := o.resetTo(0); err != nil {
		currentProp.SetValue(0)
	}

	//TODO: Заносим значение в графики

	ws.I.Send("object", model.ObjectForWS{ID: o.GetID(), Value: totalValue})
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
