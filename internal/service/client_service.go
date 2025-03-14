package service

import (
	"core/internal/api"
	"core/internal/database"
	"core/internal/database/repository"
	"errors"
)

type ClientService interface {
	Signup(request *api.ClientRegister) error
}

type clientService struct {
	repository repository.ClientRepository
}

func (service *clientService) Signup(request *api.ClientRegister) error {
	exists := service.repository.ExistsByEmail(request.Email)
	if exists {
		return errors.New("email already exists")
	}
	service.repository.Save(&database.ClientDB{
		FullName:     request.FullName,
		Email:        request.Email,
		Phone:        request.Phone,
		PasswordHash: request.PasswordHash,
		Photo:        request.Photo,
		Type:         request.Type,
	})
	return nil
}

func NewClientService(repository repository.ClientRepository) ClientService {
	return &clientService{
		repository: repository,
	}
}
