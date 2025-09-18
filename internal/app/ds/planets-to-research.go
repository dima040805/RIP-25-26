package ds

type PlanetsResearch struct {
	ID        uint `gorm:"primaryKey"`
	// здесь создаем Unique key, указывая общий uniqueIndex
	ResearchID uint `gorm:"not null;uniqueIndex:idx_research_chat"`
	PlanetID    uint `gorm:"not null;uniqueIndex:idx_research_chat"`
	PlanetShine float64
	PlanetRadius int

	
	Research Research `gorm:"foreignKey:ResearchID"`
	Planet    Planet    `gorm:"foreignKey:PlanetID"`
}