package scripts

import (
	"bytes"
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"reflect"
	"runtime"
	"strconv"
	"time"

	"github.com/pkg/errors"
	"github.com/traefik/yaegi/interp"
	svcContext "touchon-server/internal/context"
	"touchon-server/internal/msgs"
	"touchon-server/internal/store"
	"touchon-server/lib/events/script"
	"touchon-server/lib/helpers/orderedmap"
	"touchon-server/lib/interfaces"
)

// Global instance
var I *Scripts

type ObjectMethodExecutor func(id int, objCategory, objType, method string, args map[string]interface{}) ([]interface{}, string)

func NewScripts(timeout time.Duration, objectMethodExecutor ObjectMethodExecutor) *Scripts {
	return &Scripts{
		timeout:              timeout,
		objectMethodExecutor: objectMethodExecutor,
	}
}

type Scripts struct {
	timeout time.Duration

	// Обработчик методов
	objectMethodExecutor ObjectMethodExecutor
}

func (o *Scripts) GetScript(id int) (*Script, error) {
	s, err := store.I.ScriptRepository().GetScript(id)
	if err != nil {
		return nil, errors.Wrap(err, "GetScript")
	}

	script := &Script{
		ID:          s.ID,
		Code:        s.Code,
		Name:        s.Name,
		Description: s.Description,
		Body:        s.Body,
	}

	script.Params = &Params{
		m: orderedmap.New[string, *Param](10),
	}

	if err := json.Unmarshal(s.Params, &script.Params); err != nil {
		return nil, errors.Wrap(err, "GetScript")
	}

	// Проверяем скрипт перед отдачей
	if err := script.Check(); err != nil {
		return nil, errors.Wrap(err, "GetScript")
	}

	return script, nil
}

func (o *Scripts) GetScriptByCode(code string) (*Script, error) {
	s, err := store.I.ScriptRepository().GetScriptByCode(code)
	if err != nil {
		return nil, errors.Wrap(err, "GetScriptByCode")
	}

	script := &Script{
		ID:          s.ID,
		Code:        s.Code,
		Name:        s.Name,
		Description: s.Description,
		Body:        s.Body,
	}

	script.Params = &Params{
		m: orderedmap.New[string, *Param](10),
	}

	if err := json.Unmarshal(s.Params, &script.Params); err != nil {
		return nil, errors.Wrap(err, "GetScriptByCode")
	}

	// Проверяем скрипт перед отдачей
	if err := script.Check(); err != nil {
		return nil, errors.Wrap(err, "GetScriptByCode")
	}

	return script, nil
}

func (o *Scripts) ExecScript(s *Script, args map[string]interface{}) (interface{}, error) {
	if s.Params.Len() != len(args) {
		return nil, errors.Wrapf(errors.New("недостаточно аргументов скрипта"), "ExecScript(%s)", s.Code)
	}

	// Выставляем значения параметров, для проверки типов
	for _, p := range s.Params.m.GetValueList() {
		v, ok := args[p.Code]
		if !ok {
			return nil, errors.Wrapf(errors.Errorf("arg %q not found", p.Code), "ExecScript(%s)", s.Code)
		}

		if err := p.SetValue(v); err != nil {
			return nil, errors.Wrapf(err, "ExecScript(%s)", s.Code)
		}
	}

	r, err := o.execScript(s, args)
	if err != nil {
		return nil, errors.Wrapf(err, "ExecScript(%s)", s.Code)
	}

	return r, nil
}

func (o *Scripts) Exec(code string, args map[string]interface{}) (interface{}, error) {
	s, err := o.GetScriptByCode(code)
	if err != nil {
		return nil, errors.Wrapf(err, "Exec(%s)", code)
	}

	r, err := o.ExecScript(s, args)
	if err != nil {
		return nil, errors.Wrapf(err, "Exec(%s)", code)
	}

	return r, nil
}

type scriptResponse struct {
	Data  interface{} `json:"data,omitempty"`
	Error string      `json:"error,omitempty"`
}

