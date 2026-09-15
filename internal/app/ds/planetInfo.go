package ds

type PlanetInfo struct {
	ID          int `gorm:"primaryKey"`
	Image       string
	Name        string `gorm:"type:varchar(25);not null"`

	StarRadius  int

	PlanetShine float64
	PlanetRadius int
	
}