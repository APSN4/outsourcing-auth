package repository

import (
	"core/internal/database"
	"errors"
	"fmt"
	"gorm.io/gorm"
)

type CompanyRepository interface {
	Save(company *database.CompanyDB)
	ExistsByEmail(email string) (bool, database.CompanyDB)
	GetByID(id uint) (*database.CompanyDB, error)
}

type companyRepository struct {
	db *gorm.DB
}

func (repository *companyRepository) Save(company *database.CompanyDB) {
	repository.db.Save(company)
}

func (repository *companyRepository) ExistsByEmail(email string) (bool, database.CompanyDB) {
	var company database.CompanyDB
	result := repository.db.Model(&database.CompanyDB{}).Where("email = ?", email).First(&company)
	exists := !errors.Is(result.Error, gorm.ErrRecordNotFound)
	return exists, company
}

func (repository *companyRepository) GetByID(id uint) (*database.CompanyDB, error) {
	var company database.CompanyDB

	result := repository.db.Model(&database.CompanyDB{}).Where("id = ?", id).First(&company)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("company with ID %d not found", id)
		}
		return nil, result.Error
	}
	return &company, nil
}

func NewCompanyRepository(db *gorm.DB) CompanyRepository {
	var repository CompanyRepository

	repository = &companyRepository{
		db: db,
	}

	return repository
}
