package service

import (
	"core/internal/api"
	"core/internal/database"
	"core/internal/database/repository"
	"errors"
)

type ClientService interface {
	Signup(request *api.ClientRegister) (database.ClientDB, error)
	GetClient(id uint) (database.ClientDB, error)
}

type clientService struct {
	repository repository.ClientRepository
}

func (service *clientService) Signup(request *api.ClientRegister) (database.ClientDB, error) {
	exists := service.repository.ExistsByEmail(request.Email)
	if exists {
		return database.ClientDB{}, errors.New("email already exists")
	}

	client := &database.ClientDB{
		FullName:     request.FullName,
		Email:        request.Email,
		Phone:        request.Phone,
		PasswordHash: request.PasswordHash,
		Photo:        request.Photo,
		Type:         request.Type,
	}
	service.repository.Save(client)
	return *client, nil
}

func (service *clientService) GetClient(id uint) (database.ClientDB, error) {
	client, err := service.repository.GetByID(id)
	if err != nil {
		return database.ClientDB{}, err
	}
	return *client, nil
}

func NewClientService(repository repository.ClientRepository) ClientService {
	return &clientService{
		repository: repository,
	}
}
