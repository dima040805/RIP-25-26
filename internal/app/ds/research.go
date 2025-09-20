package ds

import (
	"database/sql"
	"time"
)

type Research struct {
	ID           uint   `gorm:"primaryKey"`
	DateResearch time.Time
	Status       string `gorm:"type:varchar(15);not null"`
	DateCreate   time.Time `gorm:"not null"`
	DateForm    sql.NullTime  `gorm:"default:null"`
	DateFinish  sql.NullTime  `gorm:"default:null"`
	CreatorID    int       `gorm:"not null"`
	ModeratorID sql.NullInt64 `gorm:"default:null"`

	Creator Users `gorm:"foreignKey:CreatorID"`
	Moderator Users `gorm:"foreignKey:ModeratorID"`
}
