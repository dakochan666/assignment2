package models

import (
	"time"

	"github.com/asaskevich/govalidator"
	"gorm.io/gorm"
)

type Item struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	ItemCode    string         `gorm:"not null;type:varchar(100)" json:"item_code" valid:"required~item code is required"`
	Description string         `gorm:"not null" json:"description" valid:"required~description is required"`
	Quantity    int            `gorm:"not null" json:"quantity" valid:"required~quantity is required"`
	OrderID     uint           `json:"order_id,omitempty"`
	CreatedAt   *time.Time     `json:"created_at,omitempty"`
	UpdatedAt   *time.Time     `json:"updated_at,omitempty"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

func (i *Item) BeforeCreate(tx *gorm.DB) (err error) {
	_, errCreate := govalidator.ValidateStruct(i)
	if errCreate != nil {
		return errCreate
	}

	return
}
