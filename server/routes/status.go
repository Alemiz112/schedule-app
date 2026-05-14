package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"schej.it/server/db"
)

func InitStatus(router *gin.RouterGroup) {
	router.GET("/status", getInstanceStatus)
}

// getInstanceStatus returns public instance settings (e.g. whether registration is open)
func getInstanceStatus(c *gin.Context) {
	settings := db.GetInstanceSettings()
	c.JSON(http.StatusOK, gin.H{
		"allowRegistration": settings.AllowRegistration,
	})
}
