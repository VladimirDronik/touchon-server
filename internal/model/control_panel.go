package model

// ViewItem Структура основных итемов (кнопок панели управления)
type ViewItem struct {
	ID            int    `json:"item_id"`                           //
	ParentID      *int   `json:"id_group,omitempty"`                //
	ZoneID        *int   `json:"zone_id,omitempty"`                 //
	Type          string `json:"type,omitempty"`                    //
	Status        string `json:"status,omitempty"`                  //
	Icon          string `json:"icon,omitempty"`                    //
	Title         string `json:"title,omitempty"`                   //
	Sort          int    `json:"sort"`                              //
	Params        string `json:"-"`                                 // Используется для выборки из БД, но не для выдачи в ответе сервера
	Color         string `json:"color,omitempty"`                   //
	Auth          string `json:"auth,omitempty"`                    //
	Description   string `json:"description"`                       //
	PositionLeft  int    `json:"position_left"`                     //
	PositionTop   int    `json:"position_top"`                      //
	Scene         int    `json:"scene"`                             //
	Enabled       bool   `json:"enabled"`                           //
	ControlObject int    `json:"control_object,omitempty" gorm:"-"` // Объект, статус которого влияет на статус итема

	Value        float32     `json:"value,omitempty" gorm:"-"`  //
	ParamsOutput interface{} `json:"params,omitempty" gorm:"-"` // Используется для вывода в ответе сервера, но не для выборки из БД

	GroupElements []*ViewItem    `json:"group_elements,omitempty" gorm:"-"` //
	History       *HistoryPoints `json:"history,omitempty" gorm:"-"`        //
}

// Scenario Структура сценария
type Scenario struct {
	ID          int    `json:"-"`
	ViewItemID  int    `json:"item_id"`
	Type        string `json:"type"`
	Description string `json:"description"`
	Icon        string `json:"icon,omitempty"`
	Title       string `json:"title,omitempty"`
	Sort        int    `json:"sort"`
	Color       string `json:"color,omitempty"`
	Auth        string `json:"auth,omitempty"`
	Enabled     bool   `json:"enabled"`

	// Backward compatibility
	Status string `json:"status"`
}

// GroupRoom Структура группы для вывода в responseGroup
type GroupRoom struct {
	ID       int    `json:"id"`
	ParentID int    `json:"parent_id"`
	Name     string `json:"name"`
	Style    string `json:"style"`
	Sort     int    `json:"sort"`
	IsGroup  bool   `json:"is_group"`

	Scenarios []*ViewItem   `json:"scenario,omitempty"`
	Sensors   []*SensorItem `json:"sensors,omitempty"`
	Items     []*ViewItem   `json:"items,omitempty"`
}

// Zone Структура для вывода в помещений
type Zone struct {
	ID       int    `json:"id"`
	ParentID *int   `json:"parent_id"`
	Name     string `json:"name"`
	Style    string `json:"style"`
	Sort     int    `json:"sort"`
	IsGroup  bool   `json:"is_group"`

	Children []*Zone `json:"rooms_in_group,omitempty" gorm:"-"`
}

func (o *Zone) TableName() string {
	return "zones"
}

// Children Структура комнаты для вывода внутри группы
// Когда пользователь нажимает кнопку настроек для помещения, то ему может вывестись список комнат, которые входят
// в группу этого помещения
//type Children struct {
//	RoomID    int    `json:"room_id" gorm:"room_id"`
//	RoomName  string `json:"room_name" gorm:"room_name"`
//	RoomImage string `json:"room_image" gorm:"room_image"`
//}
