package models

import (
	"time"

	"github.com/asaskevich/govalidator"
	"gorm.io/gorm"
)

type Order struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	CustomerName string         `gorm:"not null;type:varchar(100)" json:"customer_name" valid:"required~customer name is required"`
	OrderedAt    *time.Time     `gorm:"autoCreateTime" json:"ordered_at"`
	Items        []Item         `json:"items,omtiempty"`
	CreatedAt    *time.Time     `json:"created_at,omitempty"`
	UpdatedAt    *time.Time     `json:"updated_at,omitempty"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

func (o *Order) BeforeCreate(tx *gorm.DB) (err error) {
	_, errCreate := govalidator.ValidateStruct(o)
	if errCreate != nil {
		return errCreate
	}

	return
}
