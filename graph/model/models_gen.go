// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

import (
	"fmt"
	"io"
	"strconv"
)

type EditListenEventCollectionResult interface {
	IsEditListenEventCollectionResult()
}

type EmployeesResult interface {
	IsEmployeesResult()
}

type EventResult interface {
	IsEventResult()
}

type LoginResult interface {
	IsLoginResult()
}

type MeResult interface {
	IsMeResult()
}

type MessagesResult interface {
	IsMessagesResult()
}

type MoveRoomResult interface {
	IsMoveRoomResult()
}

type MutationResult interface {
	IsMutationResult()
}

type ReadMsgResult interface {
	IsReadMsgResult()
}

type RefreshTokensResult interface {
	IsRefreshTokensResult()
}

type RegisterResult interface {
	IsRegisterResult()
}

type RoomsResult interface {
	IsRoomsResult()
}

type SendMsgResult interface {
	IsSendMsgResult()
}

type TagsResult interface {
	IsTagsResult()
}

type AdvancedError struct {
	Code  string `json:"code"`
	Error string `json:"error"`
}

func (AdvancedError) IsMutationResult()                  {}
func (AdvancedError) IsMeResult()                        {}
func (AdvancedError) IsRoomsResult()                     {}
func (AdvancedError) IsTagsResult()                      {}
func (AdvancedError) IsEmployeesResult()                 {}
func (AdvancedError) IsLoginResult()                     {}
func (AdvancedError) IsRefreshTokensResult()             {}
func (AdvancedError) IsRegisterResult()                  {}
func (AdvancedError) IsSendMsgResult()                   {}
func (AdvancedError) IsReadMsgResult()                   {}
func (AdvancedError) IsMoveRoomResult()                  {}
func (AdvancedError) IsMessagesResult()                  {}
func (AdvancedError) IsEditListenEventCollectionResult() {}

type ByCreated struct {
	RoomID   int        `json:"roomID"`
	StartMsg int        `json:"startMsg"`
	Created  MsgCreated `json:"created"`
	Count    int        `json:"count"`
}

type ByRange struct {
	RoomID int `json:"roomID"`
	Start  int `json:"start"`
	End    int `json:"end"`
}

type CreateMessageInput struct {
	RoomID      int    `json:"roomID"`
	TargetMsgID *int   `json:"targetMsgID"`
	Body        string `json:"body"`
}

type DropRoom struct {
	RoomID int `json:"roomID"`
}

func (DropRoom) IsEventResult() {}

type DropTag struct {
	TagID int `json:"tagID"`
}

func (DropTag) IsEventResult() {}

type EmpTagAction struct {
	Action Action `json:"action"`
	EmpID  int    `json:"empID"`
	TagIDs []int  `json:"tagIDs"`
}

func (EmpTagAction) IsEventResult() {}

