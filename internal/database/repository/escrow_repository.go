package repository

import (
	"core/internal/database"
	"crypto/rand"
	"encoding/hex"
	"time"
	"gorm.io/gorm"
)

type EscrowRepository interface {
	CreateTransaction(transaction *database.EscrowTransaction) error
	GetByOrderID(orderID uint) ([]database.EscrowTransaction, error)
	UpdateTransactionStatus(id uint, status string) error
}

type WorkerLinkRepository interface {
	Create(link *database.WorkerLink) error
	GetByToken(token string) (*database.WorkerLink, error)
	MarkAsUsed(token string) error
	GenerateWorkerLink(orderID uint) (*database.WorkerLink, error)
}

type escrowRepository struct {
	db *gorm.DB
}

type workerLinkRepository struct {
	db *gorm.DB
}

func (r *escrowRepository) CreateTransaction(transaction *database.EscrowTransaction) error {
	return r.db.Create(transaction).Error
}

func (r *escrowRepository) GetByOrderID(orderID uint) ([]database.EscrowTransaction, error) {
	var transactions []database.EscrowTransaction
	err := r.db.Where("order_id = ?", orderID).Order("created_at DESC").Find(&transactions).Error
	return transactions, err
}

func (r *escrowRepository) UpdateTransactionStatus(id uint, status string) error {
	return r.db.Model(&database.EscrowTransaction{}).Where("id = ?", id).Update("status", status).Error
}

func (r *workerLinkRepository) Create(link *database.WorkerLink) error {
	return r.db.Create(link).Error
}

func (r *workerLinkRepository) GetByToken(token string) (*database.WorkerLink, error) {
	var link database.WorkerLink
	err := r.db.Where("token = ? AND is_used = false AND expires_at > ?", token, time.Now()).First(&link).Error
	if err != nil {
		return nil, err
	}
	return &link, nil
}

func (r *workerLinkRepository) MarkAsUsed(token string) error {
	return r.db.Model(&database.WorkerLink{}).Where("token = ?", token).Update("is_used", true).Error
}

func (r *workerLinkRepository) GenerateWorkerLink(orderID uint) (*database.WorkerLink, error) {
	// Генерируем случайный токен
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return nil, err
	}
	token := hex.EncodeToString(bytes)

	link := &database.WorkerLink{
		OrderID:   orderID,
		Token:     token,
		IsUsed:    false,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour), // Действует 7 дней
	}

	err := r.db.Create(link).Error
	if err != nil {
		return nil, err
	}

	return link, nil
}

func NewEscrowRepository(db *gorm.DB) EscrowRepository {
	return &escrowRepository{db: db}
}

func NewWorkerLinkRepository(db *gorm.DB) WorkerLinkRepository {
	return &workerLinkRepository{db: db}
}
