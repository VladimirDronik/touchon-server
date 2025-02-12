package script

import (
	"github.com/pkg/errors"
	"touchon-server/lib/interfaces"
	"touchon-server/lib/messages"
)

func NewOnComplete(scriptID int, result interface{}) (interfaces.Message, error) {
	msg, err := messages.NewEvent(interfaces.TargetTypeScript, scriptID)
	if err != nil {
		return nil, errors.Wrap(err, "NewOnComplete")
	}

	o := &OnComplete{
		MessageImpl: msg,
		Result:      result,
	}

	return o, nil
}

type OnComplete struct {
	*messages.MessageImpl
	Result interface{} `json:"result"` // Результат работы скрипта
}
