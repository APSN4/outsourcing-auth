package repository

import (
	"core/internal/database"
	"errors"
	"fmt"
	"gorm.io/gorm"
)

type CardRepository interface {
	Create(card *database.Card) error
	GetByID(id uint) (*database.Card, error)
	GetByCompanyID(companyID uint, limit, offset int) ([]database.Card, error)
	GetAll(limit, offset int) ([]database.Card, error)
	GetByCategory(category string, limit, offset int) ([]database.Card, error)
	Update(card *database.Card) error
	Delete(id uint) error
	SearchByTitle(query string, limit, offset int) ([]database.Card, error)
	GetByPriceRange(minPrice, maxPrice float64, limit, offset int) ([]database.Card, error)
}

type cardRepository struct {
	db *gorm.DB
}

func (r *cardRepository) Create(card *database.Card) error {
	return r.db.Create(card).Error
}

func (r *cardRepository) GetByID(id uint) (*database.Card, error) {
	var card database.Card
	err := r.db.Preload("Company").First(&card, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("card with ID %d not found", id)
		}
		return nil, err
	}
	return &card, nil
}

func (r *cardRepository) GetByCompanyID(companyID uint, limit, offset int) ([]database.Card, error) {
	var cards []database.Card
	err := r.db.Preload("Company").Where("company_id = ? AND is_active = true", companyID).
		Limit(limit).Offset(offset).Order("created_at DESC").Find(&cards).Error
	return cards, err
}

func (r *cardRepository) GetAll(limit, offset int) ([]database.Card, error) {
	var cards []database.Card
	err := r.db.Preload("Company").Where("is_active = true").
		Limit(limit).Offset(offset).Order("created_at DESC").Find(&cards).Error
	return cards, err
}

func (r *cardRepository) GetByCategory(category string, limit, offset int) ([]database.Card, error) {
	var cards []database.Card
	err := r.db.Preload("Company").Where("category = ? AND is_active = true", category).
		Limit(limit).Offset(offset).Order("created_at DESC").Find(&cards).Error
	return cards, err
}

func (r *cardRepository) Update(card *database.Card) error {
	return r.db.Save(card).Error
}

func (r *cardRepository) Delete(id uint) error {
	return r.db.Model(&database.Card{}).Where("id = ?", id).Update("is_active", false).Error
}

func (r *cardRepository) SearchByTitle(query string, limit, offset int) ([]database.Card, error) {
	var cards []database.Card
	searchQuery := "%" + query + "%"
	err := r.db.Preload("Company").
		Where("(title ILIKE ? OR description ILIKE ?) AND is_active = true", searchQuery, searchQuery).
		Limit(limit).Offset(offset).Order("created_at DESC").Find(&cards).Error
	return cards, err
}

func (r *cardRepository) GetByPriceRange(minPrice, maxPrice float64, limit, offset int) ([]database.Card, error) {
	var cards []database.Card
	err := r.db.Preload("Company").Where("price BETWEEN ? AND ? AND is_active = true", minPrice, maxPrice).
		Limit(limit).Offset(offset).Order("price ASC").Find(&cards).Error
	return cards, err
}

func NewCardRepository(db *gorm.DB) CardRepository {
	return &cardRepository{db: db}
}
