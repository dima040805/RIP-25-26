package repository

import (
	apitypes "LAB1/internal/app/api_types"
	"LAB1/internal/app/ds"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

func (r *Repository) DeletePlanetFromResearch(researchId int, planetId int) (ds.Research, error) {
	// userId := r.userId
    // if userId == 0 {
    //     return ds.Research{}, fmt.Errorf("%w: пользователь не авторизирован", ErrNotAllowed)
    // }
    
	// user, err := r.GetUserByID(userId)
	// if err != nil {
	// 	return ds.Research{}, err
	// }
    
	var research ds.Research
	err := r.db.Where("id = ?", researchId).First(&research).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ds.Research{}, fmt.Errorf("%w: исследование с id %d", ErrNotFound, researchId)
		}
		return ds.Research{}, err
	}
    
	// if research.CreatorID != r.userId && !user.IsModerator{
	// 	return ds.Research{}, fmt.Errorf("%w: Вы не создатель этого исследования", ErrNotAllowed)
	// }
    
	err = r.db.Where("planet_id = ? and Research_id = ?", planetId, researchId).Delete(&ds.PlanetsResearch{}).Error
	if err != nil {
		return ds.Research{}, err
	}
	return research, nil
}

func (r *Repository) ChangePlanetResearch(researchId int, planetId int, planetsResearchJSON apitypes.PlanetsResearchJSON) (ds.PlanetsResearch, error) {
	var planetsResearch ds.PlanetsResearch
	err := r.db.Model(&planetsResearch).Where("planet_id = ? and research_id = ?", planetId, researchId).Updates(apitypes.PlanetsResearchFromJSON(planetsResearchJSON)).First(&planetsResearch).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ds.PlanetsResearch{}, fmt.Errorf("%w: планеты в исследовании", ErrNotFound)
		}
		return ds.PlanetsResearch{}, err
	}
	return planetsResearch, nil
}