package ds

import "github.com/google/uuid"

type User struct {
	ID          uuid.UUID    `gorm:"primaryKey;autoIncrement"`
	Login       string `gorm:"type:varchar(50);not null;unique"`
	Password    string `gorm:"type:varchar(100);not null"`
	IsModerator bool   `gorm:"not null;default:false"`
}