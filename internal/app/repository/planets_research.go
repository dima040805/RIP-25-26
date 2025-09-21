package repository

import (
	// apitypes "LAB1/internal/app/api_types"
	apitypes "LAB1/internal/app/api_types"
	"LAB1/internal/app/ds"
	"errors"
)

func (r *Repository) DeletePlanetFromResearch(researchId int, planetId int) (ds.Research, error) {
	userId := r.userId
	user, _ := r.GetUserByID(userId)
	var research ds.Research
	err := r.db.Where("id = ?", researchId).First(&research).Error
	if err != nil {
		return ds.Research{}, err
	}
	if research.CreatorID != r.userId && !user.IsModerator{
		return ds.Research{}, errors.New("you are not creator")
	}
	err = r.db.Where("planet_id = ? and Research_id = ?", planetId, researchId).Delete(&ds.PlanetsResearch{}).Error
	if err != nil {
		return ds.Research{}, err
	}
	return research, nil
}

// func (r *Repository) GetPlanetsResearch(researchId int) ([]ds.PlanetsResearch, error) {
// 	var planetsResearch []ds.PlanetsResearch
// 	err := r.db.Where("research_id = ?", researchId).Find(&planetsResearch).Error
// 	if err != nil {
// 		return nil, err
// 	}
// 	return planetsResearch, nil
// }

func (r *Repository) ChangePlanetResearch(reacherId int, planetId int, planetsResearchJSON apitypes.PlanetsResearchJSON) (ds.PlanetsResearch, error) {
	var planetsResearch ds.PlanetsResearch
	err := r.db.Model(&planetsResearch).Where("planet_id = ? and research_id = ?", planetId, reacherId).Updates(apitypes.PlanetsResearchFromJSON(planetsResearchJSON)).First(&planetsResearch).Error
	if err != nil {
		return ds.PlanetsResearch{}, err
	}
	return planetsResearch, nil
}