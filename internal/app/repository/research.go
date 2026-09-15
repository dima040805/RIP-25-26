package repository

import (
	"LAB1/internal/app/api_types"
	"LAB1/internal/app/ds"
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)


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
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ds.PlanetsResearch{}, fmt.Errorf("%w: planets research not found", ErrNotFound)
		}
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
	err = r.db.Order("id").Where("id IN (?)", sub.Select("planet_id")).Find(&planets).Error

	if err != nil {
		return []ds.Planet{}, ds.Research{}, err
	}

	return planets, research, nil
}

func (r *Repository) CheckCurrentResearchDraft(creatorID uuid.UUID) (ds.Research, error) {
    // if creatorID == 0 {
    //     return ds.Research{}, fmt.Errorf("%w: user not authenticated", ErrNotAllowed)
    // }
    
	var research ds.Research
	res := r.db.Where("creator_id = ? AND status = ?", creatorID, "draft").Limit(1).Find(&research)
	if res.Error != nil {
		return ds.Research{}, res.Error
	} else if res.RowsAffected == 0 {
		return ds.Research{}, ErrNoDraft
	}
	return research, nil
}

func (r *Repository) GetResearchDraft(creatorID uuid.UUID) (ds.Research, bool, error) {
    // if creatorID == 0 {
    //     return ds.Research{}, false, fmt.Errorf("%w: user not authenticated", ErrNotAllowed)
    // }
    
	research, err := r.CheckCurrentResearchDraft(creatorID)
	if errors.Is(err, ErrNoDraft) {
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

func (r *Repository) GetResearchCount(creatorID uuid.UUID) int64 {
    
	var count int64
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
		return ds.Research{}, errors.New("неверное id, должно быть >= 0")
	}
    
    // userId := r.GetUserID()
    // if userId == 0 {
    //     return ds.Research{}, fmt.Errorf("%w: пользователь не авторизирован", ErrNotAllowed)
    // }
    
	// user, err := r.GetUserByID(userId)
	// if err != nil {
	// 	return ds.Research{}, err
	// }
    
	var research ds.Research
	err := r.db.Where("id = ?", id).First(&research).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ds.Research{}, fmt.Errorf("%w: заявка с id %d", ErrNotFound, id)
		}
		return ds.Research{}, err
	} else if research.Status == "deleted"  {
		return ds.Research{}, fmt.Errorf("%w: заявка удалена", ErrNotAllowed)
	}
	return research, nil
}

func (r *Repository) FormResearch(researchId int, status string) (ds.Research, error) {
	research, err := r.GetSingleResearch(researchId)
	if err != nil {
		return ds.Research{}, err
	}


	if research.Status != "draft" {
		return ds.Research{}, fmt.Errorf("эта заявка не может быть %s", status)
	}
	
	if status != "deleted"{
		if research.DateResearch == "" {
			return ds.Research{}, errors.New("вы не написали дату исследования")
		}
		planetsResearch, _ := r.GetPlanetsResearches(research.ID)
		for _, planetResearch := range planetsResearch {
				if planetResearch.PlanetShine == 0{
					return ds.Research{}, errors.New("вы не написали блеск планеты" )			
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
		return ds.Research{}, errors.New("неправильное id, должно быть >= 0")
	}
	if researchJSON.DateResearch == "" {  
		return ds.Research{}, errors.New("неправильная дата исследования")
	}
	err := r.db.Where("id = ? and status != 'deleted'", id).First(&research).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ds.Research{}, fmt.Errorf("%w: исследование с id %d", ErrNotFound, id)
		}
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
		return 0, errors.New("неправильная дата исследования")
	}
	if planetShine < 0 || planetShine > 7 {
		return 0, errors.New("неправильный блеск")
	}
	return float64(starRadius) * math.Sqrt(float64(planetShine / 100)) , nil
}

