package repository

import (
	"fmt"

	"LAB1/internal/app/ds"
)

func (r *Repository) GetPlanets() ([]ds.Planet, error) {
	var planets []ds.Planet
	err := r.db.Find(&planets).Error
	if err != nil {
		return nil, err
	}
	if len(planets) == 0 {
		return nil, fmt.Errorf("массив пустой")
	}

	return planets, nil
}

func (r *Repository) GetPlanet(id int) (ds.Planet, error) {
	planet := ds.Planet{}
	err := r.db.Where("id = ?", id).First(&planet).Error
	if err != nil {
		return ds.Planet{}, err
	}
	return planet, nil
}


func (r *Repository) GetPlanetsByName(name string) ([]ds.Planet, error) {
	var planets []ds.Planet
	err := r.db.Where("name ILIKE ?", "%"+name+"%").Find(&planets).Error
	if err != nil {
		return nil, err
	}
	return planets, nil
}