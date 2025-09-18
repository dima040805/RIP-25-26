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

	// Возвращаем объект Repository с подключенной базой данных
	return &Repository{
		db: db,
	}, nil
}

type PlanetsToResearch struct {
	PlanetId     int
	PlanetShine  float64
	PlanetRadius float64
}

type Research struct {
	ResearchDate     string
	PlanetsParametrs []PlanetsToResearch
}

// func (r *Repository) GetPlanets() ([]Planet, error) {

// 	planets := []Planet{
// 		{
// 			ID:          1,
// 			Image:       "http://localhost:9000/test/HD_209458_b.png",
// 			Name:        "HD_209458_b",
// 			Description: "HD 209458 b - это газовый гигант, похожий на Юпитер, который вращается вокруг звезды солнечного типа. Его масса составляет 71.36 масс Земли, он совершает один оборот вокруг своей звезды за 3.5 дня и находится на расстоянии 0.047 а.е. от своей звезды. Его открытие было объявлено в 1999 году.",
// 			Distance:    548000,
// 			Mass:        71.36,
// 			Discovery:   1999,
// 			StarRadius:  403500,
// 		},
// 		{
// 			ID:          2,
// 			Image:       "http://localhost:9000/test/GJ_1214_b.png",
// 			Name:        "GJ_1214_b",
// 			Description: "GJ 1214 b - это мини-нептун или супер-земля, который вращается вокруг красного карлика. Его масса составляет 22.7 масс Земли, он совершает один оборот вокруг своей звезды за 1.6 дня и находится на расстоянии 0.014 а.е. от своей звезды. Его открытие было объявлено в 2009 году.",
// 			Distance:    475000,
// 			Mass:        22.7,
// 			Discovery:   2009,
// 			StarRadius:  73100,
// 		},
// 		{
// 			ID:          3,
// 			Image:       "http://localhost:9000/test/WASP-12_b.png",
// 			Name:        "WASP-12_b",
// 			Description: "WASP-12 b - это горячий юпитер, который вращается вокруг звезды солнечного типа. Его масса составляет 145.5 масс Земли, он совершает один оборот вокруг своей звезды за 1.1 дня и находится на расстоянии 0.023 а.е. от своей звезды. Его открытие было объявлено в 2008 году.",
// 			Distance:    527000,
// 			Mass:        145.5,
// 			Discovery:   2008,
// 			StarRadius:  577500,
// 		},
// 		{
// 			ID:          4,
// 			Image:       "http://localhost:9000/test/KELT-9_b.png",
// 			Name:        "KELT-9_b",
// 			Description: "KELT-9 b - это ультра-горячий юпитер, который вращается вокруг горячей звезды типа A. Его масса составляет 50.5 масс Земли, он совершает один оборот вокруг своей звезды за 1.5 дня и находится на расстоянии 0.035 а.е. от своей звезды. Его открытие было объявлено в 2016 году.",
// 			Distance:    621000,
// 			Mass:        50.5,
// 			Discovery:   2016,
// 			StarRadius:  679500,
// 		},
// 		{
// 			ID:          5,
// 			Image:       "http://localhost:9000/test/HAT-P-7_b.png",
// 			Name:        "HAT-P-7_b",
// 			Description: "HAT-P-7 b - это горячий юпитер, который вращается вокруг звезды солнечного типа. Его масса составляет 30.5 масс Земли, он совершает один оборот вокруг своей звезды за 2.2 дня и находится на расстоянии 0.038 а.е. от своей звезды. Его открытие было объявлено в 2008 году.",
// 			Distance:    104400,
// 			Mass:        30.5,
// 			Discovery:   2008,
// 			StarRadius:  512000,
// 		},
// 		{
// 			ID:          6,
// 			Image:       "http://localhost:9000/test/WASP-121_b.png",
// 			Name:        "WASP-121_b",
// 			Description: "WASP-121 b - это горячий юпитер, который вращается вокруг звезды солнечного типа. Его масса составляет 38.5 масс Земли, он совершает один оборот вокруг своей звезды за 1.3 дня и находится на расстоянии 0.025 а.е. от своей звезды. Его открытие было объявлено в 2015 году.",
// 			Distance:    858000,
// 			Mass:        38.5,
// 			Discovery:   2015,
// 			StarRadius:  665000,
// 		},
// 		{
// 			ID:          7,
// 			Image:       "http://localhost:9000/test/WASP-17_b.png",
// 			Name:        "WASP-17_b",
// 			Description: "WASP-17 b - это раздутый горячий юпитер, который вращается вокруг звезды солнечного типа. Его масса составляет 77.4 масс Земли, он совершает один оборот вокруг своей звезды за 3.7 дня и находится на расстоянии 0.051 а.е. от своей звезды. Его открытие было объявлено в 2009 году.",
// 			Distance:    1340000,
// 			Mass:        77.4,
// 			Discovery:   2009,
// 			StarRadius:  1075000,
// 		},
// 		{
// 			ID:          8,
// 			Image:       "http://localhost:9000/test/WASP-19_b.png",
// 			Name:        "WASP-19_b",
// 			Description: "WASP-19 b - это горячий юпитер, который вращается вокруг звезды солнечного типа. Его масса составляет 63.9 масс Земли, он совершает один оборот вокруг своей звезды за 0.8 дня и находится на расстоянии 0.016 а.е. от своей звезды. Его открытие было объявлено в 2009 году.",
// 			Distance:    879000,
// 			Mass:        63.9,
// 			Discovery:   2009,
// 			StarRadius:  668000,
// 		},
// 		{
// 			ID:          9,
// 			Image:       "http://localhost:9000/test/WASP-43_b.png",
// 			Name:        "WASP-43_b",
// 			Description: "WASP-43 b - это горячий юпитер, который вращается вокруг красного карлика. Его масса составляет 30.0 масс Земли, он совершает один оборот вокруг своей звезды за 0.8 дня и находится на расстоянии 0.014 а.е. от своей звезды. Его открытие было объявлено в 2011 году.",
// 			Distance:    742000,
// 			Mass:        30.0,
// 			Discovery:   2011,
// 			StarRadius:  617000,
// 		},
// 	}
// 	if len(planets) == 0 {
// 		return nil, fmt.Errorf("массив пустой")
// 	}

