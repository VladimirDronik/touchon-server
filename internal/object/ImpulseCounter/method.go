package ImpulseCounter

import (
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
	"touchon-server/internal/g"
	helpersObj "touchon-server/internal/helpers"
	"touchon-server/internal/model"
	"touchon-server/internal/objects"
	"touchon-server/internal/store"
	"touchon-server/lib/events/object/impulse_counter"
	"touchon-server/lib/interfaces"
)

const ValueUpdateAtFormat = "02.01.2006 15:04:05"

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
	current, err := o.GetProps().GetIntValue("current")
	if err != nil {
		return errors.Wrap(err, "Property 'current' not found for object")
	}
	total, err := o.GetProps().GetFloatValue("total")
	if err != nil {
		return errors.Wrap(err, "Property 'total' not found for object")
	}
	lastUpdate, err := o.GetProps().GetStringValue("last_update")
	if err != nil {
		return errors.Wrap(err, "Property 'lastUpdate' not found for object")
	}
	today, err := o.GetProps().GetFloatValue("today")
	if err != nil {
		return errors.Wrap(err, "Property 'today' not found for object")
	}
	hour, err := o.GetProps().GetFloatValue("hour")
	if err != nil {
		return errors.Wrap(err, "Property 'hour' not found for object")
	}

	multiplier, err := o.GetProps().GetFloatValue("multiplier")
	if err != nil {
		return errors.Wrap(err, "Property 'multiplier' error getFloatValue")
	}

	// если кол-во снятых импульсов меньше хранимых, значит счетчик сбросили из вне
	d := count - current
	if d < 0 {
		d = count
	}

	// если last_update отсутствует, значит ни разу не получали еще значения со счетчика и можно проигнорить его текущие значения
	if lastUpdate == "" {
		d = 0
	}

	cur := multiplier * float32(d)
	total = total + cur
	today = today + cur
	hour = hour + cur
	//total64 := float64(totalValue)
	//ratio := math.Pow(10, float64(1))
	//totalValue = float32(math.Round(total64*ratio) / ratio)
	o.GetProps().Set("total", total)
	o.GetProps().Set("hour", hour)
	o.GetProps().Set("today", today)
	o.GetProps().Set("last_update", time.Now().Format(ValueUpdateAtFormat))

	if err := o.resetTo(0); err != nil {
		o.GetProps().Set("current", 0)
	}

	//TODO: Заносим значение в графики
	err = o.saveGraph(lastUpdate, cur, today, hour)

	o.SetStatus(model.StatusAvailable)
	g.WSServer.Send("object", model.ObjectForWS{ID: o.GetID(), Value: total})

	if err := helpersObj.SaveAndSendStatus(o, model.StatusAvailable); err != nil {
		return errors.Wrap(err, "saveImpulses")
	}

	return err
}

func (o *ImpulseCounter) saveGraph(lastUpdate string, current float32, today float32, hour float32) error {
	datetime, err := time.Parse("2006-01-02", lastUpdate)
	if err != nil {
		return err
	}

	now := time.Now()
	//Если наступил новый час, то за предыдущий сохраняем данные в БД
	if datetime.Hour() != now.Hour() {
		dateTimeMinus := now.Add(time.Duration(-1) * time.Hour)
		store.I.History().SetValue(o.GetID(), dateTimeMinus.Format("2006-01-02 15:04:05"), hour, model.TableDailyHistory)
		o.GetProps().Set("hour", 0.0)
		//Очищаем таблицу от старых данных, больше 2х недель
	}

	//Если наступил новый день, то за предыдущий сохраняем данные в БД
	if datetime.Day() != now.Day() {
		dateTimeMinus := now.Add(time.Duration(-24) * time.Hour)
		store.I.History().SetValue(o.GetID(), dateTimeMinus.Format("2006-01-02 15:04:05"), today, model.TableMonthlyHistory)
		o.GetProps().Set("today", 0.0)
		//Очищаем таблицу от старых данных больше 1 года
	}

	return nil
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
