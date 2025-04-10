package repository

import (
	"core/internal/database"
	"errors"
	"fmt"
	"gorm.io/gorm"
)

type ClientRepository interface {
	Save(client *database.ClientDB)
	ExistsByEmail(email string) (bool, database.ClientDB)
	GetByID(id uint) (*database.ClientDB, error)
}

type clientRepository struct {
	db *gorm.DB
}

func (repository *clientRepository) Save(client *database.ClientDB) {
	repository.db.Save(client)
}

func (repository *clientRepository) ExistsByEmail(email string) (bool, database.ClientDB) {
	var client database.ClientDB
	result := repository.db.Model(&database.ClientDB{}).Where("email = ?", email).First(&client)
	exists := !errors.Is(result.Error, gorm.ErrRecordNotFound)
	return exists, client
}

func (repository *clientRepository) GetByID(id uint) (*database.ClientDB, error) {
	var client database.ClientDB

	result := repository.db.Model(&database.ClientDB{}).Where("id = ?", id).First(&client)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("client with ID %d not found", id)
		}
		return nil, result.Error
	}
	return &client, nil
}

func NewClientRepository(db *gorm.DB) ClientRepository {
	var repository ClientRepository

	repository = &clientRepository{
		db: db,
	}

	return repository
}