func (r *Repository) ModerateResearch(id int, status string, currUserId uuid.UUID) (ds.Research, error) {
	if status != "completed" && status != "rejected" {
		return ds.Research{}, errors.New("неверный статус")
	}

	research, err := r.GetSingleResearch(id)
	if err != nil {
		return ds.Research{}, err
	} else if research.Status != "formed" {
		return ds.Research{}, errors.New("this calculation can not be " + status)
	}

	// Обновляем статус исследования
	err = r.db.Model(&research).Updates(ds.Research{
		Status: status,
		DateFinish: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
		ModeratorID: uuid.NullUUID{
			UUID:  currUserId,
			Valid: true,
		},
	}).Error
	if err != nil {
		return ds.Research{}, err
	}

	// Если исследование одобрено - запускаем асинхронный расчет радиусов для КАЖДОЙ планеты
	if status == "completed" {
		planetsResearch, err := r.GetPlanetsResearches(research.ID)
		if err != nil {
			return ds.Research{}, err
		}
		
		// Для каждой планеты запускаем асинхронный расчет
		for _, planetResearch := range planetsResearch {
			planet, err := r.GetPlanet(int(planetResearch.PlanetID))
			if err != nil {
				logrus.Errorf("Error getting planet %d: %v", planetResearch.PlanetID, err)
				continue
			}
			
			// Запускаем асинхронный расчет радиуса для этой планеты
			go r.calculateSinglePlanetRadiusAsync(research.ID, planet, planetResearch)
		}
		
		logrus.Infof("Started async radius calculations for research %d, planets: %d", research.ID, len(planetsResearch))
	}

	return research, nil
}

func (r *Repository) calculateSinglePlanetRadiusAsync(researchId int, planet *ds.Planet, planetResearch ds.PlanetsResearch) {
    // URL Django сервиса для расчета одного радиуса
    asyncServiceURL := "http://localhost:8000/api/calculate-radius/"
    
    // Данные для расчета ОДНОГО радиуса
    requestData := map[string]interface{}{
        "research_id": researchId,
        "planet_id":   planet.ID,
        "star_radius": planet.StarRadius,
        "planet_shine": planetResearch.PlanetShine,
    }

    jsonData, err := json.Marshal(requestData)
    if err != nil {
        logrus.Errorf("Error marshaling request data for planet %d: %v", planet.ID, err)
        return
    }

    // Отправляем запрос в асинхронный сервис
    resp, err := http.Post(asyncServiceURL, "application/json", bytes.NewBuffer(jsonData))
    if err != nil {
        logrus.Errorf("Error sending request to async service for planet %d: %v", planet.ID, err)
        return
    }
    defer resp.Body.Close()

    if resp.StatusCode != 202 {
        logrus.Errorf("Async service returned status %d for planet %d", resp.StatusCode, planet.ID)
        return
    }

    logrus.Infof("Successfully started async calculation for research %d, planet %d", researchId, planet.ID)
}

func (r *Repository) UpdatePlanetRadius(researchId int, planetId int, planetRadius int) error {
    var planetsResearch ds.PlanetsResearch
    err := r.db.Where("research_id = ? AND planet_id = ?", researchId, planetId).First(&planetsResearch).Error
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return fmt.Errorf("%w: связь планеты с исследованием не найдена", ErrNotFound)
        }
        return err
    }

    err = r.db.Model(&planetsResearch).Update("planet_radius", planetRadius).Error
    if err != nil {
        return err
    }

    logrus.Printf("Updated planet radius: research=%d, planet=%d, radius=%d", researchId, planetId, planetRadius)
    return nil
}

func (r *Repository) GetCalculatedPlanetsCount(researchId int) (int, error) {
    var count int64
    err := r.db.Model(&ds.PlanetsResearch{}).
        Where("research_id = ? AND planet_radius IS NOT NULL AND planet_radius > 0", researchId).
        Count(&count).Error
    if err != nil {
        return 0, err
    }
    return int(count), nil
}