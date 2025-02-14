package objects

import (
	"encoding/json"
	"sort"

	"github.com/pkg/errors"
	"touchon-server/internal/g"
	"touchon-server/internal/model"
	"touchon-server/internal/store"
	"touchon-server/lib/helpers"
	"touchon-server/lib/interfaces"
)

// Implementation of Object interface

func NewObjectModelImpl(category model.Category, objType string, internal bool, name string, props []*Prop, children []Object, events []interfaces.Event, methods []*Method, tags []string) (*ObjectModelImpl, error) {
	o := &ObjectModelImpl{
		category: category,
		objType:  objType,
		internal: internal,
		name:     name,
		status:   model.StatusNA,
		props:    NewProps(),
		children: NewChildren(),
		events:   NewEvents(),
		methods:  NewMethods(),
		tags:     make(map[string]bool, len(tags)),
		enabled:  true,
	}

	if err := o.GetProps().Add(props...); err != nil {
		return nil, errors.Wrap(err, "NewObjectModelImpl")
	}

	o.GetChildren().Add(children...)
	o.GetMethods().Add(methods...)

	if err := o.GetEvents().Add(events...); err != nil {
		return nil, errors.Wrap(err, "NewObjectModelImpl")
	}

	o.SetTags(tags...)

	return o, nil
}

type ObjectModelImpl struct {
	id       int
	parentID *int
	zoneID   *int

	category model.Category
	objType  string
	internal bool
	name     string
	status   model.ObjectStatus

	props    *Props
	children *Children
	events   *Events
	methods  *Methods
	tags     map[string]bool
	enabled  bool

	msgHandlerIDs []int
}

func (o *ObjectModelImpl) GetID() int {
	return o.id
}

func (o *ObjectModelImpl) SetID(id int) {
	o.id = id
}

func (o *ObjectModelImpl) GetParentID() *int {
	return o.parentID
}

func (o *ObjectModelImpl) SetParentID(parentID *int) {
	o.parentID = parentID
}

func (o *ObjectModelImpl) GetZoneID() *int {
	return o.zoneID
}

func (o *ObjectModelImpl) SetZoneID(zoneID *int) {
	o.zoneID = zoneID
}

func (o *ObjectModelImpl) GetCategory() model.Category {
	return o.category
}

func (o *ObjectModelImpl) SetCategory(v model.Category) {
	o.category = v
}

func (o *ObjectModelImpl) GetType() string {
	return o.objType
}

func (o *ObjectModelImpl) SetType(v string) {
	o.objType = v
}

func (o *ObjectModelImpl) GetInternal() bool {
	return o.internal
}

func (o *ObjectModelImpl) SetInternal(v bool) {
	o.internal = v
}

func (o *ObjectModelImpl) GetName() string {
	return o.name
}

func (o *ObjectModelImpl) SetName(v string) {
	o.name = v
}

func (o *ObjectModelImpl) GetStatus() model.ObjectStatus {
	return o.status
}

func (o *ObjectModelImpl) SetStatus(v model.ObjectStatus) {
	o.status = v
}

func (o *ObjectModelImpl) GetProps() *Props {
	return o.props
}

func (o *ObjectModelImpl) GetChildren() *Children {
	return o.children
}

func (o *ObjectModelImpl) MarshalJSON() ([]byte, error) {
	props := o.GetProps()
	if props.Len() == 0 {
		props = nil
	}

	children := o.GetChildren()
	if children.Len() == 0 {
		children = nil
	}

	events := o.GetEvents()
	if events.Len() == 0 {
		events = nil
	}

	methods := o.GetMethods()
	if methods.Len() == 0 {
		methods = nil
	}

	return json.Marshal(&ObjectModel{
		ID:       o.GetID(),
		ParentID: o.GetParentID(),
		ZoneID:   o.GetZoneID(),
		Category: o.GetCategory(),
		Type:     o.GetType(),
		Internal: o.GetInternal(),
		Name:     o.GetName(),
		Status:   o.GetStatus(),
		Props:    props,
		Children: children,
		Events:   events,
		Methods:  methods,
		Tags:     o.GetTags(),
		Enabled:  o.GetEnabled(),
	})
}

