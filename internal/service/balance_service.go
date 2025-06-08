package service

import (
	"core/internal/database"
	"core/internal/database/repository"
	"errors"
	"fmt"
)

type BalanceService interface {
	GetClientBalance(clientID uint) (float64, error)
	GetCompanyBalance(companyID uint) (float64, error)
	DepositClientBalance(clientID uint, amount float64) error
	WithdrawCompanyBalance(companyID uint, amount float64) error
	GetClientTransactions(clientID uint, page, limit int) ([]database.BalanceTransaction, error)
	GetCompanyTransactions(companyID uint, page, limit int) ([]database.BalanceTransaction, error)
}

type balanceService struct {
	balanceRepo repository.BalanceRepository
}

func (s *balanceService) GetClientBalance(clientID uint) (float64, error) {
	return s.balanceRepo.GetClientBalance(clientID)
}

func (s *balanceService) GetCompanyBalance(companyID uint) (float64, error) {
	return s.balanceRepo.GetCompanyBalance(companyID)
}

func (s *balanceService) DepositClientBalance(clientID uint, amount float64) error {
	if amount <= 0 {
		return errors.New("amount must be greater than 0")
	}

	// В реальной системе здесь была бы интеграция с платежным шлюзом
	// Пока просто добавляем деньги на счет
	err := s.balanceRepo.UpdateClientBalance(clientID, amount)
	if err != nil {
		return err
	}

	// Создаем транзакцию
	transaction := &database.BalanceTransaction{
		UserID:      clientID,
		UserType:    "client",
		Amount:      amount,
		Type:        "deposit",
		Status:      "completed",
		Description: fmt.Sprintf("Пополнение баланса на %.2f руб.", amount),
	}

	return s.balanceRepo.CreateTransaction(transaction)
}

func (s *balanceService) WithdrawCompanyBalance(companyID uint, amount float64) error {
	if amount <= 0 {
		return errors.New("amount must be greater than 0")
	}

	// Проверяем баланс
	balance, err := s.balanceRepo.GetCompanyBalance(companyID)
	if err != nil {
		return err
	}

	if balance < amount {
		return errors.New("insufficient balance")
	}

	// Списываем деньги
	err = s.balanceRepo.UpdateCompanyBalance(companyID, -amount)
	if err != nil {
		return err
	}

	// Создаем транзакцию
	transaction := &database.BalanceTransaction{
		UserID:      companyID,
		UserType:    "company",
		Amount:      -amount,
		Type:        "withdrawal",
		Status:      "completed",
		Description: fmt.Sprintf("Вывод средств %.2f руб.", amount),
	}

	return s.balanceRepo.CreateTransaction(transaction)
}

func (s *balanceService) GetClientTransactions(clientID uint, page, limit int) ([]database.BalanceTransaction, error) {
	offset := (page - 1) * limit
	return s.balanceRepo.GetTransactionsByUser(clientID, "client", limit, offset)
}

func (s *balanceService) GetCompanyTransactions(companyID uint, page, limit int) ([]database.BalanceTransaction, error) {
	offset := (page - 1) * limit
	return s.balanceRepo.GetTransactionsByUser(companyID, "company", limit, offset)
}

func NewBalanceService(balanceRepo repository.BalanceRepository) BalanceService {
	return &balanceService{balanceRepo: balanceRepo}
}
