package ds

import (
		"database/sql"
	"time"
	"github.com/google/uuid"

)


type Research struct {
	ID           int       `gorm:"primaryKey;autoIncrement"`
	Status       string    `gorm:"type:varchar(20);not null"` 
	DateResearch string `gorm:"default:null"`                 
	DateCreate   time.Time `gorm:"not null"`                  
	CreatorID    uuid.UUID       `gorm:"not null"`                 
	DateForm     sql.NullTime `gorm:"default:null"`                  
	DateFinish   sql.NullTime `gorm:"default:null"`                 
	ModeratorID  uuid.NullUUID       `gorm:"default:null"`                  
	
	Creator    User `gorm:"foreignKey:CreatorID"`
	Moderator  User `gorm:"foreignKey:ModeratorID"`
}