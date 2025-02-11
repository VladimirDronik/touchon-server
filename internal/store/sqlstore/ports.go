package sqlstore

import (
	"strings"

	"github.com/pkg/errors"
)

type PortRepository struct {
	store *Store
}

// GetPortObjectID Найти ИД порта контроллера по его номеру и адресу контроллера
func (o *PortRepository) GetPortObjectID(controllerID, portNumber string) (int, error) {
	type R struct {
		ID int
	}
	rows := make([]R, 0, 10)

	q := `
SELECT o.id
FROM props as p
    inner join objects as o on o.id = p.object_id
WHERE o.category == 'port' AND p.code = 'number' AND p.value = ? AND o.parent_id in (
    SELECT o.object_id FROM props as o
    WHERE o.code == 'id' AND o.value = ?
);`

	if err := o.store.db.Raw(q, portNumber, strings.TrimSpace(controllerID)).Scan(&rows).Error; err != nil {
		return 0, errors.Wrapf(err, "GetPortObjectID(%s, %s)", controllerID, portNumber)
	}

	switch len(rows) {
	case 0:
		return 0, errors.Wrapf(errors.New("port not found"), "GetPortObjectID(%s, %s)", controllerID, portNumber)
	case 1:
		return rows[0].ID, nil
	default:
		return 0, errors.Wrapf(errors.Errorf("%d port found", len(rows)), "GetPortObjectID(%s, %s)", controllerID, portNumber)
	}
}
