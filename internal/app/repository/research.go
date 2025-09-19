package repository

import (
	"LAB1/internal/app/ds"
	"errors"
	"time"

	"github.com/sirupsen/logrus"
)

var noDraftError = errors.New("no draft for this user")


func (r *Repository) GetPlanetsResearch(id int) ([]ds.PlanetsResearch, ds.Research, error) {

	creatorID := r.GetUser()
	// пока что мы захардкодили id создателя заявки, в последующем вы сделаете авторизацию и будете получать его из JWT

	var research ds.Research
	err := r.db.Where("id = ?", id).First(&research).Error
	if err != nil {
		return []ds.ReactionInfo{}, ds.Calculation{}, err
	} else if creatorID != calculation.CreatorID {
		return []ds.ReactionInfo{}, ds.Calculation{}, errors.New("you are not allowed")
	} else if calculation.Status == "deleted" {
		return []ds.ReactionInfo{}, ds.Calculation{}, errors.New("you can`t watch deleted calculations")
	}

	var reactions []ds.Reaction
	var reactionCalculations []ds.ReactionCalculation
	sub := r.db.Table("reaction_calculations").Where("calculation_id = ?", calculation.ID).Find(&reactionCalculations)
	err = r.db.Where("id IN (?)", sub.Select("reaction_id")).Find(&reactions).Error

	var reactionInfos []ds.ReactionInfo
	for _, reaction := range reactions {
		for _, reactionCalculation := range reactionCalculations {
			if reaction.ID == reactionCalculation.ReactionID {
				reactionInfos = append(reactionInfos, ds.ReactionInfo{
					ID:                 reaction.ID,
					Title:              reaction.Title,
					Reagent:            reaction.Reagent,
					Product:            reaction.Product,
					ConversationFactor: reaction.ConversationFactor,
					ImgLink:            reaction.ImgLink,

					OutputMass: reactionCalculation.OutputMass,
					InputMass:  reactionCalculation.InputMass,
				})
				break
			}
		}
	}

	if err != nil {
		return []ds.ReactionInfo{}, ds.Calculation{}, err
	}

	return reactionInfos, calculation, nil
}