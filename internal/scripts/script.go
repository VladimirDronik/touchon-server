package scripts

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"touchon-server/lib/models"
)

func NewScript(code, name, description string, params []*Param, body string) (*Script, error) {
	s := &Script{
		Code:        code,
		Name:        name,
		Description: description,
		Params:      NewParams(),
		Body:        body,
	}

	for _, p := range params {
		if err := s.Params.Add(p); err != nil {
			return nil, errors.Wrap(err, "NewScript")
		}
	}

	if s.Body == "" {
		s.Body = s.GenerateStub()
	}

	return s, nil
}

type Script struct {
	ID          int     `json:"id"`
	Code        string  `json:"code"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Params      *Params `json:"params"`
	Body        string  `json:"body"`
}

func (o *Script) Check() error {
	switch {
	case o.Code == "":
		return errors.New("script Code is empty")
	case o.Name == "":
		return errors.New("script Name is empty")
	case o.Params == nil:
		return errors.New("script Params is nil")
	case o.Body == "":
		return errors.New("script Body is empty")
	}

	if err := o.Params.Check(); err != nil {
		return err
	}

	return nil
}

func (o *Script) GenerateStub() string {
	v := make([]string, 0, o.Params.Len())
	i := make([]string, 0, o.Params.Len())

	for _, p := range o.Params.m.GetValueList() {
		v = append(v, fmt.Sprintf("var %s %s // %s\n", p.Code, models.DataTypeToGoType[p.Type], p.Name))
		i = append(i, fmt.Sprintf("%[1]s, _ = args[%[1]q].(%s)\n", p.Code, models.DataTypeToGoType[p.Type]))
	}

	s := `package main
import s "scripts"

// Аргументы сценария
%[1]s

func init() {
	args := s.Args()
	%[2]s
}

func main() {
	// Код сценария
	// ...

	s.Ok("sleep " + s.ToString(timeout))
	// s.Err("status is err")

	// Вспомогательные функции преобразования типов
	// --------------------------------------------
	// s.ToInt("123")
	// s.ToFloat("1.5")
	// s.ToBool("true")
	// s.ToString(123)
	// s.ToString(false)
	//
	// Выполнение другого скрипта
	// --------------------------
	// Сигнатура:
	// func(code string, args map[string]interface{}) (interface{}, string)
	// Пример:
	// r, err := s.Exec("my_second_script", map[string]interface{}{"ObjectID":88})
	// if err != "" {
	//     s.Err(err)	
	// }
	//
	// Выполнение метода объекта
	// -------------------------
	// Выполнение метода объекта с ID=1
	// r, err := s.ExecObjectMethod(1, "", "", "method", map[string]interface{}{"logicArg":false})
	// s.Ok(r)
	//
	// Выполнение метода объектов категории controller и с типом mega_d
	// r, err := s.ExecObjectMethod(0, "controller", "mega_d", "method", nil)
	// s.Err(err)
}
`

	s = strings.Replace(s, "%[1]s", strings.TrimSpace(strings.Join(v, "")), 1)
	s = strings.Replace(s, "%[2]s", strings.TrimSpace(strings.Join(i, "")), 1)

	return s
}
