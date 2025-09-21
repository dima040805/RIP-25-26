package ds



type Planet struct {
	ID          int     `gorm:"primaryKey;autoIncrement"`
	Image       string  `gorm:"type:varchar(50)"`             
	Name        string  `gorm:"type:varchar(25);not null"`     
	Description string  `gorm:"type:varchar(512)"`             
	Distance    int     `gorm:"default: null"`                    
	Mass        float64 `gorm:"default: null"`                      
	Discovery   int     `gorm:"default: null"`                      
	StarRadius  int     `gorm:"not null"`                      
	IsDelete    bool    `gorm:"not null;default:false"`       
}