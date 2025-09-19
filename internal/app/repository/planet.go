package repository

import (
	"database/sql"
	"errors"
	"fmt"

	"LAB1/internal/app/ds"

	"github.com/sirupsen/logrus"
)

func (r *Repository) GetPlanets() ([]ds.Planet, error) {
	var planets []ds.Planet
	err := r.db.Order("id").Find(&planets).Error
	if err != nil {
		return nil, err
	}
	if len(planets) == 0 {
		return nil, fmt.Errorf("массив пустой")
	}

	return planets, nil
}

func (r *Repository) GetPlanet(id int) (*ds.Planet, error) {
	query := "SELECT id, image, name, description, distance, mass, discovery, star_radius FROM planets WHERE id = $1"
	row := r.db.Raw(query, id).Row()
	planet := &ds.Planet{}

   err := row.Scan(
		&planet.ID,
      &planet.Image,
      &planet.Name,
      &planet.Description,
      &planet.Distance,
      &planet.Mass,
	  &planet.Discovery,
      &planet.StarRadius,
   )

   if err != nil {
      if errors.Is(err, sql.ErrNoRows) {
         return nil, nil // Возвращаем nil, если записи нет
      }
      return nil, err
   }
	return planet, nil
}

func (r *Repository) GetPlanetsByName(name string) ([]ds.Planet, error) {
	var planets []ds.Planet
	err := r.db.Where("name ILIKE ?", "%"+name+"%").Find(&planets).Error
	if err != nil {
		return nil, err
	}
	return planets, nil
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
	fmt.Println(researchID)
	err = r.db.Model(&ds.PlanetsResearch{}).Where("research_id = ?", researchID).Count(&count).Error
	if err != nil {
		logrus.Println("Error counting records in lists_planets:", err)
	}

	return count
}
