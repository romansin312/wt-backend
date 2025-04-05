package models

import "github.com/google/uuid"

type ActionType int

const (
	Play ActionType = iota
	Pause
	Progress
	UserConnected
	UserDisconnected
)

type ActionMessage struct {
	SenderUserId int32
	Timestamp    int64
	ActionType   ActionType
	ActionInfo   string
	RoomId       uuid.UUID
}
