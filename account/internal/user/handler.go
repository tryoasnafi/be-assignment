package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func SetHandlers(router *gin.RouterGroup, userService Service) {
	handler := UserHandler{userService: userService}
	routerGroup := router.Group("users")
	routerGroup.GET("/:id", handler.GetByID)
}

type UserHandler struct {
	userService Service
}

func (h UserHandler) GetByID(c *gin.Context) {
	userUUID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	u, err := h.userService.GetByID(userUUID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, u)
}
