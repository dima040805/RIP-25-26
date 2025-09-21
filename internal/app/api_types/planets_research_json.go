package apitypes

import "LAB1/internal/app/ds"

type PlanetsResearchJSON struct {
	ID           uint    	`json:"id"`
	ResearchID   uint     	`json:"research_id"`
	PlanetID     uint    	`json:"planet_id"`
	PlanetShine  float64 	`json:"planet_shine"`                  
	PlanetRadius int 		`json:"planet_radius"`                   
}

func PlanetsResearchToJSON(planetsResearch ds.PlanetsResearch) PlanetsResearchJSON {
	return PlanetsResearchJSON{
		ID:            planetsResearch.ID,
		ResearchID:    planetsResearch.ResearchID,
		PlanetID: planetsResearch.PlanetID,
		PlanetShine:    planetsResearch.PlanetShine,
		PlanetRadius:     planetsResearch.PlanetRadius,
	}
}

func PlanetsResearchFromJSON(planetsResearchJSON PlanetsResearchJSON) ds.PlanetsResearch {
	return ds.PlanetsResearch{
		PlanetShine:     planetsResearchJSON.PlanetShine,
	}
}