package controller

import (
	"core/internal/api"
	"core/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type CardController interface {
	GetAllCards(c *gin.Context)
	GetCardsByCategory(c *gin.Context)
	GetCardByID(c *gin.Context)
	SearchCards(c *gin.Context)
	GetCardsByPriceRange(c *gin.Context)
	CreateCard(c *gin.Context, request *api.TokenCreateCard)
	UpdateCard(c *gin.Context, request *api.TokenUpdateCard)
	DeleteCard(c *gin.Context, request *api.TokenDeleteCard)
	GetCompanyCards(c *gin.Context, request *api.TokenListCard)
}

type cardController struct {
	cardService service.CardService
}

func (ctrl *cardController) GetAllCards(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	cards, err := ctrl.cardService.GetAllCards(page, limit)
	if err != nil {
		api.GetErrorJSON(c, http.StatusInternalServerError, "Failed to get cards")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"cards": cards,
		"page":  page,
		"limit": limit,
	})
}

func (ctrl *cardController) GetCardsByCategory(c *gin.Context) {
	category := c.Param("category")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	cards, err := ctrl.cardService.GetCardsByCategory(category, page, limit)
	if err != nil {
		api.GetErrorJSON(c, http.StatusInternalServerError, "Failed to get cards")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"cards":    cards,
		"category": category,
		"page":     page,
		"limit":    limit,
	})
}

func (ctrl *cardController) GetCardByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		api.GetErrorJSON(c, http.StatusBadRequest, "Invalid card ID")
		return
	}

	card, err := ctrl.cardService.GetCardByID(uint(id))
	if err != nil {
		api.GetErrorJSON(c, http.StatusNotFound, "Card not found")
		return
	}

	c.JSON(http.StatusOK, gin.H{"card": card})
}

func (ctrl *cardController) SearchCards(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		api.GetErrorJSON(c, http.StatusBadRequest, "Search query is required")
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

	cards, err := ctrl.cardService.SearchCards(query, page, limit)
	if err != nil {
		api.GetErrorJSON(c, http.StatusInternalServerError, "Failed to search cards")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"cards": cards,
		"query": query,
		"page":  page,
		"limit": limit,
	})
}

func (ctrl *cardController) GetCardsByPriceRange(c *gin.Context) {
	minPriceStr := c.Query("min_price")
	maxPriceStr := c.Query("max_price")

	var minPrice, maxPrice float64
	var err error

	if minPriceStr != "" {
		minPrice, err = strconv.ParseFloat(minPriceStr, 64)
		if err != nil {
			api.GetErrorJSON(c, http.StatusBadRequest, "Invalid min_price")
			return
		}
	}

	if maxPriceStr != "" {
		maxPrice, err = strconv.ParseFloat(maxPriceStr, 64)
		if err != nil {
			api.GetErrorJSON(c, http.StatusBadRequest, "Invalid max_price")
			return
		}
	} else {
		maxPrice = 999999999 // Большое число если максимум не указан
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	cards, err := ctrl.cardService.GetCardsByPriceRange(minPrice, maxPrice, page, limit)
	if err != nil {
		api.GetErrorJSON(c, http.StatusInternalServerError, "Failed to get cards")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"cards":     cards,
		"min_price": minPrice,
		"max_price": maxPrice,
		"page":      page,
		"limit":     limit,
	})
}

func (ctrl *cardController) CreateCard(c *gin.Context, request *api.TokenCreateCard) {
	userInfo, err := ExtractUserFromToken(request.User.Login.Token)
	if err != nil {
		api.GetErrorJSON(c, http.StatusUnauthorized, err.Error())
		return
	}

	if !userInfo.IsCompany {
		api.GetErrorJSON(c, http.StatusForbidden, "Only companies can create cards")
		return
	}

	card, err := ctrl.cardService.CreateCard(
		userInfo.UserID,
		request.Card.Title,
		request.Card.Description,
		request.Card.Category,
		request.Card.Location,
		request.Card.Price,
	)
	if err != nil {
		api.GetErrorJSON(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusCreated, gin.H{"card": card})
}

func (ctrl *cardController) UpdateCard(c *gin.Context, request *api.TokenUpdateCard) {
	userInfo, err := ExtractUserFromToken(request.User.Login.Token)
	if err != nil {
		api.GetErrorJSON(c, http.StatusUnauthorized, err.Error())
		return
	}

	if !userInfo.IsCompany {
		api.GetErrorJSON(c, http.StatusForbidden, "Only companies can update cards")
		return
	}

	card, err := ctrl.cardService.UpdateCard(
		request.CardID,
		userInfo.UserID,
		request.Card.Title,
		request.Card.Description,
		request.Card.Category,
		request.Card.Location,
		request.Card.Price,
	)
	if err != nil {
		api.GetErrorJSON(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"card": card})
}

func (ctrl *cardController) DeleteCard(c *gin.Context, request *api.TokenDeleteCard) {
	userInfo, err := ExtractUserFromToken(request.User.Login.Token)
	if err != nil {
		api.GetErrorJSON(c, http.StatusUnauthorized, err.Error())
		return
	}

	if !userInfo.IsCompany {
		api.GetErrorJSON(c, http.StatusForbidden, "Only companies can delete cards")
		return
	}

	err = ctrl.cardService.DeleteCard(request.CardID, userInfo.UserID)
	if err != nil {
		api.GetErrorJSON(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Card deleted successfully"})
}

func (ctrl *cardController) GetCompanyCards(c *gin.Context, request *api.TokenListCard) {
	userInfo, err := ExtractUserFromToken(request.User.Login.Token)
	if err != nil {
		api.GetErrorJSON(c, http.StatusUnauthorized, err.Error())
		return
	}

	if !userInfo.IsCompany {
		api.GetErrorJSON(c, http.StatusForbidden, "Only companies can access their cards")
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

	cards, err := ctrl.cardService.GetCardsByCompany(userInfo.UserID, page, limit)
	if err != nil {
		api.GetErrorJSON(c, http.StatusInternalServerError, "Failed to get cards")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"cards": cards,
		"page":  page,
		"limit": limit,
	})
}

func NewCardController(cardService service.CardService) CardController {
	return &cardController{cardService: cardService}
}
