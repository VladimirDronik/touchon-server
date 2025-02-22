package objects

import (
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"touchon-server/internal/model"
	"touchon-server/internal/store"
	"touchon-server/lib/interfaces"
)

// LoadObject создает модель объекта и заполняет его данными из БД
func LoadObject(objectID int, objCat model.Category, objType string, childType model.ChildType) (Object, error) {
	// Если создаем новый объект, то создаем модель и возвращаем ее сразу
	if objectID <= 0 {
		objModel, err := GetObjectModel(objCat, objType)
		if err != nil {
			return nil, errors.Wrap(err, "LoadObject")
		}

		return objModel, nil
	}

	// Если объект уже существует - создаем модель и заполняем данными из БД
	storeObj, err := store.I.ObjectRepository().GetObject(objectID)
	if err != nil {
		return nil, errors.Wrap(err, "LoadObject")
	}

	objModel, err := GetObjectModel(storeObj.Category, storeObj.Type)
	if err != nil {
		return nil, errors.Wrap(err, "LoadObject")
	}

	if err := objModel.Init(storeObj, childType); err != nil {
		return nil, errors.Wrap(err, "LoadObject")
	}

	return objModel, nil
}

func LoadPort(objectID int, childType model.ChildType) (interfaces.Port, error) {
	portObj, err := LoadObject(objectID, model.CategoryPort, "port_mega_d", childType)
	if err != nil {
		return nil, errors.Wrap(err, "LoadPort")
	}

	port, ok := portObj.(interfaces.Port)
	if !ok {
		return nil, errors.Wrap(errors.Errorf("object %T not implements interface Port", portObj), "LoadPort")
	}

	return port, nil
}

func Checks(checks ...PropValueChecker) PropValueChecker {
	return func(prop *Prop, allProps map[string]*Prop) error {
		for _, check := range checks {
			if err := check(prop, allProps); err != nil {
				return errors.Wrap(err, "CheckAll")
			}
		}

		return nil
	}
}

func AboveOrEqual1() PropValueChecker {
	return func(prop *Prop, allProps map[string]*Prop) error {
		if v, err := prop.GetFloatValue(); err == nil && v < 1 {
			return errors.Wrap(errors.Errorf("%q < 1", prop.Name), "AboveOrEqual1")
		}

		return nil
	}
}

func BelowOrEqual(prop2Code string) PropValueChecker {
	return func(prop *Prop, allProps map[string]*Prop) error {
		prop2, ok := allProps[prop2Code]
		if !ok {
			return errors.Wrap(errors.Errorf("prop %q not found", prop2Code), "BelowOrEqual")
		}

		v1, err1 := prop.GetFloatValue()
		v2, err2 := prop2.GetFloatValue()

		// Если оба значения заданы
		if err1 == nil && err2 == nil && v1 > v2 {
			return errors.Wrap(errors.Errorf("%q > %q", prop.Name, prop2.Name), "BelowOrEqual")
		}

		return nil
	}
}

func Between(minValue, maxValue float32) PropValueChecker {
	return func(prop *Prop, allProps map[string]*Prop) error {
		if v, err := prop.GetFloatValue(); err == nil && (v < minValue || maxValue < v) {
			return errors.Wrap(errors.Errorf("value %v of %s is not in [%v, %v]", v, prop.Name, minValue, maxValue), "Between")
		}

		return nil
	}
}

func Above(minLimit float32) PropValueChecker {
	return func(prop *Prop, allProps map[string]*Prop) error {
		if v, err := prop.GetFloatValue(); err == nil && v <= minLimit {
			return errors.Wrap(errors.Errorf("value %v of %s is below or equal %v", v, prop.Name, minLimit), "Above")
		}

		return nil
	}
}

