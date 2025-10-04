package repository

import (
	// "database/sql"
	// "errors"
	"fmt"

	"LAB1/internal/app/ds"
)

func (r *Repository) GetPlanets() ([]ds.Planet, error) {
	var planets []ds.Planet
	err := r.db.Order("id").Where("is_delete = false").Find(&planets).Error
	if err != nil {
		return nil, err
	}
	if len(planets) == 0 {
		return nil, fmt.Errorf("массив пустой")
	}

	return planets, nil
}

func (r *Repository) GetPlanet(id int) (*ds.Planet, error) {
// 	query := "SELECT id, image, name, description, distance, mass, discovery, star_radius FROM planets WHERE id = $1"
// 	row := r.db.Raw(query, id).Row()
// 	planet := &ds.Planet{}

//    err := row.Scan(
// 		&planet.ID,
//       &planet.Image,
//       &planet.Name,
//       &planet.Description,
//       &planet.Distance,
//       &planet.Mass,
// 	  &planet.Discovery,
//       &planet.StarRadius,
//    )

//    if err != nil {
//       if errors.Is(err, sql.ErrNoRows) {
//          return nil, nil // Возвращаем nil, если записи нет
//       }
//       return nil, err
//    }
// 	return planet, nil

	planet := ds.Planet{}
	err := r.db.Order("id").Where("id = ? and is_delete = ?", id, false).First(&planet).Error
	if err != nil {
		return &ds.Planet{}, err
	}
	return &planet, nil
}

func (r *Repository) GetPlanetsByName(name string) ([]ds.Planet, error) {
	var planets []ds.Planet
	err := r.db.Order("id").Where("name ILIKE ? and is_delete = ?", "%"+name+"%", false).Find(&planets).Error
	if err != nil {
		return nil, err
	}
	return planets, nil
}


func (r *Repository) AddPlanetToResearch(researchId int, planetId int) error {
	var planet ds.Planet
	if err := r.db.First(&planet, planetId).Error; err != nil {
		return err
	}

	var research ds.Research
	if err := r.db.First(&research, researchId).Error; err != nil {
		return err
	}
	planetsResearch := ds.PlanetsResearch{}
	result := r.db.Where("planet_id = ? and research_id = ?", planetId, researchId).Find(&planetsResearch)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected != 0 {
		return nil
	}
	return r.db.Create(&ds.PlanetsResearch{
		PlanetID:    uint(planetId),
		ResearchID: uint(researchId),
	}).Error

}