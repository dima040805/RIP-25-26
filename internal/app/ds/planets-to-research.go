package ds


type PlanetsResearch struct {
	ID           uint     `gorm:"primaryKey;autoIncrement"`
	ResearchID   uint     `gorm:"not null;uniqueIndex:idx_research_chat"`
	PlanetID     uint     `gorm:"not null;uniqueIndex:idx_research_chat"`
	PlanetShine  float64 `gorm:"default: null"`                    
	PlanetRadius int `gorm:"default: null"`                    
	
	Research Research `gorm:"foreignKey:ResearchID"`
	Planet   Planet   `gorm:"foreignKey:PlanetID"`
}