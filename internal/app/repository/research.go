package repository

import (
	"LAB1/internal/app/api_types"
	"LAB1/internal/app/ds"
	"database/sql"
	"errors"
	"math"
	"time"

	"github.com/sirupsen/logrus"
	// "time"
	// "github.com/sirupsen/logrus"
)

var errNoDraft = errors.New("no draft for this user")

func (r *Repository) GetResearches(from, to time.Time, status string) ([]ds.Research, error) {
	var researches []ds.Research
	sub := r.db.Where("status != 'deleted' and status != 'draft'")
	if !from.IsZero() {
		sub = sub.Where("date_create > ?", from)
	}
	if !to.IsZero() {
		sub = sub.Where("date_create < ?", to.Add(time.Hour*24))
	}
	if status != "" {
		sub = sub.Where("status = ?", status)
	}
	err := sub.Order("id").Find(&researches).Error
	if err != nil {
		return nil, err
	}
	return researches, nil
}

func (r *Repository) GetPlanetsResearches(researchId int) ([]ds.PlanetsResearch, error) {
	var planetsResearches []ds.PlanetsResearch
	err := r.db.Where("research_id = ?", researchId).Find(&planetsResearches).Error
	if err != nil {
		return nil, err
	}
	return planetsResearches, nil
}

func (r *Repository) GetPlanetsResearch(planetId int, researchId int) (ds.PlanetsResearch, error) {
	var planetsResearch ds.PlanetsResearch
	err := r.db.Where("planet_id = ? and research_id = ?", planetId, researchId).First(&planetsResearch).Error
	if err != nil {
		return ds.PlanetsResearch{}, err
	}
	return planetsResearch, nil
}

