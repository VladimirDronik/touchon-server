package interfaces

import "time"

type Port interface {
	GetPortState(command string, params map[string]string, timeout time.Duration) (string, error)
	On(args map[string]interface{}) ([]Message, error)
	Off(args map[string]interface{}) ([]Message, error)
	Toggle(args map[string]interface{}) ([]Message, error)
	SetTypeMode(typePt string, modePt string, title string, params map[string]string) error
	SetPortParams(params map[string]string) (int, error)
}
