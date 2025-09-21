package ds

import (
	"time"
)


type Research struct {
	ID           int       `gorm:"primaryKey;autoIncrement"`
	Status       string    `gorm:"type:varchar(20);not null"` 
	DateResearch time.Time `gorm:"default:null"`                 
	DateCreate   time.Time `gorm:"not null"`                  
	CreatorID    int       `gorm:"not null"`                 
	DateForm     time.Time `gorm:"default:null"`                  
	DateFinish   time.Time `gorm:"default:null"`                 
	ModeratorID  int       `gorm:"default:null"`                  
	
	Creator    User `gorm:"foreignKey:CreatorID"`
	Moderator  User `gorm:"foreignKey:ModeratorID"`
}