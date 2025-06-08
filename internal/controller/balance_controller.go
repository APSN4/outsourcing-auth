package controller

import (
	"core/internal/api"
	"core/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type BalanceController interface {
	GetClientBalance(c *gin.Context, request *api.TokenAccess)
	GetCompanyBalance(c *gin.Context, request *api.TokenAccess)
	DepositClientBalance(c *gin.Context, request *api.TokenDepositBalance)
	WithdrawCompanyBalance(c *gin.Context, request *api.TokenWithdrawBalance)
	GetClientTransactions(c *gin.Context, request *api.TokenAccess)
	GetCompanyTransactions(c *gin.Context, request *api.TokenAccess)
}

type balanceController struct {
	balanceService service.BalanceService
}

func (ctrl *balanceController) GetClientBalance(c *gin.Context, request *api.TokenAccess) {
	userInfo, err := ExtractUserFromToken(request.User.Login.Token)
	if err != nil {
		api.GetErrorJSON(c, http.StatusUnauthorized, err.Error())
		return
	}

	if userInfo.IsCompany {
		api.GetErrorJSON(c, http.StatusForbidden, "Only clients can access client balance")
		return
	}

	balance, err := ctrl.balanceService.GetClientBalance(userInfo.UserID)
	if err != nil {
		api.GetErrorJSON(c, http.StatusInternalServerError, "Failed to get balance")
		return
	}

	c.JSON(http.StatusOK, gin.H{"balance": balance})
}

func (ctrl *balanceController) GetCompanyBalance(c *gin.Context, request *api.TokenAccess) {
	userInfo, err := ExtractUserFromToken(request.User.Login.Token)
	if err != nil {
		api.GetErrorJSON(c, http.StatusUnauthorized, err.Error())
		return
	}

	if !userInfo.IsCompany {
		api.GetErrorJSON(c, http.StatusForbidden, "Only companies can access company balance")
		return
	}

	balance, err := ctrl.balanceService.GetCompanyBalance(userInfo.UserID)
	if err != nil {
		api.GetErrorJSON(c, http.StatusInternalServerError, "Failed to get balance")
		return
	}

	c.JSON(http.StatusOK, gin.H{"balance": balance})
}

func (ctrl *balanceController) DepositClientBalance(c *gin.Context, request *api.TokenDepositBalance) {
	userInfo, err := ExtractUserFromToken(request.User.Login.Token)
	if err != nil {
		api.GetErrorJSON(c, http.StatusUnauthorized, err.Error())
		return
	}

	if userInfo.IsCompany {
		api.GetErrorJSON(c, http.StatusForbidden, "Only clients can deposit to client balance")
		return
	}

	if request.Amount <= 0 {
		api.GetErrorJSON(c, http.StatusBadRequest, "Amount must be greater than 0")
		return
	}

	err = ctrl.balanceService.DepositClientBalance(userInfo.UserID, request.Amount)
	if err != nil {
		api.GetErrorJSON(c, http.StatusBadRequest, err.Error())
		return
	}

	// Получаем новый баланс
	newBalance, err := ctrl.balanceService.GetClientBalance(userInfo.UserID)
	if err != nil {
		api.GetErrorJSON(c, http.StatusInternalServerError, "Failed to get updated balance")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Balance deposited successfully",
		"balance": newBalance,
		"amount":  request.Amount,
	})
}

func (ctrl *balanceController) WithdrawCompanyBalance(c *gin.Context, request *api.TokenWithdrawBalance) {
	userInfo, err := ExtractUserFromToken(request.User.Login.Token)
	if err != nil {
		api.GetErrorJSON(c, http.StatusUnauthorized, err.Error())
		return
	}

	if !userInfo.IsCompany {
		api.GetErrorJSON(c, http.StatusForbidden, "Only companies can withdraw from company balance")
		return
	}

	if request.Amount <= 0 {
		api.GetErrorJSON(c, http.StatusBadRequest, "Amount must be greater than 0")
		return
	}

	err = ctrl.balanceService.WithdrawCompanyBalance(userInfo.UserID, request.Amount)
	if err != nil {
		api.GetErrorJSON(c, http.StatusBadRequest, err.Error())
		return
	}

	// Получаем новый баланс
	newBalance, err := ctrl.balanceService.GetCompanyBalance(userInfo.UserID)
	if err != nil {
		api.GetErrorJSON(c, http.StatusInternalServerError, "Failed to get updated balance")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Balance withdrawn successfully",
		"balance": newBalance,
		"amount":  request.Amount,
	})
}

func (ctrl *balanceController) GetClientTransactions(c *gin.Context, request *api.TokenAccess) {
	userInfo, err := ExtractUserFromToken(request.User.Login.Token)
	if err != nil {
		api.GetErrorJSON(c, http.StatusUnauthorized, err.Error())
		return
	}

	if userInfo.IsCompany {
		api.GetErrorJSON(c, http.StatusForbidden, "Only clients can access client transactions")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	transactions, err := ctrl.balanceService.GetClientTransactions(userInfo.UserID, page, limit)
	if err != nil {
		api.GetErrorJSON(c, http.StatusInternalServerError, "Failed to get transactions")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"transactions": transactions,
		"page":         page,
		"limit":        limit,
	})
}

func (ctrl *balanceController) GetCompanyTransactions(c *gin.Context, request *api.TokenAccess) {
	userInfo, err := ExtractUserFromToken(request.User.Login.Token)
	if err != nil {
		api.GetErrorJSON(c, http.StatusUnauthorized, err.Error())
		return
	}

	if !userInfo.IsCompany {
		api.GetErrorJSON(c, http.StatusForbidden, "Only companies can access company transactions")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	transactions, err := ctrl.balanceService.GetCompanyTransactions(userInfo.UserID, page, limit)
	if err != nil {
		api.GetErrorJSON(c, http.StatusInternalServerError, "Failed to get transactions")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"transactions": transactions,
		"page":         page,
		"limit":        limit,
	})
}

func NewBalanceController(balanceService service.BalanceService) BalanceController {
	return &balanceController{balanceService: balanceService}
}
