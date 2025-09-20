package ds

type Planet struct {
	ID          int `gorm:"primaryKey"`
	IsDelete    bool   `gorm:"type:boolean not null;default:false"`
	Image       string `gorm:"type: varchar(32)"`
	Name        string `gorm:"type:varchar(25);not null"`
	Distance    int
	Description string `gorm:"type: varchar(512)"`
	Mass        float64 `gorm:"type: float"`
	Discovery   int
	StarRadius  int
}