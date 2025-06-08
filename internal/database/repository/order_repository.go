package repository

import (
	"core/internal/database"
	"errors"
	"fmt"
	"gorm.io/gorm"
)

type OrderRepository interface {
	Create(order *database.Order) error
	GetByID(id uint) (*database.Order, error)
	GetByClientID(clientID uint, limit, offset int) ([]database.Order, error)
	GetByCompanyID(companyID uint, limit, offset int) ([]database.Order, error)
	UpdateStatus(id uint, status string) error
	GetByWorkerToken(token string) (*database.Order, error)
	Update(order *database.Order) error
	GetAllActive(limit, offset int) ([]database.Order, error)
}

type orderRepository struct {
	db *gorm.DB
}

func (r *orderRepository) Create(order *database.Order) error {
	return r.db.Create(order).Error
}

func (r *orderRepository) GetByID(id uint) (*database.Order, error) {
	var order database.Order
	err := r.db.Preload("Client").Preload("Company").Preload("Card").First(&order, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("order with ID %d not found", id)
		}
		return nil, err
	}
	return &order, nil
}

func (r *orderRepository) GetByClientID(clientID uint, limit, offset int) ([]database.Order, error) {
	var orders []database.Order
	err := r.db.Preload("Company").Preload("Card").Where("client_id = ?", clientID).
		Limit(limit).Offset(offset).Order("created_at DESC").Find(&orders).Error
	return orders, err
}

func (r *orderRepository) GetByCompanyID(companyID uint, limit, offset int) ([]database.Order, error) {
	var orders []database.Order
	err := r.db.Preload("Client").Preload("Card").Where("company_id = ?", companyID).
		Limit(limit).Offset(offset).Order("created_at DESC").Find(&orders).Error
	return orders, err
}

func (r *orderRepository) UpdateStatus(id uint, status string) error {
	return r.db.Model(&database.Order{}).Where("id = ?", id).Update("status", status).Error
}

func (r *orderRepository) GetByWorkerToken(token string) (*database.Order, error) {
	var workerLink database.WorkerLink
	err := r.db.Where("token = ? AND is_used = false AND expires_at > NOW()", token).First(&workerLink).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("worker link not found or expired")
		}
		return nil, err
	}

	var order database.Order
	err = r.db.Preload("Client").Preload("Company").Preload("Card").First(&order, workerLink.OrderID).Error
	if err != nil {
		return nil, err
	}
	
	return &order, nil
}

func (r *orderRepository) Update(order *database.Order) error {
	return r.db.Save(order).Error
}

func (r *orderRepository) GetAllActive(limit, offset int) ([]database.Order, error) {
	var orders []database.Order
	err := r.db.Preload("Client").Preload("Company").Preload("Card").
		Where("status IN ?", []string{"created", "paid", "in_progress"}).
		Limit(limit).Offset(offset).Order("created_at DESC").Find(&orders).Error
	return orders, err
}

func NewOrderRepository(db *gorm.DB) OrderRepository {
	return &orderRepository{db: db}
}