func (o *ObjectModelImpl) UnmarshalJSON(data []byte) error {
	v := &ObjectModel{
		Props:    o.props,
		Children: o.children,
		// Нельзя переопределять с фронта
		// Events:   o.events,
		// Methods:  o.methods,
	}

	if err := json.Unmarshal(data, v); err != nil {
		return errors.Wrap(err, "ObjectModelImpl.UnmarshalJSON")
	}

	switch {
	case o.GetCategory() != v.Category:
		return errors.Wrap(errors.Errorf("can't unmarshal object of category %s to object of category %s", v.Category, o.GetCategory()), "ObjectModelImpl.UnmarshalJSON")
	case o.GetType() != v.Type:
		return errors.Wrap(errors.Errorf("can't unmarshal object %s/%s to object %s/%s", v.Category, v.Type, o.GetCategory(), o.GetType()), "ObjectModelImpl.UnmarshalJSON")
	case o.GetID() != v.ID:
		return errors.Wrap(errors.Errorf("can't unmarshal object %s/%s %d to object %s/%s %d", v.Category, v.Type, v.ID, o.GetCategory(), o.GetType(), o.GetID()), "ObjectModelImpl.UnmarshalJSON")
	}

	o.SetID(v.ID)
	o.SetParentID(v.ParentID)
	o.SetZoneID(v.ZoneID)
	// o.SetCategory(v.Category) // Нельзя переопределять с фронта
	// o.SetType(v.Type)         // Нельзя переопределять с фронта
	// o.SetInternal(v.Internal) // Нельзя переопределять с фронта
	o.SetName(v.Name)
	// o.SetStatus(v.Status)     // Нельзя переопределять с фронта
	o.SetTags(v.Tags...)
	o.SetEnabled(v.Enabled)

	return nil
}

func (o *ObjectModelImpl) GetEvents() *Events {
	return o.events
}

func (o *ObjectModelImpl) GetMethods() *Methods {
	return o.methods
}

func (o *ObjectModelImpl) Init(storeObj *model.StoreObject, childType model.ChildType) error {
	o.SetID(storeObj.ID)
	o.SetParentID(storeObj.ParentID)
	o.SetZoneID(storeObj.ZoneID)
	o.SetCategory(storeObj.Category)
	o.SetType(storeObj.Type)
	o.SetInternal(storeObj.Internal)
	o.SetName(storeObj.Name)
	o.SetStatus(storeObj.Status)
	o.SetTagsMap(storeObj.Tags)
	o.SetEnabled(storeObj.Enabled)

	// Загружаем свойства объекта
	props, err := store.I.ObjectRepository().GetProps(storeObj.ID)
	if err != nil {
		return errors.Wrapf(err, "ObjectModelImpl.Init (%d)", storeObj.ID)
	}

	for _, prop := range props {
		if err := o.GetProps().Set(prop.Code, prop.Value); err != nil {
			return errors.Wrapf(err, "ObjectModelImpl.Init (%d)", storeObj.ID)
		}
	}

	// Очищаем список детей перед загрузкой из базы
	o.GetChildren().DeleteAll()

	if childType == model.ChildTypeNobody {
		return nil
	}

	// Загружаем дочерние объекты
	children, err := store.I.ObjectRepository().GetObjectChildren(childType, storeObj.ID)
	if err != nil {
		return errors.Wrapf(err, "ObjectModelImpl.Init (%d)", storeObj.ID)
	}

	for _, childStoreObj := range children {
		childObjModel, err := LoadObject(childStoreObj.ID, "", "", childType)
		if err != nil {
			return errors.Wrapf(err, "ObjectModelImpl.Init (%d)", storeObj.ID)
		}

		o.GetChildren().Add(childObjModel)
	}

	return nil
}

func (o *ObjectModelImpl) Save() error {
	// Сохраняем поля объекта
	storeObj := o.GetStoreObject()

	if err := store.I.ObjectRepository().SaveObject(storeObj); err != nil {
		return errors.Wrapf(err, "ObjectModelImpl.Save(%s/%s)", o.category, o.objType)
	}

	o.SetID(storeObj.ID)
	for _, child := range o.GetChildren().items {
		child.SetParentID(&storeObj.ID)
	}

	// Сохраняем свойства объекта
	props := make(map[string]string)
	for _, prop := range o.GetProps().m.GetValueList() {
		props[prop.Code] = prop.StringValue()
	}

	if err := store.I.ObjectRepository().SetProps(o.GetID(), props); err != nil {
		return errors.Wrapf(err, "ObjectModelImpl.Save(%s/%s)", o.category, o.objType)
	}

	// Сохраняем детей
	for _, child := range o.GetChildren().items {
		if err := child.Save(); err != nil {
			return errors.Wrapf(err, "ObjectModelImpl.Save(%s/%s)", o.category, o.objType)
		}
	}

	return nil
}

