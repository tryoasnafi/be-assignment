package account

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/tryoasnafi/be-assignment/account/internal/auth"
	"gorm.io/gorm"
)

func SetHandlers(router *gin.RouterGroup, accountService Service) {
	handler := AccountHandler{accountService: accountService}
	routerGroup := router.Group("accounts")
	routerGroup.Use(auth.VerifySession(nil))
	routerGroup.GET("", handler.GetAllAccounts)
	routerGroup.POST("", handler.CreateAccount)
	routerGroup.GET("/:id", handler.GetAccountByID)
}

type AccountHandler struct {
	accountService Service
}

// GetAllAccounts
// @Summary get all account of user
// @Schemes
// @Description get all related accounts of a user
// @Tags accounts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} model.Account
// @Router /accounts [get]
func (h AccountHandler) GetAllAccounts(c *gin.Context) {
	id, err := auth.GetUserIDFromRequest(c.Request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
	}
	accounts, err := h.accountService.GetAllAccounts(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, accounts)
}

// GetAccountByID
// @Summary get account and all payment histories
// @Schemes
// @Description get account by id and all transaction histories
// @Tags accounts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} model.Account
// @Router /accounts/:id [get]
// @Param id path int true "Account ID"
func (h AccountHandler) GetAccountByID(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	account, err := h.accountService.GetAccountByID(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "Account not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, account)
}

// CreateAccount
// @Summary create bank account
// @Schemes
// @Description create bank account
// @Tags accounts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} model.Account
// @Router /accounts [post]
func (h AccountHandler) CreateAccount(c *gin.Context) {
	accountReq := CreateAccountRequest{}
	if err := c.ShouldBindJSON(&accountReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "missing/invalid account request",
		})
		c.Abort()
		return
	}
	id, err := auth.GetUserIDFromRequest(c.Request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
	}
	accounts, err := h.accountService.CreateAccount(id, accountReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusCreated, accounts)
}
