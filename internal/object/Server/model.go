package Server

import (
	"strconv"

	"github.com/pkg/errors"
	"touchon-server/internal/helpers"
	"touchon-server/internal/model"
	"touchon-server/internal/objects"
	"touchon-server/lib/models"
)

type Server interface {
	objects.Object
}

func init() {
	_ = objects.Register(MakeModel)
}

func MakeModel(withChildren bool) (objects.Object, error) {
	props := []*objects.Prop{
		{
			Code:        "server_id",
			Name:        "Уникальный ID сервера",
			Description: "",
			Item: &models.Item{
				Type:         models.DataTypeString,
				DefaultValue: "",
			},
			Required: objects.True(),
			Editable: objects.True(),
			Visible:  objects.True(),
		},
		{
			Code:        "eco_mode",
			Name:        "Режим экономии",
			Description: "",
			Item: &models.Item{
				Type:         models.DataTypeBool,
				DefaultValue: false,
			},
			Required: objects.True(),
			Editable: objects.True(),
			Visible:  objects.True(),
		},
		{
			Code:        "guard_mode",
			Name:        "Режим охраны",
			Description: "",
			Item: &models.Item{
				Type:         models.DataTypeBool,
				DefaultValue: false,
			},
			Required: objects.True(),
			Editable: objects.True(),
			Visible:  objects.True(),
		},
		{
			Code:        "night_mode",
			Name:        "Ночной режим",
			Description: "",
			Item: &models.Item{
				Type:         models.DataTypeBool,
				DefaultValue: false,
			},
			Required: objects.True(),
			Editable: objects.True(),
			Visible:  objects.True(),
		},
		{
			Code:        "heating_mode",
			Name:        "План отопления дома",
			Description: "",
			Item: &models.Item{
				Type: models.DataTypeEnum,
				Values: map[string]string{
					"eco":    "eco",
					"night":  "night",
					"normal": "normal",
				},
				DefaultValue: "eco",
			},
			Required: objects.True(),
			Editable: objects.True(),
			Visible:  objects.True(),
		},
		{
			Code:        "light_mode",
			Name:        "Режим освещения",
			Description: "",
			Item: &models.Item{
				Type: models.DataTypeEnum,
				Values: map[string]string{
					"night":   "night",
					"day":     "day",
					"evening": "evening",
				},
				DefaultValue: "day",
			},
			Required: objects.True(),
			Editable: objects.True(),
			Visible:  objects.True(),
		},
		{
			Code:        "logging",
			Name:        "Где хранить логи",
			Description: "",
			Item: &models.Item{
				Type: models.DataTypeEnum,
				Values: map[string]string{
					"file": "file",
					"db":   "DB",
				},
				DefaultValue: "file",
			},
			Required: objects.True(),
			Editable: objects.True(),
			Visible:  objects.True(),
		},
		{
			Code:        "storage_logs",
			Name:        "Количество дней хранения логов",
			Description: "",
			Item: &models.Item{
				Type:         models.DataTypeInt,
				DefaultValue: 30,
			},
			Required: objects.True(),
			Editable: objects.True(),
			Visible:  objects.True(),
		},
		{
			Code:        "graph_date",
			Name:        "Сколько дней хранить информацию в графиках",
			Description: "",
			Item: &models.Item{
				Type:         models.DataTypeInt,
				DefaultValue: 365,
			},
			Required: objects.True(),
			Editable: objects.True(),
			Visible:  objects.True(),
		},
		{
			Code:        "time_zone",
			Name:        "Часовой пояс приложения",
			Description: "",
			Item: &models.Item{
				Type:         models.DataTypeString,
				DefaultValue: "Europe/Moscow",
			},
			Required: objects.True(),
			Editable: objects.True(),
			Visible:  objects.True(),
		},
	}

	impl, err := objects.NewObjectModelImpl(
		model.CategoryServer,
		"server",
		objects.CreationForbidden|objects.DeletionForbidden,
		"Сервер",
		props,
		nil,
		nil,
		nil,
		[]string{"server"},
	)
	if err != nil {
		return nil, errors.Wrap(err, "Server.MakeModel")
	}

	o := &ServerImpl{
		Object: impl,
	}

	if err := o.GetProps().Set("server_id", helpers.MD5(strconv.Itoa(helpers.Rnd.Int()))); err != nil {
		return nil, errors.Wrap(err, "Server.MakeModel")
	}

	return o, nil
}

type ServerImpl struct {
	objects.Object
}

func (o *ServerImpl) Start() error {
	if err := o.Object.Start(); err != nil {
		return errors.Wrap(err, "ServerImpl.Start")
	}

	// todo...

	return nil
}

func (o *ServerImpl) Shutdown() error {
	if err := o.Object.Shutdown(); err != nil {
		return errors.Wrap(err, "ServerImpl.Shutdown")
	}

	// todo...

	return nil
}
