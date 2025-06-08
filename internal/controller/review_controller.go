package controller

import (
	"core/internal/api"
	"core/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type ReviewController interface {
	CreateReview(c *gin.Context, request *api.TokenCreateReview)
	GetCompanyReviews(c *gin.Context)
	GetOrderReview(c *gin.Context)
	GetCompanyRating(c *gin.Context)
}

type reviewController struct {
	reviewService service.ReviewService
}

func (ctrl *reviewController) CreateReview(c *gin.Context, request *api.TokenCreateReview) {
	userInfo, err := ExtractUserFromToken(request.User.Login.Token)
	if err != nil {
		api.GetErrorJSON(c, http.StatusUnauthorized, err.Error())
		return
	}

	if userInfo.IsCompany {
		api.GetErrorJSON(c, http.StatusForbidden, "Only clients can create reviews")
		return
	}

	review, err := ctrl.reviewService.CreateReview(
		userInfo.UserID,
		request.Review.CompanyID,
		request.Review.OrderID,
		request.Review.Rating,
		request.Review.Comment,
	)
	if err != nil {
		api.GetErrorJSON(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusCreated, gin.H{"review": review})
}

func (ctrl *reviewController) GetCompanyReviews(c *gin.Context) {
	companyIDStr := c.Param("company_id")
	companyID, err := strconv.ParseUint(companyIDStr, 10, 32)
	if err != nil {
		api.GetErrorJSON(c, http.StatusBadRequest, "Invalid company ID")
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

	reviews, err := ctrl.reviewService.GetReviewsByCompany(uint(companyID), page, limit)
	if err != nil {
		api.GetErrorJSON(c, http.StatusInternalServerError, "Failed to get reviews")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"reviews":    reviews,
		"company_id": companyID,
		"page":       page,
		"limit":      limit,
	})
}

func (ctrl *reviewController) GetOrderReview(c *gin.Context) {
	orderIDStr := c.Param("order_id")
	orderID, err := strconv.ParseUint(orderIDStr, 10, 32)
	if err != nil {
		api.GetErrorJSON(c, http.StatusBadRequest, "Invalid order ID")
		return
	}

	review, err := ctrl.reviewService.GetReviewByOrder(uint(orderID))
	if err != nil {
		api.GetErrorJSON(c, http.StatusNotFound, "Review not found")
		return
	}

	c.JSON(http.StatusOK, gin.H{"review": review})
}

func (ctrl *reviewController) GetCompanyRating(c *gin.Context) {
	companyIDStr := c.Param("company_id")
	companyID, err := strconv.ParseUint(companyIDStr, 10, 32)
	if err != nil {
		api.GetErrorJSON(c, http.StatusBadRequest, "Invalid company ID")
		return
	}

	rating, count, err := ctrl.reviewService.GetCompanyRating(uint(companyID))
	if err != nil {
		api.GetErrorJSON(c, http.StatusInternalServerError, "Failed to get rating")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"company_id":   companyID,
		"rating":       rating,
		"review_count": count,
	})
}

func NewReviewController(reviewService service.ReviewService) ReviewController {
	return &reviewController{reviewService: reviewService}
}
