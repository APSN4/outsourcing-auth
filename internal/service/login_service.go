package service

import (
	"core/internal"
	"core/internal/api"
	"core/internal/database"
	"core/internal/database/repository"
	"core/internal/security"
	"errors"
)

type LoginService interface {
	Login(request *api.UserCompanyRegister) (dbUser database.CompanyDB, jwtToken string, err error)
	AccessByToken(request *api.TokenAccess) (*api.ResponseSuccessAccess, database.CompanyDB, error)
}

type loginService struct {
	repository repository.LoginRepository
}

func (service *loginService) Login(request *api.GeneralAuth) (dbUser database.CompanyDB, jwtToken string, err error) {
	dbUser, err = service.repository.CheckPassword(request.GeneralLogin.GeneralLoginAttributes.Email, request.GeneralLogin.GeneralLoginAttributes.PasswordHash)
	if err != nil {
		return dbUser, "", err
	}
	jwtToken = security.CreateToken(dbUser.Type == "company", dbUser.ID, internal.LifeTimeJWT)
	if jwtToken == "" {
		return dbUser, "", errors.New("the created jwt was faulty")
	}
	return dbUser, jwtToken, nil
}

func (service *loginService) AccessByToken(request *api.TokenAccess) (*api.ResponseSuccessAccess, database.CompanyDB, error) {
	result, tokenStructure := security.CheckToken(request.User.Login.Token)
	company, err := service.GetCompany(uint(tokenStructure["accessID"].(float64)))
	if err != nil {
		return nil, company, err
	}

	if result {
		response := api.ResponseSuccessAccess{
			StatusResponse: &internal.StatusResponse{Status: "success"},
			ResponseUser: &api.ResponseUser{
				ID:    company.ID,
				Token: request.User.Login.Token,
				Type:  company.Type,
			},
		}
		return &response, company, nil
	} else {
		response := api.ResponseSuccessAccess{
			StatusResponse: &internal.StatusResponse{Status: "success"},
			ResponseUser: &api.ResponseUser{
				ID:    company.ID,
				Token: "expired",
				Type:  company.Type,
			},
		}
		return &response, company, nil
	}
}

func NewLoginService(repository repository.LoginRepository) LoginService {
	return &loginService{
		repository: repository,
	}
}
