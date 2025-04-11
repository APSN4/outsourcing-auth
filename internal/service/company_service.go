package service

import (
	"core/internal"
	"core/internal/api"
	"core/internal/database"
	"core/internal/database/repository"
	"core/internal/security"
	"errors"
)

type CompanyService interface {
	Signup(request *api.UserCompanyRegister) (database.CompanyDB, error)
	GetCompany(id uint) (database.CompanyDB, error)
	Login(request *api.GeneralAuth) (dbUser database.CompanyDB, jwtToken string, err error)
	AccessByToken(request *api.TokenAccess) (*api.ResponseSuccessAccess, database.CompanyDB, error)
}

type companyService struct {
	repository repository.CompanyRepository
}

func (service *companyService) Signup(request *api.UserCompanyRegister) (database.CompanyDB, error) {
	exists, existsCompany, _ := service.repository.ExistsByEmail(request.Email)
	if exists || existsCompany {
		return database.CompanyDB{}, errors.New("email already exists")
	}

	company := &database.CompanyDB{
		CompanyName:   request.CompanyName,
		FullName:      request.FullName,
		PositionAgent: request.PositionAgent,
		IDCompany:     request.IDCompany,
		Email:         request.Email,
		Phone:         request.Phone,
		Address:       request.Address,
		TypeService:   request.TypeService,
		PasswordHash:  request.PasswordHash,
		Photo:         request.Photo,
		Documents:     request.Documents,
		Stars:         5,
		Type:          "company",
	}
	service.repository.Save(company)
	return *company, nil
}

func (service *companyService) GetCompany(id uint) (database.CompanyDB, error) {
	company, err := service.repository.GetByID(id)
	if err != nil {
		return database.CompanyDB{}, err
	}
	return *company, nil
}

func (service *companyService) Login(request *api.GeneralAuth) (dbUser database.CompanyDB, jwtToken string, err error) {
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

func (service *companyService) AccessByToken(request *api.TokenAccess) (*api.ResponseSuccessAccess, database.CompanyDB, error) {
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

func NewCompanyService(repository repository.CompanyRepository) CompanyService {
	return &companyService{
		repository: repository,
	}
}
