package repository

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
	userId int
}

func NewRepository(dsn string) (*Repository, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{}) // подключаемся к БД
	if err != nil {
		return nil, err
	}

	return &Repository{
		db: db,
		userId: 0,
	}, nil
}


func (r *Repository) GetUserID() (int) {
	return r.userId
}

func (r *Repository) SetUserID(id int) {
	r.userId = id
}

func (r *Repository) SignOut() {
	r.userId = 0
}