package apitypes

import "LAB1/internal/app/ds"

type PlanetJSON struct {
	ID          int 	`json:"id"`
	Name        string 	`json:"name"`
	Image       string 	`json:"image"`
	Distance    int		`json:"distance"`
	Description string 	`json:"description"`
	Mass        float64 `json:"mass"`
	Discovery   int		`json:"discovery"`
	StarRadius  int		`json:"star_radius"`
	IsDelete    bool   	`json:"is_delete"`
}

func PlanetToJSON(planet ds.Planet) PlanetJSON {
	return PlanetJSON{
		ID:			 planet.ID,
		IsDelete:  	 planet.IsDelete,  
		Image:       planet.Image,
		Name:        planet.Name,
		Distance:    planet.Distance,
		Description: planet.Description,
		Mass:        planet.Mass,
		Discovery:   planet.Discovery,
		StarRadius:  planet.StarRadius,
	}
}

func PlanetFromJSON(planetJSON PlanetJSON) ds.Planet {
	return ds.Planet{
		IsDelete:  	 planetJSON.IsDelete,  
		Name:        planetJSON.Name,
		Distance:    planetJSON.Distance,
		Description: planetJSON.Description,
		Mass:        planetJSON.Mass,
		Discovery:   planetJSON.Discovery,
		StarRadius:  planetJSON.StarRadius,
	}
}