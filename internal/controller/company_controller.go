package controller

import (
	"core/internal"
	"core/internal/api"
	"core/internal/security"
	"core/internal/service"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
)

type CompanyController interface {
	Signup(c *gin.Context)
	Login(c *gin.Context, request *api.GeneralAuth)
	GetAccount(c *gin.Context, request *api.TokenAccess)
}

type companyController struct {
	service service.CompanyService
}

func (controller companyController) Signup(c *gin.Context) {
	request := &api.UserCompanyRegister{}
	if err := c.ShouldBind(request); err != nil && errors.As(err, &validator.ValidationErrors{}) {
		api.GetErrorJSON(c, http.StatusBadRequest, "JSON is invalid")
		return
	}
	company, err := controller.service.Signup(request)
	if err != nil {
		api.GetErrorJSON(c, http.StatusPreconditionFailed, err.Error())
		return
	}
	tokenGenerated := security.CreateToken(true, company.ID, 60)
	c.JSON(http.StatusOK, api.ResponseSuccessAccess{
		StatusResponse: &internal.StatusResponse{Status: "success"},
		ResponseUser: &api.ResponseUser{
			ID:    company.ID,
			Token: tokenGenerated,
			Type:  company.Type,
		},
	})
}

func (controller companyController) Login(c *gin.Context, request *api.GeneralAuth) {
	dbUser, jwtToken, err := controller.service.Login(request)
	if err != nil {
		api.GetErrorJSON(c, http.StatusBadRequest, "the created jwt was faulty")
		return
	}
	c.JSON(http.StatusOK, api.ResponseSuccessAccess{
		StatusResponse: &internal.StatusResponse{Status: "success"},
		ResponseUser: &api.ResponseUser{
			ID:    dbUser.ID,
			Token: jwtToken,
			Type:  dbUser.Type,
		},
	})
}

func (controller companyController) GetAccount(c *gin.Context, request *api.TokenAccess) {
	response, user, err := controller.service.AccessByToken(request)
	if err != nil {
		api.GetErrorJSON(c, http.StatusBadRequest, "The token is incorrect")
		return
	}
	if response.ResponseUser.Token == "expired" {
		api.GetErrorJSON(c, http.StatusForbidden, "The token has expired")
		return
	}
	c.JSON(http.StatusOK, api.ResponseAccountCompany{
		StatusResponse: &internal.StatusResponse{Status: "success"},
		User: struct {
			Account api.CompanyInfo `json:"account"`
		}{Account: api.CompanyInfo{
			ID:            user.ID,
			CompanyName:   user.CompanyName,
			Email:         user.Email,
			Phone:         user.Phone,
			FullName:      user.FullName,
			PositionAgent: user.PositionAgent,
			IDCompany:     user.IDCompany,
			Address:       user.Address,
			TypeService:   user.TypeService,
			PasswordHash:  user.PasswordHash,
			Photo:         user.Photo,
			Documents:     user.Documents,
			Type:          user.Type,
		}},
	})
}

func NewCompanyController(service service.CompanyService) CompanyController {
	return &companyController{
		service: service,
	}
}
