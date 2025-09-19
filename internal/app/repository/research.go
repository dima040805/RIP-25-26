package repository

import (
	"LAB1/internal/app/ds"
	"errors"
	"time"

	"github.com/sirupsen/logrus"
	// "time"
	// "github.com/sirupsen/logrus"
)

var errNoDraft = errors.New("no draft for this user")


func (r *Repository) GetPlanetsResearch(id int) ([]ds.PlanetInfo, ds.Research, error) {

	creatorID := r.GetUser()
	// пока что мы захардкодили id создателя заявки, в последующем вы сделаете авторизацию и будете получать его из JWT

	var research ds.Research
	err := r.db.Where("id = ?", id).First(&research).Error
	if err != nil {
		return []ds.PlanetInfo{}, ds.Research{}, err
	} else if creatorID != int(research.CreatorID) {
		return []ds.PlanetInfo{}, ds.Research{}, errors.New("you are not allowed")
	} else if research.Status == "deleted" {
		return []ds.PlanetInfo{}, ds.Research{}, errors.New("you can`t watch deleted calculations")
	}

	var planets []ds.Planet
	var planetsResearches []ds.PlanetsResearch
	sub := r.db.Table("planets_researches").Where("research_id = ?", research.ID).Find(&planetsResearches)
	err = r.db.Where("id IN (?)", sub.Select("planet_id")).Find(&planets).Error
	if err != nil {
		return []ds.PlanetInfo{}, ds.Research{}, err
	}


	var planetsResult []ds.PlanetInfo
	for _, planet := range planets {
		for _, planetsResearch := range planetsResearches {
			if planet.ID == int(planetsResearch.PlanetID) {
				planetsResult = append(planetsResult, ds.PlanetInfo{
					ID:                 planet.ID,
					Name:              	planet.Name,
					Image:            	planet.Image,
					StarRadius:         planet.StarRadius,
					
					PlanetShine: 		planetsResearch.PlanetShine,
					PlanetRadius:       planetsResearch.PlanetRadius,
				})
				break
			}
		}
	}

	return planetsResult, research, nil
}


func (r *Repository) CheckCurrentResearchDraft(creatorID int) (ds.Research, error) {
	var research ds.Research

	res := r.db.Where("creator_id = ? AND status = ?", creatorID, "draft").Limit(1).Find(&research)
	if res.Error != nil {
		return ds.Research{}, res.Error
	} else if res.RowsAffected == 0 {
		return ds.Research{}, errNoDraft
	}
	return research, nil
}


func (r *Repository) GetResearchDraft(creatorID int) (ds.Research, error) {
	research, err := r.CheckCurrentResearchDraft(creatorID)
	if err == errNoDraft {
		research = ds.Research{
			Status:     "draft",
			CreatorID:  creatorID,
			DateCreate: time.Now(),
		}
		result := r.db.Create(&research)
		if result.Error != nil {
			return ds.Research{}, result.Error
		}
		return research, nil
	} else if err != nil {
		return ds.Research{}, err
	}
	return research, nil
}

func (r *Repository) GetResearchCount() int64 {
	var researchID uint
	var count int64
	creatorID := 1
	// пока что мы захардкодили id создателя заявки, в последующем вы сделаете авторизацию и будете получать его из JWT

	err := r.db.Model(&ds.Research{}).Where("creator_id = ? AND status = ?", creatorID, "draft").Select("id").First(&researchID).Error
	if err != nil {
		return 0
	}
	err = r.db.Model(&ds.PlanetsResearch{}).Where("research_id = ?", researchID).Count(&count).Error
	if err != nil {
		logrus.Println("Error counting records in lists_planets:", err)
	}

	return count
}


func (r *Repository) DeleteCalculation(researchId int) error{
	return r.db.Exec("UPDATE researches SET status = 'deleted' WHERE id = ?", researchId).Error
}