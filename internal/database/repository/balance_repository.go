package repository

import (
	"core/internal/database"
	"gorm.io/gorm"
)

type BalanceRepository interface {
	GetClientBalance(clientID uint) (float64, error)
	GetCompanyBalance(companyID uint) (float64, error)
	UpdateClientBalance(clientID uint, amount float64) error
	UpdateCompanyBalance(companyID uint, amount float64) error
	CreateTransaction(transaction *database.BalanceTransaction) error
	GetTransactionsByUser(userID uint, userType string, limit, offset int) ([]database.BalanceTransaction, error)
}

type balanceRepository struct {
	db *gorm.DB
}

func (r *balanceRepository) GetClientBalance(clientID uint) (float64, error) {
	var client database.ClientDB
	err := r.db.Select("balance").First(&client, clientID).Error
	if err != nil {
		return 0, err
	}
	return client.Balance, nil
}

func (r *balanceRepository) GetCompanyBalance(companyID uint) (float64, error) {
	var company database.CompanyDB
	err := r.db.Select("balance").First(&company, companyID).Error
	if err != nil {
		return 0, err
	}
	return company.Balance, nil
}

func (r *balanceRepository) UpdateClientBalance(clientID uint, amount float64) error {
	return r.db.Model(&database.ClientDB{}).Where("id = ?", clientID).
		Update("balance", gorm.Expr("balance + ?", amount)).Error
}

func (r *balanceRepository) UpdateCompanyBalance(companyID uint, amount float64) error {
	return r.db.Model(&database.CompanyDB{}).Where("id = ?", companyID).
		Update("balance", gorm.Expr("balance + ?", amount)).Error
}

func (r *balanceRepository) CreateTransaction(transaction *database.BalanceTransaction) error {
	return r.db.Create(transaction).Error
}

func (r *balanceRepository) GetTransactionsByUser(userID uint, userType string, limit, offset int) ([]database.BalanceTransaction, error) {
	var transactions []database.BalanceTransaction
	err := r.db.Where("user_id = ? AND user_type = ?", userID, userType).
		Limit(limit).Offset(offset).Order("created_at DESC").Find(&transactions).Error
	return transactions, err
}

func NewBalanceRepository(db *gorm.DB) BalanceRepository {
	return &balanceRepository{db: db}
}
