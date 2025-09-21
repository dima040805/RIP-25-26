package apitypes

import "LAB1/internal/app/ds"

type PlanetJSON struct {
	ID          int 	`json:"id"`
	IsDelete    bool   	`json:"is_delete"`
	Image       string 	`json:"image"`
	Name        string 	`json:"name"`
	Distance    int		`json:"distance"`
	Description string 	`json:"description"`
	Mass        float64 `json:"mass"`
	Discovery   int		`json:"discovery"`
	StarRadius  int		`json:"star_radius"`
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
		Image:       planetJSON.Image,
		Name:        planetJSON.Description,
		Distance:    planetJSON.Distance,
		Description: planetJSON.Description,
		Mass:        planetJSON.Mass,
		Discovery:   planetJSON.Discovery,
		StarRadius:  planetJSON.StarRadius,
	}
}