package repository

import (
	// "database/sql"
	// "errors"
	"errors"
	"fmt"

	"LAB1/internal/app/api_types"
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

func (r *Repository) CreatePlanet(planetJSON apitypes.PlanetJSON) (ds.Planet, error) {
	planet := apitypes.PlanetFromJSON(planetJSON)
	if planet.StarRadius <= 0 {
		return ds.Planet{}, errors.New("invalid star radius")
	}
	if planet.Mass <= 0 {
		return ds.Planet{}, errors.New("invalid mass")
	}
	if planet.Mass <= 0 {
		return ds.Planet{}, errors.New("invalid mass")
	}
	err := r.db.Create(&planet).First(&planet).Error
	if err != nil {
		return ds.Planet{}, err
	}
	return planet, nil
}

func (r *Repository) ChangePlanet(id int, planetJSON apitypes.PlanetJSON) (ds.Planet, error) {
	planet := ds.Planet{}
	if id < 0 {
		return ds.Planet{}, errors.New("invalid id, it must be >= 0")
	}
	err := r.db.Where("id = ? and is_delete = ?", id, false).First(&planet).Error
	if err != nil {
		return ds.Planet{}, err
	}
	if planetJSON.StarRadius <= 0 {
		return ds.Planet{}, errors.New("invalid star radius")
	}
	err = r.db.Model(&planet).Updates(apitypes.PlanetFromJSON(planetJSON)).Error
	if err != nil {
		return ds.Planet{}, err
	}
	return planet, nil
}

func (r *Repository) DeletePlanet(id int) error {
	planet := ds.Planet{}
	if id < 0 {
		return errors.New("invalid id, it must be >= 0")
	}

	err := r.db.Where("id = ? and is_delete = ?", id, false).First(&planet).Error
	if err != nil {
		return err
	}
	// if reaction.ImgLink != "" {
	// 	err = minioInclude.DeleteObject(context.Background(), r.mc, minioInclude.GetImgBucket(), reaction.ImgLink)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	err = r.db.Model(&ds.Planet{}).Where("id = ?", id).Update("is_delete", true).Error
	if err != nil {
		return err
	}
	return nil
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

func (r *Repository) GetModeratorAndCreatorLogin(research ds.Research) (string, string, error) {
	var creator ds.User
	var moderator ds.User

	err := r.db.Where("id = ?", research.CreatorID).First(&creator).Error
	if err != nil {
		return "", "", err
	}

	var moderatorLogin string
	if research.ModeratorID.Valid {
		err = r.db.Where("id = ?", research.ModeratorID.Int64).First(&moderator).Error
		if err != nil {
			return "", "", err
		}
		moderatorLogin = moderator.Login
	}
	
	return creator.Login, moderatorLogin, nil
}