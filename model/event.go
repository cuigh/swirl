package model

import (
	"fmt"
	"time"
)

type EventType string

const (
	// swarm
	EventTypeRegistry     EventType = "Registry"
	EventTypeNode         EventType = "Node"
	EventTypeNetwork      EventType = "Network"
	EventTypeService      EventType = "Service"
	EventTypeStackTask    EventType = "Stack Task"
	EventTypeStackArchive EventType = "Stack Archive"
	EventTypeSecret       EventType = "Secret"
	EventTypeConfig       EventType = "Config"

	// local
	EventTypeVolume EventType = "Volume"

	// system
	EventTypeAuthentication EventType = "Authentication"
	EventTypeRole           EventType = "Role"
	EventTypeUser           EventType = "User"
	EventTypeSetting        EventType = "Setting"
)

type EventAction string

const (
	EventActionLogin      EventAction = "Login"
	EventActionLogout     EventAction = "Logout"
	EventActionCreate     EventAction = "Create"
	EventActionDelete     EventAction = "Delete"
	EventActionUpdate     EventAction = "Update"
	EventActionScale      EventAction = "Scale"
	EventActionDisconnect EventAction = "Disconnect"
)

type Event struct {
	ID       string      `bson:"_id"`
	Type     EventType   `bson:"type"`
	Action   EventAction `bson:"action"`
	Code     string      `bson:"code"`
	Name     string      `bson:"name"`
	UserID   string      `bson:"user_id"`
	Username string      `bson:"username"`
	Time     time.Time   `bson:"time"`
}

func (e *Event) URL(et EventType, code string) string {
	switch et {
	case EventTypeAuthentication:
		return fmt.Sprintf("/system/user/%s/detail", code)
	case EventTypeNode:
		return fmt.Sprintf("/node/%s/detail", code)
	case EventTypeNetwork:
		return fmt.Sprintf("/network/%s/detail", code)
	case EventTypeService:
		return fmt.Sprintf("/service/%s/detail", code)
	case EventTypeStackArchive:
		return fmt.Sprintf("/stack/archive/%s/detail", code)
	case EventTypeVolume:
		return fmt.Sprintf("/volume/%s/detail", code)
	case EventTypeRole:
		return fmt.Sprintf("/system/role/%s/detail", code)
	case EventTypeUser:
		return fmt.Sprintf("/system/user/%s/detail", code)
	case EventTypeSetting:
		return "/system/setting/"
	}
	return ""
}

type EventListArgs struct {
	Type      string `query:"type"`
	Name      string `query:"name"`
	PageIndex int    `query:"page"`
	PageSize  int    `query:"size"`
}
