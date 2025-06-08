package repository

import (
	"core/internal/database"
	"gorm.io/gorm"
)

type ReviewRepository interface {
	Create(review *database.Review) error
	GetByCompanyID(companyID uint, limit, offset int) ([]database.Review, error)
	GetByOrderID(orderID uint) (*database.Review, error)
	UpdateCompanyRating(companyID uint) error
	GetCompanyAverageRating(companyID uint) (float64, int, error)
}

type reviewRepository struct {
	db *gorm.DB
}

func (r *reviewRepository) Create(review *database.Review) error {
	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Create(review).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := r.updateCompanyRatingInTx(tx, review.CompanyID); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (r *reviewRepository) GetByCompanyID(companyID uint, limit, offset int) ([]database.Review, error) {
	var reviews []database.Review
	err := r.db.Preload("Client").Where("company_id = ?", companyID).
		Limit(limit).Offset(offset).Order("created_at DESC").Find(&reviews).Error
	return reviews, err
}

func (r *reviewRepository) GetByOrderID(orderID uint) (*database.Review, error) {
	var review database.Review
	err := r.db.Preload("Client").Where("order_id = ?", orderID).First(&review).Error
	if err != nil {
		return nil, err
	}
	return &review, nil
}

func (r *reviewRepository) UpdateCompanyRating(companyID uint) error {
	return r.updateCompanyRatingInTx(r.db, companyID)
}

func (r *reviewRepository) updateCompanyRatingInTx(tx *gorm.DB, companyID uint) error {
	var result struct {
		AvgRating float64
		Count     int64
	}

	err := tx.Model(&database.Review{}).Select("AVG(rating) as avg_rating, COUNT(*) as count").
		Where("company_id = ?", companyID).Scan(&result).Error
	if err != nil {
		return err
	}

	return tx.Model(&database.CompanyDB{}).Where("id = ?", companyID).
		Updates(map[string]interface{}{
			"stars":        result.AvgRating,
			"review_count": result.Count,
		}).Error
}

func (r *reviewRepository) GetCompanyAverageRating(companyID uint) (float64, int, error) {
	var result struct {
		AvgRating float64
		Count     int64
	}

	err := r.db.Model(&database.Review{}).Select("AVG(rating) as avg_rating, COUNT(*) as count").
		Where("company_id = ?", companyID).Scan(&result).Error
	
	return result.AvgRating, int(result.Count), err
}

func NewReviewRepository(db *gorm.DB) ReviewRepository {
	return &reviewRepository{db: db}
}
