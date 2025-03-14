package repository

import (
	"core/internal/database"
	"errors"
	"gorm.io/gorm"
)

type ClientRepository interface {
	Save(client *database.ClientDB)
	ExistsByEmail(email string) bool
}

type clientRepository struct {
	db *gorm.DB
}

func (repository *clientRepository) Save(client *database.ClientDB) {
	repository.db.Save(client)
}

func (repository *clientRepository) ExistsByEmail(email string) bool {
	var client database.ClientDB
	result := repository.db.Model(&database.ClientDB{}).Where("email = ?", email).First(&client)
	exists := !errors.Is(result.Error, gorm.ErrRecordNotFound)
	return exists
}

func NewClientRepository(db *gorm.DB) ClientRepository {
	var repository ClientRepository

	repository = &clientRepository{
		db: db,
	}

	return repository
}
