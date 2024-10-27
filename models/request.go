package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Request struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	Type      string    `gorm:"type:varchar(20);not null"`
	Amount    float64   `gorm:"type:numeric(18,2);not null"`
	Recipient string    `gorm:"type:varchar(255);not null"`
	Status    string    `gorm:"type:varchar(20);"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	RunID     string    `gorm:"type:varchar(20)"`
}

func (r *Request) BeforeCreate(tx *gorm.DB) (err error) {
	r.ID = uuid.New()
	return
}
