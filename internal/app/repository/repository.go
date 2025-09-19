package repository

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(dsn string) (*Repository, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{}) // подключаемся к БД
	if err != nil {
		return nil, err
	}

	return &Repository{
		db: db,
	}, nil
}


func (r *Repository) GetUser() (int) {
	return 1
}


// func (r *Repository) GetResearch(id int) Research {
// 	researchPlanets := map[int]Research{
// 		1: { // ID Заявки
// 			ResearchDate: "16.09.25",
// 			PlanetsParametrs: []PlanetsToResearch{ // Поля м-м
// 				{
// 					PlanetId:      1,
// 					PlanetShine:   1.5, // Вводится со странички заявки
// 					PlanetRadius:  98500, // Поле расчета, выводится на страничку заявки
// 				},
// 				{
// 					PlanetId:      3,
// 					PlanetShine:   1.3,
// 					PlanetRadius:  12,
// 				},
// 				{
// 					PlanetId:      5,
// 					PlanetShine:   1.5,
// 					PlanetRadius:  13,
// 				},
// 			},
// 		},
// 	}

// 	return researchPlanets[id]
// }

// func (r *Repository) GetResearchPlanets(id int) []Planet{
// 	research := r.GetResearch(id)

// 	var planetsInGroup []Planet
// 	for _, _planet := range research.PlanetsParametrs {
// 		planet, err := r.GetPlanet(_planet.PlanetId)
// 		if err == nil {
// 			planetsInGroup = append(planetsInGroup, planet)
// 		}
// 	}
// 	return planetsInGroup
// }


