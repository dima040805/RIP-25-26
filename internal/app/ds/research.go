package ds

import (
	"time"
)

type Research struct {
	ID          uint           `gorm:"primaryKey"`
	Status      string         `gorm:"type:varchar(15);not null"`
	DateResearch time.Time 		`gorm:"not null"`
	DateCreate  time.Time      `gorm:"not null"`
	CreatorID   int         `gorm:"not null"`

	Creator   Users `gorm:"foreignKey:CreatorID"`
}