func (r *Repository) GetResearchPlanets(id int) ([]ds.Planet, ds.Research, error) {
	research, err := r.GetSingleResearch(id)
	if err != nil {
		return []ds.Planet{}, ds.Research{}, err
	}

	var planets []ds.Planet
	sub := r.db.Table("planets_researches").Where("Research_id = ?", research.ID)
	err = r.db.Order("id DESC").Where("id IN (?)", sub.Select("planet_id")).Find(&planets).Error

	if err != nil {
		return []ds.Planet{}, ds.Research{}, err
	}

	return planets, research, nil
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


func (r *Repository) GetResearchDraft(creatorID int) (ds.Research, bool, error) {
	research, err := r.CheckCurrentResearchDraft(creatorID)
	if err == errNoDraft {
		research = ds.Research{
			Status:     "draft",
			CreatorID:  creatorID,
			DateCreate: time.Now(),
		}
		result := r.db.Create(&research)
		if result.Error != nil {
			return ds.Research{}, false, result.Error
		}
		return research, true, nil
	} else if err != nil {
		return ds.Research{}, false, err
	}
	return research, true, nil
}

func (r *Repository) GetResearchCount(creatorID int) int64 {
	var count int64
	// пока что мы захардкодили id создателя заявки, в последующем вы сделаете авторизацию и будете получать его из JWT

	research, err := r.CheckCurrentResearchDraft(creatorID)
	if err != nil {
		return 0
	}
	err = r.db.Model(&ds.PlanetsResearch{}).Where("research_id = ?", research.ID).Count(&count).Error
	if err != nil {
		logrus.Println("Error counting records in lists_planets:", err)
	}

	return count
}


func (r *Repository) DeleteCalculation(researchId int) error{
	return r.db.Exec("UPDATE researches SET status = 'deleted' WHERE id = ?", researchId).Error
}


func (r *Repository) GetSingleResearch(id int) (ds.Research, error) {
	if id < 0 {
		return ds.Research{}, errors.New("invalid id, it must be >= 0")
	}
	user, err := r.GetUserByID(r.GetUserID())
	if err != nil {
		return ds.Research{}, err
	}
	var research ds.Research
	err = r.db.Where("id = ?", id).First(&research).Error
	if err != nil {
		return ds.Research{}, err
	// } else if user.ID != calculation.CreatorID && !user.IsModerator {
	// 	return ds.Calculation{}, ErrorNotAllowed
	} else if research.Status == "deleted" && !user.IsModerator {
		return ds.Research{}, errors.New("calculation is deleted")
	}
	return research, nil
}

func (r *Repository) FormResearch(researchId int, status string) (ds.Research, error) {
	research, err := r.GetSingleResearch(researchId)
	if err != nil {
		return ds.Research{}, err
	}

	user, err := r.GetUserByID(r.GetUserID())

	if err != nil{
		return ds.Research{}, errors.New("you are not registered")
	}

	if research.CreatorID != r.userId && !user.IsModerator{
		return ds.Research{}, errors.New("you do not have the rights to have this research " + status)
	}

	if research.Status != "draft" {
		return ds.Research{}, errors.New("this research can not be " + status)
	}
		if status != "deleted"{
			if research.DateResearch == "" {
				return ds.Research{}, errors.New("you don't write date research")
			}
		planetsResearch, _ := r.GetPlanetsResearches(research.ID)
		for _, planetResearch := range planetsResearch {
				if planetResearch.PlanetShine == 0{
					return ds.Research{}, errors.New("you don't write planet shine" )			
				}
		}
	}	

	err = r.db.Model(&research).Updates(ds.Research{
		Status: status,
		DateForm: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
	}).Error
	if err != nil {
		return ds.Research{}, err
	}

	return research, nil
}

func (r *Repository) ChangeResearch(id int, researchJSON apitypes.ResearchJSON) (ds.Research, error) {
	research := ds.Research{}
	if id < 0 {
		return ds.Research{}, errors.New("invalid id, it must be >= 0")
	}
	if researchJSON.DateResearch == "" {  
		return ds.Research{}, errors.New("invalid date research")
	}
	err := r.db.Where("id = ? and status != 'deleted'", id).First(&research).Error
	if err != nil {
		return ds.Research{}, err
	}
	err = r.db.Model(&research).Updates(apitypes.ResearchFromJSON(researchJSON)).Error
	if err != nil {
		return ds.Research{}, err
	}
	return research, nil
}


func CalculatePlanetRadius(starRadius int, dateResearch string, planetShine float64) (float64, error) {
	if dateResearch == "" {
		return 0, errors.New("invalid conversation factor")
	}
	if planetShine < 0 || planetShine > 3 {
		return 0, errors.New("invalid output koeficient")
	}
	return float64(starRadius) * math.Sqrt(float64(planetShine / 100)) , nil
}

func (r *Repository) ModerateResearch(id int, status string) (ds.Research, error) {
	if status != "completed" && status != "rejected" {
		return ds.Research{}, errors.New("wrong status")
	}

	user, err := r.GetUserByID(r.GetUserID())
	if err != nil {
		return ds.Research{}, err
	}

	if !user.IsModerator {
		return ds.Research{}, errors.New("you are not a moderator")
	}

	research, err := r.GetSingleResearch(id)
	if err != nil {
		return ds.Research{}, err
	} else if research.Status != "formed" {
		return ds.Research{}, errors.New("this calculation can not be " + status)
	}

	err = r.db.Model(&research).Updates(ds.Research{
		Status: status,
		DateFinish: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
		ModeratorID: sql.NullInt64{
			Int64: int64(user.ID),
			Valid: true,
		},
	}).Error
	if err != nil {
		return ds.Research{}, err
	}

	if status == "completed" {
		planetsResearch, err := r.GetPlanetsResearches(research.ID)
		if err != nil {
			return ds.Research{}, err
		}
		for _, planetResearch := range planetsResearch {
			planet, err := r.GetPlanet(int(planetResearch.PlanetID))
			if err != nil {
				return ds.Research{}, err
			}
			planetRadius, err := CalculatePlanetRadius(planet.StarRadius, research.DateResearch, planetResearch.PlanetShine)
			if err != nil {
				return ds.Research{}, err
			}
			err = r.db.Model(&planetResearch).Updates(ds.PlanetsResearch{
				PlanetRadius: int(planetRadius),
			}).Error
			if err != nil {
				return ds.Research{}, err
			}
		}
	}

	return research, nil
}