type Employee struct {
	EmpID     int    `json:"empID"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Tags      *Tags  `json:"tags"`
}

type Employees struct {
	Employees []*Employee `json:"employees"`
}

func (Employees) IsEmployeesResult() {}

type FindEmployees struct {
	EmpID  *int    `json:"empID"`
	RoomID *int    `json:"roomID"`
	TagID  *int    `json:"tagID"`
	Name   *string `json:"name"`
}

type FindMessages struct {
	MsgID        *int    `json:"msgID"`
	EmpID        *int    `json:"empID"`
	RoomID       *int    `json:"roomID"`
	TargetID     *int    `json:"targetID"`
	TextFragment *string `json:"textFragment"`
}

type FindRooms struct {
	RoomID *int    `json:"roomID"`
	Name   *string `json:"name"`
}

type ListenCollection struct {
	SessionKey string          `json:"sessionKey"`
	Success    string          `json:"success"`
	Collection []*ListenedChat `json:"collection"`
}

func (ListenCollection) IsEditListenEventCollectionResult() {}

type ListenedChat struct {
	ID     int         `json:"id"`
	Events []EventType `json:"events"`
}

type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Me struct {
	Employee *Employee     `json:"employee"`
	Personal *PersonalData `json:"personal"`
	Rooms    *Rooms        `json:"rooms"`
}

func (Me) IsMeResult() {}

type Member struct {
	Employee *Employee `json:"employee"`
	Room     *Room     `json:"room"`
}

type MemberAction struct {
	Action  Action `json:"action"`
	EmpID   int    `json:"empID"`
	RoomIDs []int  `json:"roomIDs"`
}

func (MemberAction) IsEventResult() {}

type Members struct {
	Members []*Member `json:"members"`
}

type Message struct {
	Room      *Room     `json:"room"`
	MsgID     int       `json:"msgID"`
	Next      *int      `json:"next"`
	Prev      *int      `json:"prev"`
	Employee  *Employee `json:"employee"`
	TargetMsg *Message  `json:"targetMsg"`
	Body      string    `json:"body"`
	CreatedAt int64     `json:"createdAt"`
}

type Messages struct {
	Messages []*Message `json:"messages"`
}

func (Messages) IsMessagesResult() {}

type NewMessage struct {
	MsgID       int    `json:"msgID"`
	RoomID      int    `json:"roomID"`
	TargetMsgID *int   `json:"targetMsgID"`
	EmployeeID  *int   `json:"employeeID"`
	Body        string `json:"body"`
	CreatedAt   int64  `json:"createdAt"`
	Prev        *int   `json:"prev"`
}

func (NewMessage) IsEventResult() {}

type Params struct {
	Limit  *int `json:"limit"`
	Offset *int `json:"offset"`
}

type PersonalData struct {
	Email       string `json:"email"`
	PhoneNumber string `json:"phoneNumber"`
	Token       string `json:"token"`
}

type Room struct {
	RoomID          int      `json:"roomID"`
	Name            string   `json:"name"`
	View            RoomType `json:"view"`
	PrevRoomID      *int     `json:"prevRoomID"`
	LastMessageRead *int     `json:"lastMessageRead"`
	LastMessageID   *int     `json:"lastMessageID"`
	Members         *Members `json:"members"`
}

type Rooms struct {
	Rooms []*Room `json:"rooms"`
}

func (Rooms) IsRoomsResult() {}

type SubscriptionBody struct {
	Event EventType   `json:"event"`
	Body  EventResult `json:"body"`
}

type Successful struct {
	Success string `json:"success"`
}

func (Successful) IsMutationResult() {}
func (Successful) IsRegisterResult() {}
func (Successful) IsSendMsgResult()  {}
func (Successful) IsReadMsgResult()  {}
func (Successful) IsMoveRoomResult() {}

type Tag struct {
	TagID int    `json:"tagID"`
	Name  string `json:"name"`
}

type Tags struct {
	Tags []*Tag `json:"tags"`
}

func (Tags) IsTagsResult() {}

type TokenExpired struct {
	Message string `json:"message"`
}

func (TokenExpired) IsEventResult() {}

type TokenPair struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

func (TokenPair) IsLoginResult()         {}
func (TokenPair) IsRefreshTokensResult() {}

type Action string

const (
	ActionAdd Action = "ADD"
	ActionDel Action = "DEL"
)

var AllAction = []Action{
	ActionAdd,
	ActionDel,
}

func (e Action) IsValid() bool {
	switch e {
	case ActionAdd, ActionDel:
		return true
	}
	return false
}

func (e Action) String() string {
	return string(e)
}

func (e *Action) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = Action(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid Action", str)
	}
	return nil
}

func (e Action) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type ActionType string

const (
	ActionTypeRead  ActionType = "READ"
	ActionTypeWrite ActionType = "WRITE"
)

var AllActionType = []ActionType{
	ActionTypeRead,
	ActionTypeWrite,
}

func (e ActionType) IsValid() bool {
	switch e {
	case ActionTypeRead, ActionTypeWrite:
		return true
	}
	return false
}

func (e ActionType) String() string {
	return string(e)
}

func (e *ActionType) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = ActionType(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid ActionType", str)
	}
	return nil
}

func (e ActionType) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type EventSubjectAction string

const (
	EventSubjectActionAdd    EventSubjectAction = "ADD"
	EventSubjectActionDelete EventSubjectAction = "DELETE"
)

var AllEventSubjectAction = []EventSubjectAction{
	EventSubjectActionAdd,
	EventSubjectActionDelete,
}

func (e EventSubjectAction) IsValid() bool {
	switch e {
	case EventSubjectActionAdd, EventSubjectActionDelete:
		return true
	}
	return false
}

func (e EventSubjectAction) String() string {
	return string(e)
}

func (e *EventSubjectAction) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = EventSubjectAction(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid EventSubjectAction", str)
	}
	return nil
}

func (e EventSubjectAction) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type EventType string

const (
	EventTypeAll          EventType = "all"
	EventTypeNewMessage   EventType = "NewMessage"
	EventTypeDropTag      EventType = "DropTag"
	EventTypeEmpTagAction EventType = "EmpTagAction"
	EventTypeMemberAction EventType = "MemberAction"
	EventTypeDropRoom     EventType = "DropRoom"
	EventTypeTokenExpired EventType = "TokenExpired"
)

var AllEventType = []EventType{
	EventTypeAll,
	EventTypeNewMessage,
	EventTypeDropTag,
	EventTypeEmpTagAction,
	EventTypeMemberAction,
	EventTypeDropRoom,
	EventTypeTokenExpired,
}

func (e EventType) IsValid() bool {
	switch e {
	case EventTypeAll, EventTypeNewMessage, EventTypeDropTag, EventTypeEmpTagAction, EventTypeMemberAction, EventTypeDropRoom, EventTypeTokenExpired:
		return true
	}
	return false
}

func (e EventType) String() string {
	return string(e)
}

func (e *EventType) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = EventType(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid EventType", str)
	}
	return nil
}

func (e EventType) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type FetchType string

const (
	FetchTypePositive FetchType = "POSITIVE"
	FetchTypeNeutral  FetchType = "NEUTRAL"
	FetchTypeNegative FetchType = "NEGATIVE"
)

var AllFetchType = []FetchType{
	FetchTypePositive,
	FetchTypeNeutral,
	FetchTypeNegative,
}

func (e FetchType) IsValid() bool {
	switch e {
	case FetchTypePositive, FetchTypeNeutral, FetchTypeNegative:
		return true
	}
	return false
}

func (e FetchType) String() string {
	return string(e)
}

func (e *FetchType) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = FetchType(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid FetchType", str)
	}
	return nil
}

func (e FetchType) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type MsgCreated string

const (
	MsgCreatedAfter  MsgCreated = "AFTER"
	MsgCreatedBefore MsgCreated = "BEFORE"
)

var AllMsgCreated = []MsgCreated{
	MsgCreatedAfter,
	MsgCreatedBefore,
}

func (e MsgCreated) IsValid() bool {
	switch e {
	case MsgCreatedAfter, MsgCreatedBefore:
		return true
	}
	return false
}

func (e MsgCreated) String() string {
	return string(e)
}

func (e *MsgCreated) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = MsgCreated(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid MsgCreated", str)
	}
	return nil
}

func (e MsgCreated) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type RoomType string

const (
	RoomTypeTalk RoomType = "TALK"
	RoomTypeBlog RoomType = "BLOG"
)

var AllRoomType = []RoomType{
	RoomTypeTalk,
	RoomTypeBlog,
}

func (e RoomType) IsValid() bool {
	switch e {
	case RoomTypeTalk, RoomTypeBlog:
		return true
	}
	return false
}

func (e RoomType) String() string {
	return string(e)
}

func (e *RoomType) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = RoomType(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid RoomType", str)
	}
	return nil
}

func (e RoomType) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}
