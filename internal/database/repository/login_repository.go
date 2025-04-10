package repository

import (
	"core/internal/database"
	"errors"
	"gorm.io/gorm"
)

type LoginRepository interface {
	CheckPassword(email string, PasswordHash string) (database.CompanyDB, error)
	ExistsByEmail(email string) (bool, any)
}

func (repository *loginRepository) CheckPassword(entity any, email string, PasswordHash string) (database.CompanyDB, error) {
	ok, dbUser := ExistsByEmail[entity](repository.db, email)
	if ok {
		if dbUser.PasswordHash == PasswordHash {
			return dbUser, nil
		} else {
			return database.CompanyDB{}, errors.New("bad password hash")
		}
	} else {
		return database.CompanyDB{}, errors.New("user not found")
	}
}

func ExistsByEmail[T any](db *gorm.DB, email string) (bool, T) {
	var entity T
	result := db.Where("email = ?", email).First(&entity)
	exists := result.Error == nil
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return false, entity
	}
	return exists, entity
}

type loginRepository struct {
	db *gorm.DB
}
