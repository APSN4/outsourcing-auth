package controller

import (
	"core/internal/security"
	"errors"
	"github.com/golang-jwt/jwt/v5"
)

// UserInfo представляет информацию о пользователе из токена
type UserInfo struct {
	UserID    uint
	IsCompany bool
	UserType  string // "client" или "company"
}

// ExtractUserFromToken извлекает информацию о пользователе из токена
func ExtractUserFromToken(token string) (*UserInfo, error) {
	isValid, claims := security.CheckToken(token)
	if !isValid || claims == nil {
		return nil, errors.New("invalid or expired token")
	}

	// Извлекаем accessID
	accessIDFloat, ok := claims["accessID"].(float64)
	if !ok {
		return nil, errors.New("invalid accessID in token")
	}
	accessID := uint(accessIDFloat)

	// Извлекаем isCompany
	isCompany, ok := claims["isCompany"].(bool)
	if !ok {
		return nil, errors.New("invalid isCompany in token")
	}

	userType := "client"
	if isCompany {
		userType = "company"
	}

	return &UserInfo{
		UserID:    accessID,
		IsCompany: isCompany,
		UserType:  userType,
	}, nil
}

// GetUserIDFromClaims получает ID пользователя из JWT claims
func GetUserIDFromClaims(claims jwt.MapClaims) (uint, error) {
	accessIDFloat, ok := claims["accessID"].(float64)
	if !ok {
		return 0, errors.New("invalid accessID in token")
	}
	return uint(accessIDFloat), nil
}

// IsCompanyFromClaims проверяет, является ли пользователь компанией
func IsCompanyFromClaims(claims jwt.MapClaims) (bool, error) {
	isCompany, ok := claims["isCompany"].(bool)
	if !ok {
		return false, errors.New("invalid isCompany in token")
	}
	return isCompany, nil
}