func (o *ObjectModelImpl) Subscribe(msgType interfaces.MessageType, name string, targetType interfaces.TargetType, targetID *int, handler interfaces.MsgHandler) error {
	handlerID, err := g.Msgs.Subscribe(msgType, name, targetType, targetID, handler)
	if err != nil {
		return errors.Wrap(err, "ObjectModelImpl.Subscribe")
	}

	o.msgHandlerIDs = append(o.msgHandlerIDs, handlerID)

	return nil
}

func (o *ObjectModelImpl) CheckEnabled() error {
	if !o.enabled {
		return ErrObjectDisabled
	}

	return nil
}

func (o *ObjectModelImpl) Start() error {
	if err := o.CheckEnabled(); err != nil {
		return errors.Wrap(err, "ObjectModelImpl.Start")
	}

	// Подписываем объект на обработку команд (вызов методов)
	err := o.Subscribe(interfaces.MessageTypeCommand, "", interfaces.TargetTypeObject, &o.id, o.commandHandler)
	if err != nil {
		return errors.Wrap(err, "ObjectModelImpl.Start")
	}

	return nil
}

func (o *ObjectModelImpl) commandHandler(svc interfaces.MessageSender, msg interfaces.Message) {
	cmd, ok := msg.(interfaces.Command)
	if !ok {
		g.Logger.Error(errors.Wrap(errors.Errorf("msg is not command: %T", msg), "ObjectModelImpl.commandHandler"))
		return
	}

	method, err := o.GetMethods().Get(cmd.GetName())
	if err != nil {
		g.Logger.Error(errors.Wrap(err, "ObjectModelImpl.commandHandler"))
		return
	}

	msgsList, err := method.Func(cmd.GetArgs())
	if err != nil {
		g.Logger.Error(errors.Wrap(err, "ObjectModelImpl.commandHandler"))
		return
	}

	if err := svc.Send(msgsList...); err != nil {
		g.Logger.Error(errors.Wrap(err, "ObjectModelImpl.commandHandler"))
	}
}

func (o *ObjectModelImpl) Shutdown() error {
	if err := o.CheckEnabled(); err != nil {
		return errors.Wrap(err, "ObjectModelImpl.Shutdown")
	}

	g.Msgs.Unsubscribe(o.msgHandlerIDs...)

	return nil
}

func (o *ObjectModelImpl) GetStoreObject() *model.StoreObject {
	return &model.StoreObject{
		ID:       o.id,
		ParentID: o.parentID,
		ZoneID:   o.zoneID,
		Category: o.category,
		Type:     o.objType,
		Internal: o.internal,
		Name:     o.name,
		Status:   o.status,
		Tags:     o.GetTagsMap(),
		Enabled:  o.enabled,
	}
}

func (o *ObjectModelImpl) GetTags() []string {
	r := make([]string, 0, len(o.tags))

	for tag := range o.tags {
		r = append(r, tag)
	}

	sort.Strings(r)

	return r
}

func (o *ObjectModelImpl) ReplaceTags(tags ...string) {
	for tag := range o.tags {
		delete(o.tags, tag)
	}

	o.SetTags(tags...)
}

func (o *ObjectModelImpl) SetTags(tags ...string) {
	for _, tag := range tags {
		o.tags[helpers.PrepareTag(tag)] = true
	}
}

func (o *ObjectModelImpl) DeleteTags(tags ...string) {
	for _, tag := range tags {
		delete(o.tags, helpers.PrepareTag(tag))
	}
}

func (o *ObjectModelImpl) GetTagsMap() map[string]bool {
	return o.tags
}

func (o *ObjectModelImpl) SetTagsMap(tags map[string]bool) {
	for tag := range o.tags {
		delete(o.tags, tag)
	}

	for tag := range tags {
		o.SetTags(tag)
	}
}

func (o *ObjectModelImpl) GetEnabled() bool {
	return o.enabled
}

func (o *ObjectModelImpl) SetEnabled(v bool) {
	o.enabled = v
}

func (o *ObjectModelImpl) DeleteChildren() error {
	for _, child := range o.GetChildren().GetAll() {
		if err := child.DeleteChildren(); err != nil {
			return errors.Wrap(err, "DeleteChildren")
		}
	}

	return nil
}
