package service

import (
	"core/internal/database"
	"core/internal/database/repository"
	"errors"
)

type CardService interface {
	CreateCard(companyID uint, title, description, category, location string, price float64) (*database.Card, error)
	GetCardByID(id uint) (*database.Card, error)
	GetCardsByCompany(companyID uint, page, limit int) ([]database.Card, error)
	GetAllCards(page, limit int) ([]database.Card, error)
	GetCardsByCategory(category string, page, limit int) ([]database.Card, error)
	UpdateCard(cardID, companyID uint, title, description, category, location string, price float64) (*database.Card, error)
	DeleteCard(cardID, companyID uint) error
	SearchCards(query string, page, limit int) ([]database.Card, error)
	GetCardsByPriceRange(minPrice, maxPrice float64, page, limit int) ([]database.Card, error)
}

type cardService struct {
	cardRepo repository.CardRepository
}

func (s *cardService) CreateCard(companyID uint, title, description, category, location string, price float64) (*database.Card, error) {
	if title == "" {
		return nil, errors.New("title cannot be empty")
	}
	if description == "" {
		return nil, errors.New("description cannot be empty")
	}
	if price <= 0 {
		return nil, errors.New("price must be greater than 0")
	}

	card := &database.Card{
		Title:       title,
		Description: description,
		Category:    category,
		Location:    location,
		Price:       price,
		CompanyID:   companyID,
		IsActive:    true,
	}

	err := s.cardRepo.Create(card)
	if err != nil {
		return nil, err
	}

	return s.cardRepo.GetByID(card.ID)
}

func (s *cardService) GetCardByID(id uint) (*database.Card, error) {
	return s.cardRepo.GetByID(id)
}

func (s *cardService) GetCardsByCompany(companyID uint, page, limit int) ([]database.Card, error) {
	offset := (page - 1) * limit
	return s.cardRepo.GetByCompanyID(companyID, limit, offset)
}

func (s *cardService) GetAllCards(page, limit int) ([]database.Card, error) {
	offset := (page - 1) * limit
	return s.cardRepo.GetAll(limit, offset)
}

func (s *cardService) GetCardsByCategory(category string, page, limit int) ([]database.Card, error) {
	offset := (page - 1) * limit
	return s.cardRepo.GetByCategory(category, limit, offset)
}

func (s *cardService) UpdateCard(cardID, companyID uint, title, description, category, location string, price float64) (*database.Card, error) {
	card, err := s.cardRepo.GetByID(cardID)
	if err != nil {
		return nil, err
	}

	if card.CompanyID != companyID {
		return nil, errors.New("unauthorized: card does not belong to this company")
	}

	if title != "" {
		card.Title = title
	}
	if description != "" {
		card.Description = description
	}
	if category != "" {
		card.Category = category
	}
	if location != "" {
		card.Location = location
	}
	if price > 0 {
		card.Price = price
	}

	err = s.cardRepo.Update(card)
	if err != nil {
		return nil, err
	}

	return card, nil
}

func (s *cardService) DeleteCard(cardID, companyID uint) error {
	card, err := s.cardRepo.GetByID(cardID)
	if err != nil {
		return err
	}

	if card.CompanyID != companyID {
		return errors.New("unauthorized: card does not belong to this company")
	}

	return s.cardRepo.Delete(cardID)
}

func (s *cardService) SearchCards(query string, page, limit int) ([]database.Card, error) {
	offset := (page - 1) * limit
	return s.cardRepo.SearchByTitle(query, limit, offset)
}

func (s *cardService) GetCardsByPriceRange(minPrice, maxPrice float64, page, limit int) ([]database.Card, error) {
	offset := (page - 1) * limit
	return s.cardRepo.GetByPriceRange(minPrice, maxPrice, limit, offset)
}

func NewCardService(cardRepo repository.CardRepository) CardService {
	return &cardService{cardRepo: cardRepo}
}
