package main

import (
	"github.com/dima040805/RIP-25-26/internal/app/ds"
	"github.com/dima040805/RIP-25-26/internal/app/dsn"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	_ = godotenv.Load()
	db, err := gorm.Open(postgres.Open(dsn.FromEnv()), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	err = db.AutoMigrate(
		&ds.Planet{},
		&ds.PlanetsResearch{},
		&ds.Research{},
		&ds.User{},
	)
	if err != nil {
		panic("cant migrate db")
	}
}