// 	return planets, nil
// }

// func (r *Repository) GetPlanet(id int) (Planet, error) {
// 	// тут у вас будет логика получения нужной услуги, тоже наверное через цикл в первой лабе, и через запрос к БД начиная со второй
// 	planets, err := r.GetPlanets()
// 	if err != nil {
// 		return Planet{}, err // тут у нас уже есть кастомная ошибка из нашего метода, поэтому мы можем просто вернуть ее
// 	}

// 	for _, planet := range planets {
// 		if int(planet.ID) == id { // если нашли, то просто возвращаем найденный заказ (услугу) без ошибок
// 			return planet, nil
// 		}
// 	}
// 	return Planet{}, fmt.Errorf("заказ не найден") // тут нужна кастомная ошибка, чтобы понимать на каком этапе возникла ошибка и что произошло
// }

// func (r *Repository) GetPlanetsByName(name string) ([]Planet, error) {
// 	planets, err := r.GetPlanets()
// 	if err != nil {
// 		return []Planet{}, err
// 	}

// 	var result []Planet
// 	for _, planet := range planets {
// 		print(planet.Name, name)
// 		if strings.Contains(strings.ToLower(planet.Name), strings.ToLower(name)) {
// 			result = append(result, planet)
// 		}
// 	}

// 	return result, nil
// }








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

// func (r *Repository) GetResearchCount(id int) int {
// 	researchPlanets := r.GetResearch(id)
// 	return len(researchPlanets.PlanetsParametrs)
// }

// func (r *Repository) GetResearchId() int {
// 	return 1
// }