func (o *Scripts) execScript(script *Script, args map[string]interface{}) (interface{}, error) {
	m := script.Params.m.GetUnorderedMap()

	for k, arg := range args {
		p, ok := m[k]
		if !ok {
			return nil, errors.Wrap(errors.Errorf("unknown arg %q", k), "execScript")
		}

		// Проверяем тип и преобразуем его в нужный формат
		if err := p.SetValue(arg); err != nil {
			return nil, errors.Wrap(err, "execScript")
		}
		args[k] = p.GetValue()

		log.Printf("%[1]v %[1]T", args[k])
	}

	var stdOut bytes.Buffer

	opts := interp.Options{
		BuildTags:            nil,
		Stdin:                &bytes.Buffer{},
		Stdout:               &stdOut,
		Stderr:               &bytes.Buffer{},
		SourcecodeFilesystem: embed.FS{},
		Unrestricted:         false,
	}

	i := interp.New(opts)
	export := make(map[string]map[string]reflect.Value)
	export["scripts/scripts"] = make(map[string]reflect.Value)

	// Функция запуска другого скрипта
	export["scripts/scripts"]["Exec"] = reflect.ValueOf(func(code string, args map[string]interface{}) (interface{}, string) {
		// Преобразуем error в string при вызове из скрипта
		r, err := o.Exec(code, args)
		if err != nil {
			return r, err.Error()
		}

		return r, ""
	})

	// Передаем аргументы запуска
	export["scripts/scripts"]["Args"] = reflect.ValueOf(func() map[string]interface{} {
		// Копируем аргументы скрипта, чтобы их нельзя было изменить из самого скрипта
		a := make(map[string]interface{}, len(args))
		for k, v := range args {
			a[k] = v
		}

		return a
	})

	// Функция возврата результата
	export["scripts/scripts"]["Ok"] = reflect.ValueOf(func(data interface{}) {
		returnResult(&stdOut, data)
	})

	// Функция возврата ошибки
	export["scripts/scripts"]["Err"] = reflect.ValueOf(func(errMsg string) {
		returnError(&stdOut, errMsg)
	})

	export["scripts/scripts"]["ToInt"] = reflect.ValueOf(toInt)
	export["scripts/scripts"]["ToFloat"] = reflect.ValueOf(toFloat)
	export["scripts/scripts"]["ToBool"] = reflect.ValueOf(toBool)
	export["scripts/scripts"]["ToString"] = reflect.ValueOf(toString)

	if o.objectMethodExecutor != nil {
		export["scripts/scripts"]["ExecObjectMethod"] = reflect.ValueOf(o.objectMethodExecutor)
	}

	if err := i.Use(export); err != nil {
		return nil, errors.Wrap(err, "execScript")
	}

	ctx, cancel := context.WithTimeout(context.Background(), o.timeout)
	defer cancel()

	if _, err := i.EvalWithContext(ctx, script.Body); err != nil {
		return nil, errors.Wrap(err, "execScript")
	}

	data := stdOut.Bytes()

	// Если получили json
	if len(data) > 0 && data[0] == '{' && data[len(data)-1] == '}' {
		r := &scriptResponse{}
		if err := json.Unmarshal(data, r); err != nil {
			return nil, errors.Wrap(err, "execScript")
		}

		var err error
		if r.Error != "" {
			err = errors.New(r.Error)
		}

		return r.Data, err
	}

	// Если получили текст или пустой вывод
	return string(data), nil
}

// MsgHandler позволяет обрабатывать сообщения из брокера сообщений
func (o *Scripts) MsgHandler(msg interfaces.Message) {
	s, err := o.GetScript(msg.GetTargetID())
	if err != nil {
		svcContext.Logger.Error(errors.Wrap(err, "Scripts.MsgHandler"))
		return
	}

	// TODO
	r, err := o.ExecScript(s, nil) // msg.GetPayload()
	if err != nil {
		svcContext.Logger.Error(errors.Wrap(err, "Scripts.MsgHandler"))
		return
	}

	msg, err = script.NewOnComplete(msg.GetTargetID(), r)
	if err != nil {
		svcContext.Logger.Error(errors.Wrap(err, "Scripts.MsgHandler"))
		return
	}

	if err := msgs.I.Send(msg); err != nil {
		svcContext.Logger.Error(errors.Wrap(err, "Scripts.MsgHandler"))
	}
}

func returnResult(stdOut io.Writer, data interface{}) {
	defer runtime.Goexit()

	r, err := json.Marshal(scriptResponse{Data: data})
	if err != nil {
		_, _ = fmt.Fprintf(stdOut, `{"error":"%s"}`, err.Error())
		return
	}

	_, _ = fmt.Fprint(stdOut, string(r))
}

func returnError(stdOut io.Writer, errMsg string) {
	defer runtime.Goexit()

	r, err := json.Marshal(scriptResponse{Error: errMsg})
	if err != nil {
		_, _ = fmt.Fprintf(stdOut, `{"error":"%s"}`, err.Error())
		return
	}

	_, _ = fmt.Fprint(stdOut, string(r))
}

func toInt(v interface{}) int {
	switch v := v.(type) {
	case string:
		i, _ := strconv.Atoi(v)
		return i
	case int:
		return v
	case int64:
		return int(v)
	case float32:
		return int(v)
	case float64:
		return int(v)
	default:
		return 0
	}
}

func toFloat(v interface{}) float32 {
	switch v := v.(type) {
	case string:
		f, _ := strconv.ParseFloat(v, 32)
		return float32(f)
	case int:
		return float32(v)
	case int64:
		return float32(v)
	case float32:
		return v
	case float64:
		return float32(v) // trancate value!
	default:
		return 0
	}
}

func toBool(v interface{}) bool {
	switch v := v.(type) {
	case string:
		b, _ := strconv.ParseBool(v)
		return b
	case bool:
		return v
	default:
		return false
	}
}

func toString(v interface{}) string {
	return fmt.Sprintf("%v", v)
}
