package MegaD

import (
	"sort"

	"touchon-server/lib/helpers/orderedmap"
)

// ============================================================
// Работаем с этими переменными

// Список всех портов
var Ports = orderedmap.New[int, *Port](50)

// Список портов по группам
var Groups = orderedmap.New[string, []*Port](3)

// Список типов портов
var PortTypes = orderedmap.New[string, *PortType](10)

// ============================================================
// Здесь настраиваем группы/типы/режимы портов

var PortGroupList = []*portGroup{
	{Code: "inputs", Name: "Входы", Description: ""},
	{Code: "outputs", Name: "Выходы", Description: ""},
	{Code: "digital", Name: "Цифровые", Description: ""},
}

var portTypeList = []*portType{
	{
		Code:        "nc",
		Name:        "NC",
		Description: "",
		Modes: []*PortMode{
			{Code: "", Name: "Нет режима", Description: ""},
		},
	},
	{
		Code:        "in",
		Name:        "IN",
		Description: "",
		Modes: []*PortMode{
			{Code: "p", Name: "P", Description: "Порт срабатывает только на замыкание (press)"},
			{Code: "r", Name: "R", Description: "Порт срабатывает только на размыкание (release)"},
			{Code: "pr", Name: "P&R", Description: "Порт срабатывает как на замыкание, так и на размыкание (press & release)"},
			{Code: "c", Name: "C", Description: "Обработка одинарных и двойных кликов/нажатий (Click Mode)"},
		},
	},
	{
		Code:        "out",
		Name:        "OUT",
		Description: "",
		Modes: []*PortMode{
			{Code: "sw", Name: "SW", Description: ""},
			{Code: "swlink", Name: "SW LINK", Description: "Режим SW LINK (связанный порт) используется в тех случаях, когда необходимо предотвратить на аппаратном уровне включение двух или более портов одновременно. Это актуально для управления различными приводами (рольставен, кранов, клапанов). Данный режим предполагает, что если один из связанных портов включен, то любой другой включить уже нельзя. Для того, чтобы связать два или более портов, необходимо для каждого из них выбрать режим SW LINK и указать одинаковую группу в поле Group (любое число от 0 до 99)."},
			{Code: "pwm", Name: "PWM", Description: ""},
		},
	},
	{
		Code:        "dsen",
		Name:        "DSen",
		Description: "",
		Modes: []*PortMode{
			{Code: "dht11", Name: "DHT11", Description: ""},
			{Code: "dht22", Name: "DHT22", Description: ""},
			{Code: "1w", Name: "1W", Description: ""},
			{Code: "1wbus", Name: "1WBUS", Description: ""},
			{Code: "ib", Name: "iB", Description: ""},
			{Code: "w26", Name: "W26", Description: ""},
		},
	},
	{
		Code:        "i2c",
		Name:        "I2C",
		Description: "",
		Modes: []*PortMode{
			{Code: "nc", Name: "NC", Description: ""},
			{Code: "sda", Name: "SDA", Description: ""},
			{Code: "scl", Name: "SCL", Description: ""},
		},
	},
	{
		Code:        "adc",
		Name:        "ADC",
		Description: "",
		Modes: []*PortMode{
			{Code: "norm", Name: "NORM", Description: ""},
		},
	},
}

var portList = []*port{
	{
		Numbers: []int{0, 1, 2, 3, 4, 5},
		Group:   "inputs",
		Types:   []string{"nc", "in", "dsen", "i2c", "adc"},
		Modes:   []string{"", "p", "pr", "r", "c", "dht11", "dht22", "1w", "1wbus", "ib", "w26", "nc", "sda", "scl", "norm"},
	},
	{
		Numbers: []int{6, 15, 16, 17, 18, 19, 20, 21},
		Group:   "inputs",
		Types:   []string{"nc", "in", "dsen", "i2c"},
		Modes:   []string{"", "p", "pr", "r", "c", "dht11", "dht22", "1w", "1wbus", "ib", "w26", "nc", "sda", "scl"},
	},
	{
		Numbers: []int{7, 8, 9, 22, 23, 24, 26},
		Group:   "outputs",
		Types:   []string{"nc", "out"},
		Modes:   []string{"", "sw", "swlink"},
	},
	{
		Numbers: []int{10, 11, 12, 13, 25, 27, 28},
		Group:   "outputs",
		Types:   []string{"nc", "out"},
		Modes:   []string{"", "sw", "swlink", "pwm"},
	},
	{
		Numbers: []int{14, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45},
		Group:   "digital",
		Types:   []string{"nc", "in", "out", "dsen", "i2c", "adc"},
		Modes:   []string{"", "p", "pr", "r", "c", "sw", "swlink", "dht11", "dht22", "1w", "1wbus", "ib", "w26", "nc", "sda", "scl", "norm"},
	},
}