func Below(maxLimit float32) PropValueChecker {
	return func(prop *Prop, allProps map[string]*Prop) error {
		if v, err := prop.GetFloatValue(); err == nil && maxLimit <= v {
			return errors.Wrap(errors.Errorf("value %v of %s is above or equal %v", v, prop.Name, maxLimit), "Below")
		}

		return nil
	}
}

// ConfigureDevice настройка устройства, на котором находится объект
func ConfigureDevice(interfaceConnection string, addressObject string, options map[string]string, title string) error {
	var port [2]int
	var modePt [2]string
	var typePt [2]string
	params := make(map[int]map[string]string)
	params[0] = make(map[string]string)
	params[1] = make(map[string]string)

	switch interfaceConnection {
	case "NC":
		ports := strings.Split(addressObject, ";")
		port[0], _ = strconv.Atoi(ports[0])
		typePt[0] = "nc"
		modePt[0] = ""
		if len(ports) > 1 {
			port[1], _ = strconv.Atoi(ports[1])
			modePt[1] = ""
			typePt[1] = "nc"
		}
	case "MEGA-IN":
		ports := strings.Split(addressObject, ";")
		port[0], _ = strconv.Atoi(ports[0])
		typePt[0] = "in"
		if len(ports) > 1 {
			port[1], _ = strconv.Atoi(ports[1])
			typePt[1] = "in"
		}
		modePt[0] = options["mode"]
	case "DISCRETE":
		port[0], _ = strconv.Atoi(addressObject)
		typePt[0] = "adc"
		modePt[0] = "norm"
	case "MEGA-OUT":
		port[0], _ = strconv.Atoi(addressObject)
		typePt[0] = "out"
		modePt[0] = "sw"
	case "1W":
		port[0], _ = strconv.Atoi(addressObject)
		typePt[0] = "dsen"
		modePt[0] = "1w"
	case "1WBUS":
		port[0], _ = strconv.Atoi(strings.Split(addressObject, ";")[0])
		typePt[0] = "dsen"
		modePt[0] = "1wbus"
	case "I2C":
		ports := strings.Split(addressObject, ";")
		if len(ports) > 0 {
			port[0], _ = strconv.Atoi(ports[0])
		}
		if len(ports) > 1 {
			port[1], _ = strconv.Atoi(ports[1])
		}
		typePt[0] = "i2c"
		typePt[1] = "i2c"
		modePt[0] = "sda"
		modePt[1] = "scl"

		if portSCLObject, err := getObjects(port[1], "", ""); err == nil && len(portSCLObject) > 0 {
			if portSCL, err := portSCLObject[0].GetProps().GetIntValue("number"); err == nil {
				params[0]["misc"] = strconv.Itoa(portSCL) //указываем порт, на котором находится SCL
			}
		}
		params[0]["gr"] = options["gr"]
		params[0]["d"] = options["d"]
	}

	//Конфигурим порт контроллера
	for k, p := range port {
		if p != 0 {
			portObj, err := LoadPort(p, model.ChildTypeNobody)
			if err != nil {
				return errors.Wrap(err, "getValues")
			}

			err = portObj.SetTypeMode(typePt[k], modePt[k], title, params[k])
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// FillOptions извлечение из свойств объекта опции, которые нужны для настройки портов контроллера или другого устройства
func FillOptions(typeObject string, props map[string]interface{}) (map[string]string, error) {
	paramsI2C := map[string]map[string]string{
		"htu21d": {"gr": "1", "d": "1"},
		"htu31d": {"gr": "1", "d": "56"},
		"bme280": {"gr": "1", "d": "6"},
		"bmp280": {"gr": "1", "d": "6"},
		"bh1750": {"gr": "2", "d": "2"},
		"scd4x":  {"gr": "5", "d": "44"},
	}

	//Если тип объекта является датчиком I2C
	if paramsI2C[typeObject] != nil {
		return paramsI2C[typeObject], nil
	}

	switch typeObject {
	case "generic_input":
		if props["mode"] == nil {
			return nil, errors.New("no input parameters found: mode")
		}
		return map[string]string{"mode": props["mode"].(string)}, nil
	case "presence":
	case "motion":
		return map[string]string{"mode": "P"}, nil
	}

	return nil, nil
}

// FindRelatedObjects Ищет объекты по адресу размещения и распределяет по группам
// objResetGroup - в группе находятся объекты, которые не могут находится на одних портах с исходным
// objRelatedGroup - в группе находятся объекты, которые связаны с текущим и могут существовать на одних портах с исходным
func FindRelatedObjects(addressObject string, typeInterface string, objectID int, objectType string) (map[int]string, map[int]string, error) {
	var objResetGroup = make(map[int]string)
	var objRelatedGroup = make(map[int]string)

	if addressObject == "0" {
		return nil, nil, nil
	}

	ports := strings.Split(addressObject, ";")

	storageObjects, err := store.I.ObjectRepository().GetObjectsByAddress(ports)
	if err != nil {
		return nil, nil, errors.Wrap(err, "FindRelatedObjects")
	}

	for _, storageObject := range storageObjects {
		storageTypeInterface, err := store.I.ObjectRepository().GetProp(storageObject.ID, "interface")
		if err != nil {
			continue
		}

		addressStorageObject, err := store.I.ObjectRepository().GetProp(storageObject.ID, "address")
		if err != nil {
			continue
		}

		storageObjectPorts := strings.Split(addressStorageObject, ";")

		//Если новый объект не I2C или старый объект не I2C, то старый объект на удаление
		if typeInterface != "I2C" || storageTypeInterface != "I2C" {
			objResetGroup[storageObject.ID] = addressStorageObject
			continue
		}

		//Проверяем на полное совпадение адреса объектов I2C
		if addressObject == addressStorageObject {
			//Проверяем тип старого объекта, если совпадает с типом нового, то на удаление, т.к. на одном порту не может быть два одинаковых
			if storageObject.Type == objectType {
				objResetGroup[storageObject.ID] = addressStorageObject
			} else {
				objRelatedGroup[storageObject.ID] = addressStorageObject
			}
			continue
		}

		//Ищем есть ли объекты, у которых SCL совпадает с SCL искомого
		if len(ports) > 1 && len(storageObjectPorts) > 1 && ports[1] == storageObjectPorts[1] {
			objRelatedGroup[storageObject.ID] = addressStorageObject
			continue
		}

		//Если объект уже был в БД и у него просто меняем SDA и SCL местами
		if storageObject.ID == objectID {
			if len(ports) > 1 && len(storageObjectPorts) > 1 && storageObjectPorts[0] != ports[1] && storageObjectPorts[1] != ports[0] {
				objResetGroup[storageObject.ID] = addressStorageObject
			}

			continue
		}

		objResetGroup[storageObject.ID] = addressStorageObject
	}

	delete(objResetGroup, objectID)
	delete(objRelatedGroup, objectID)

	return objResetGroup, objRelatedGroup, nil
}

func ResetPortToDefault(objectsToReset map[int]string, relatedObjects map[int]string) {
	var portsToReset = make(map[string]bool)

	for _, addressToReset := range objectsToReset {
		ports := strings.Split(addressToReset, ";")
		portsToReset[ports[0]] = true
		if len(ports) > 1 {
			portsToReset[ports[1]] = true
		}
	}

	for _, addressToSafe := range relatedObjects {
		ports := strings.Split(addressToSafe, ";")
		delete(portsToReset, ports[0])
		if len(ports) > 1 {
			delete(portsToReset, ports[1])
		}
	}

	for portToReset, _ := range portsToReset {
		err := ConfigureDevice("NC", portToReset, nil, "")
		if err != nil {
			//TODO: тут сформировать запись в лог, что не могли изменить состояние порта на дефолтное и убрать вывод ошибки, чтобы неуспешность действия не было критичным
			//return nil, http.StatusBadRequest, err
		}
	}
}
