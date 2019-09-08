package entity

import (
	"time"

	"github.com/satori/go.uuid"
)

const (
	GroupTable GroupTableConfig = iota
	GroupId
	GroupUuid
	GroupName
	GroupCreatedAt
	GroupUpdatedAt
)

type Group struct {
	Id        int       `gorm:"primary_key"json:"id"`
	Uuid      uuid.UUID `gorm:"type:varchar(128)"json:"uuid"`
	Name      string    `gorm:"unique;type:varchar(128)"validate:"required"json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type GroupTableConfig int

func (gc GroupTableConfig) String() string {
	switch gc {
	case GroupTable:
		return "groups"
	case GroupId:
		return "id"
	case GroupUuid:
		return "uuid"
	case GroupName:
		return "name"
	case GroupCreatedAt:
		return "created_at"
	case GroupUpdatedAt:
		return "updated_at"
	}
	panic("Unknown value")
}