// ============================================================
// Здесь подготовка PortMegaD.Ports и PortMegaD.Groups

type portGroup struct {
	Code        string
	Name        string
	Description string
}

type portType struct {
	Code        string
	Name        string
	Description string
	Modes       []*PortMode
}

type PortMode struct {
	Code        string
	Name        string
	Description string
}

type port struct {
	Numbers []int
	Group   string
	Types   []string
	Modes   []string
}

type Port struct {
	Number int
	Group  string
	Types  *orderedmap.OrderedMap[string, *PortType]
}

type PortType struct {
	Code        string
	Name        string
	Description string
	Modes       *orderedmap.OrderedMap[string, *PortMode]
}

func init() {
	// Создаем карту для поиска объекта группы по его коду
	var portGroups = orderedmap.New[string, *portGroup](len(PortGroupList))
	for _, item := range PortGroupList {
		if err := portGroups.Add(item.Code, item); err != nil {
			panic(err)
		}
	}

	// Создаем карту для поиска объектов типа/режима по их кодам
	for _, item := range portTypeList {
		portType := &PortType{
			Code:        item.Code,
			Name:        item.Name,
			Description: item.Description,
			Modes:       orderedmap.New[string, *PortMode](10),
		}

		for _, mode := range item.Modes {
			if err := portType.Modes.Add(mode.Code, mode); err != nil {
				panic(err)
			}
		}

		if err := PortTypes.Add(portType.Code, portType); err != nil {
			panic(err)
		}
	}

	// Создаем карту для поиска объектов портов по их номерам
	for _, item := range portList {
		// Строим карту поддерживаемых режимов
		supportedModes := make(map[string]bool, len(item.Modes))
		for _, modeCode := range item.Modes {
			supportedModes[modeCode] = true
		}

		// Подготавливаем типы
		portTypes := orderedmap.New[string, *PortType](len(item.Types))
		for _, typeCode := range item.Types {
			v, err := PortTypes.Get(typeCode)
			if err != nil {
				panic(err)
			}

			// Copy port type
			portType := *v

			// Некоторые порты поддерживают ШИМ, некоторые нет.
			// В списке есть все режимы для данного типа.
			// Лишние не копируем.
			modes := orderedmap.New[string, *PortMode](portType.Modes.Len())
			for _, mode := range portType.Modes.GetValueList() {
				if supportedModes[mode.Code] {
					if err := modes.Add(mode.Code, mode); err != nil {
						panic(err)
					}
				}
			}
			portType.Modes = modes

			if err := portTypes.Add(typeCode, &portType); err != nil {
				panic(err)
			}
		}

		// Добавляем все порты из списка Numbers в общую карту
		for _, portNumber := range item.Numbers {
			err := Ports.Add(portNumber, &Port{
				Number: portNumber,
				Group:  item.Group,
				Types:  portTypes,
			})
			if err != nil {
				panic(err)
			}
		}
	}

	// Распределяем порты по группам
	for _, port := range Ports.GetValueList() {
		ports, _ := Groups.Get(port.Group)
		Groups.Set(port.Group, append(ports, port))
	}

	// Сортируем по порядку
	for _, portList := range Groups.GetValueList() {
		sort.Slice(portList, func(i, j int) bool {
			return portList[i].Number < portList[j].Number
		})
	}
}
