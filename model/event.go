package model

import (
	"fmt"
	"time"
)

type EventType string

const (
	EventTypeRegistry        EventType = "Registry"
	EventTypeNode            EventType = "Node"
	EventTypeNetwork         EventType = "Network"
	EventTypeService         EventType = "Service"
	EventTypeServiceTemplate EventType = "Service Template"
	EventTypeStack           EventType = "Stack"
	EventTypeSecret          EventType = "Secret"
	EventTypeConfig          EventType = "Config"
	EventTypeVolume          EventType = "Volume"
	EventTypeAuthentication  EventType = "Authentication"
	EventTypeRole            EventType = "Role"
	EventTypeUser            EventType = "User"
	EventTypeSetting         EventType = "Setting"
)

type EventAction string

const (
	EventActionLogin EventAction = "Login"
	//EventActionLogout     EventAction = "Logout"
	EventActionCreate     EventAction = "Create"
	EventActionDelete     EventAction = "Delete"
	EventActionUpdate     EventAction = "Update"
	EventActionScale      EventAction = "Scale"
	EventActionRollback   EventAction = "Rollback"
	EventActionRestart    EventAction = "Restart"
	EventActionDisconnect EventAction = "Disconnect"
	EventActionDeploy     EventAction = "Deploy"
	EventActionShutdown   EventAction = "Shutdown"
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
	case EventTypeStack:
		return fmt.Sprintf("/stack/%s/detail", code)
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
	Type      string `bind:"type"`
	Name      string `bind:"name"`
	PageIndex int    `bind:"page"`
	PageSize  int    `bind:"size"`
}
