package objects

import (
	"encoding/json"
	"sort"
	"time"

	"github.com/pkg/errors"
	"touchon-server/internal/g"
	"touchon-server/internal/model"
	"touchon-server/internal/store"
	"touchon-server/lib/helpers"
	"touchon-server/lib/interfaces"
	"touchon-server/lib/messages"
)

// Implementation of Object interface

func NewObjectModelImpl(category model.Category, objType string, flags Flags, name string, props []*Prop, children []Object, events []interfaces.Event, methods []*Method, tags []string) (*ObjectModelImpl, error) {
	o := &ObjectModelImpl{
		category: category,
		objType:  objType,
		name:     name,
		status:   model.StatusNA,
		props:    NewProps(),
		children: NewChildren(),
		events:   NewEvents(),
		methods:  NewMethods(),
		tags:     make(map[string]bool, len(tags)),
		enabled:  true,
		flags:    flags,
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
	name     string
	status   model.ObjectStatus
	props    *Props
	children *Children
	events   *Events
	methods  *Methods
	tags     map[string]bool
	enabled  bool
	flags    Flags

	msgHandlerIDs []int

	// Используется для выполнения периодических действий
	intervalTimer *helpers.Timer
}

func (o *ObjectModelImpl) GetID() int {
	return o.id
}

func (o *ObjectModelImpl) SetID(id int) {
	o.id = id
}

func (o *ObjectModelImpl) GetParentID() *int {
	if o.parentID == nil {
		return nil
	}

	// copy value
	v := *o.parentID

	return &v
}

func (o *ObjectModelImpl) SetParentID(parentID *int) {
	if parentID == nil {
		o.parentID = nil
		return
	}

	// copy value
	v := *parentID

	o.parentID = &v
}

func (o *ObjectModelImpl) GetZoneID() *int {
	if o.zoneID == nil {
		return nil
	}

	// copy value
	v := *o.zoneID

	return &v
}

func (o *ObjectModelImpl) SetZoneID(zoneID *int) {
	if zoneID == nil {
		o.zoneID = nil
		return
	}

	// copy value
	v := *zoneID

	o.zoneID = &v
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
	// Копируем состояние объекта в методы, для их отключения при необходимости
	o.methods.SetEnabled(o.enabled)

	return nil
}

func (o *ObjectModelImpl) GetEvents() *Events {
	return o.events
}

func (o *ObjectModelImpl) GetMethods() *Methods {
	return o.methods
}

func (o *ObjectModelImpl) Init(storeObj *model.StoreObject, withChildren bool) error {
	o.SetID(storeObj.ID)
	o.SetParentID(storeObj.ParentID)
	o.SetZoneID(storeObj.ZoneID)
	o.SetCategory(storeObj.Category)
	o.SetType(storeObj.Type)
	o.SetName(storeObj.Name)
	o.SetStatus(storeObj.Status)
	o.setTagsMap(storeObj.Tags)
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

	if !withChildren {
		return nil
	}

	// Загружаем дочерние объекты
	children, err := store.I.ObjectRepository().GetObjectChildren(storeObj.ID)
	if err != nil {
		return errors.Wrapf(err, "ObjectModelImpl.Init (%d)", storeObj.ID)
	}

	for _, childStoreObj := range children {
		childObjModel, err := LoadObject(childStoreObj.ID, "", "", withChildren)
		if err != nil {
			return errors.Wrapf(err, "ObjectModelImpl.Init (%d)", storeObj.ID)
		}

		o.GetChildren().Add(childObjModel)
	}

	return nil
}

func (o *ObjectModelImpl) Save() error {
	// Сохраняем поля объекта
	storeObj := o.getStoreObject()

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

	g.Logger.Debugf("%s/%s (%d, %q) starting..", o.GetCategory(), o.GetType(), o.GetID(), o.GetName())

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

	g.Logger.Debugf("%s/%s (%d, %q) stopping..", o.GetCategory(), o.GetType(), o.GetID(), o.GetName())

	g.Msgs.Unsubscribe(o.msgHandlerIDs...)

	if o.intervalTimer != nil {
		o.intervalTimer.Stop()
	}

	return nil
}

func (o *ObjectModelImpl) getStoreObject() *model.StoreObject {
	return &model.StoreObject{
		ID:       o.id,
		ParentID: o.parentID,
		ZoneID:   o.zoneID,
		Category: o.category,
		Type:     o.objType,
		Name:     o.name,
		Status:   o.status,
		Tags:     o.getTagsMap(),
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

func (o *ObjectModelImpl) getTagsMap() map[string]bool {
	return o.tags
}

func (o *ObjectModelImpl) setTagsMap(tags map[string]bool) {
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
	o.methods.SetEnabled(v)
}

func (o *ObjectModelImpl) DeleteChildren() error {
	for _, child := range o.GetChildren().GetAll() {
		if err := child.DeleteChildren(); err != nil {
			return errors.Wrap(err, "DeleteChildren")
		}
	}

	return nil
}

func (o *ObjectModelImpl) SetTimer(interval time.Duration, handler func()) {
	if o.intervalTimer != nil {
		o.intervalTimer.Stop()
	}

	o.intervalTimer = helpers.NewTimer(interval, handler)
}

func (o *ObjectModelImpl) GetTimer() *helpers.Timer {
	return o.intervalTimer
}

func (o *ObjectModelImpl) GetState() (interfaces.Message, error) {
	msg, err := messages.NewEvent("on_get_state", interfaces.TargetTypeObject, o.GetID())
	if err != nil {
		return nil, errors.Wrap(err, "ObjectModelImpl.GetState")
	}

	for _, p := range o.GetProps().GetAll().GetValueList() {
		if p.Visible.Check(o.GetProps()) {
			msg.SetValue(p.Code, p.GetValue())
		}
	}

	return msg, nil
}

func (o *ObjectModelImpl) GetFlags() Flags {
	return o.flags
}

func (o *ObjectModelImpl) SetFlags(v Flags) {
	o.flags = v
}
