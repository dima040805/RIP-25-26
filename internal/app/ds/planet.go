package ds

type Planet struct {
	ID          int `gorm:"primaryKey"`
	IsDelete    bool   `gorm:"type:boolean not null;default:false"`
	Image       string
	Name        string `gorm:"type:varchar(25);not null"`
	Distance    int
	Description string
	Mass        float64
	Discovery   int
	StarRadius  int
}