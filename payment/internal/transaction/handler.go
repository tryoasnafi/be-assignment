package transaction

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tryoasnafi/be-assignment/payment/internal/auth"
	. "github.com/tryoasnafi/be-assignment/payment/internal/transaction/dto"
	"gorm.io/gorm"
)

func SetHandlers(router *gin.RouterGroup, transactionService Service) {
	handler := AccountHandler{transactionService: transactionService}
	routerGroup := router.Group("transactions")
	routerGroup.Use(auth.Verify)
	routerGroup.GET("/test", handler.test)
	routerGroup.POST("/send", handler.Send)
	routerGroup.POST("/withdraw", handler.Withdraw)
}

type AccountHandler struct {
	transactionService Service
}

func (h AccountHandler) test(c *gin.Context) {
	c.Status(http.StatusOK)
}

func (h AccountHandler) Send(c *gin.Context) {
	sendReq := SendRequest{}
	if err := c.ShouldBindJSON(&sendReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "missing/invalid send request",
		})
		c.Abort()
		return
	}

	resp, err := h.transactionService.Send(sendReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h AccountHandler) Withdraw(c *gin.Context) {
	withdrawReq := WithdrawRequest{}
	if err := c.ShouldBindJSON(&withdrawReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "missing/invalid withdraw request",
		})
		c.Abort()
		return
	}

	resp, err := h.transactionService.Withdraw(withdrawReq)
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
	c.JSON(http.StatusOK, resp)
}
