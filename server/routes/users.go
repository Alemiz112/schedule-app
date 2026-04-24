package routes

import "github.com/gin-gonic/gin"

func InitUsers(router *gin.RouterGroup) {
	usersRouter := router.Group("/users")

	usersRouter.GET("/:userId/is-premium", getIsUserPremium)
}

// @Summary Returns whether the given user is a premium user
// @Tags users
// @Produce json
// @Param userId path string true "User ID"
// @Success 200 {object} object{isPremium=bool}
// @Router /users/{userId}/is-premium [get]
func getIsUserPremium(c *gin.Context) {
	_ = c.Param("userId")
	c.JSON(200, gin.H{"isPremium": true})
}
