package service

import (
	"core/internal/database"
	"core/internal/database/repository"
	"errors"
)

type ReviewService interface {
	CreateReview(clientID, companyID, orderID uint, rating int, comment string) (*database.Review, error)
	GetReviewsByCompany(companyID uint, page, limit int) ([]database.Review, error)
	GetReviewByOrder(orderID uint) (*database.Review, error)
	GetCompanyRating(companyID uint) (float64, int, error)
}

type reviewService struct {
	reviewRepo repository.ReviewRepository
	orderRepo  repository.OrderRepository
}

func (s *reviewService) CreateReview(clientID, companyID, orderID uint, rating int, comment string) (*database.Review, error) {
	if rating < 1 || rating > 5 {
		return nil, errors.New("rating must be between 1 and 5")
	}

	// Проверяем, что заказ существует и принадлежит клиенту
	order, err := s.orderRepo.GetByID(orderID)
	if err != nil {
		return nil, err
	}

	if order.ClientID != clientID {
		return nil, errors.New("unauthorized: order does not belong to this client")
	}

	if order.CompanyID != companyID {
		return nil, errors.New("order does not belong to this company")
	}

	if order.Status != "finished" {
		return nil, errors.New("can only review finished orders")
	}

	// Проверяем, что отзыв еще не был оставлен
	existingReview, _ := s.reviewRepo.GetByOrderID(orderID)
	if existingReview != nil {
		return nil, errors.New("review already exists for this order")
	}

	review := &database.Review{
		ClientID:  clientID,
		CompanyID: companyID,
		OrderID:   orderID,
		Rating:    rating,
		Comment:   comment,
	}

	err = s.reviewRepo.Create(review)
	if err != nil {
		return nil, err
	}

	return review, nil
}

func (s *reviewService) GetReviewsByCompany(companyID uint, page, limit int) ([]database.Review, error) {
	offset := (page - 1) * limit
	return s.reviewRepo.GetByCompanyID(companyID, limit, offset)
}

func (s *reviewService) GetReviewByOrder(orderID uint) (*database.Review, error) {
	return s.reviewRepo.GetByOrderID(orderID)
}

func (s *reviewService) GetCompanyRating(companyID uint) (float64, int, error) {
	return s.reviewRepo.GetCompanyAverageRating(companyID)
}

func NewReviewService(reviewRepo repository.ReviewRepository, orderRepo repository.OrderRepository) ReviewService {
	return &reviewService{
		reviewRepo: reviewRepo,
		orderRepo:  orderRepo,
	}
}
