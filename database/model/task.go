package model

import (
	"time"
)

type EnumStatus string

var (
	EnumStatusNew        EnumStatus = "new"
	EnumStatusInProgress EnumStatus = "in_progress"
	EnumStatusDone       EnumStatus = "done"
)

type Task struct {
	Id          int32      `json:"id"`
	Title       string     `json:"title"`
	Description *string    `json:"description"`
	Status      EnumStatus `json:"status"`
	CreatedAt   *time.Time `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
}
