/* The /analytics group contains all the routes to track analytics */
package routes

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"schej.it/server/db"
	"schej.it/server/models"
	"schej.it/server/slackbot"
)

// BasicAuth middleware for analytics routes
func AnalyticsBasicAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		analyticsUsername := os.Getenv("ANALYTICS_USERNAME")
		analyticsPassword := os.Getenv("ANALYTICS_PASSWORD")
		user, pass, hasAuth := c.Request.BasicAuth()

		if !hasAuth || user != analyticsUsername || pass != analyticsPassword {
			c.Header("WWW-Authenticate", `Basic realm="Restricted"`)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}
		c.Next()
	}
}

func InitAnalytics(router *gin.RouterGroup) {
	analyticsRouter := router.Group("/analytics")

	analyticsRouter.POST("/scanned-poster", scannedPoster)
	analyticsRouter.POST("/upgrade-dialog-viewed", upgradeDialogViewed)
	analyticsRouter.GET("/user/:email", AnalyticsBasicAuth(), getUserByEmail)
}

// @Summary Notifies us when poster QR code has been scanned
// @Tags analytics
// @Accept json
// @Produce json
// @Param payload body object{url=string,location=models.Location} true "Object containing the location that poster was scanned from and the url that was scanned"
// @Success 200
// @Router /analytics/scanned-poster [post]
func scannedPoster(c *gin.Context) {
	payload := struct {
		Url      string           `json:"url" binding:"required"`
		Location *models.Location `json:"location"`
	}{}
	if err := c.BindJSON(&payload); err != nil {
		return
	}

	if payload.Location != nil {
		slackbot.SendTextMessage(
			fmt.Sprintf(":face_with_monocle: Poster was scanned :face_with_monocle:\n*Location:* %s, %s, %s\n*URL:* %s",
				payload.Location.City,
				payload.Location.State,
				payload.Location.CountryCode,
				payload.Url,
			),
		)
	} else {
		slackbot.SendTextMessage(
			fmt.Sprintf(":face_with_monocle: Poster was scanned :face_with_monocle:\n*URL:* %s", payload.Url),
		)
	}

	c.JSON(http.StatusOK, gin.H{})
}

// @Summary Notifies us when user has viewed the upgrade dialog
// @Tags analytics
// @Accept json
// @Produce json
// @Param payload body object{userId=string} true "Object containing the user id"
// @Success 200
// @Router /analytics/upgrade-dialog-viewed [post]
func upgradeDialogViewed(c *gin.Context) {
	payload := struct {
		UserId string `json:"userId" binding:"required"`
		Price  string `json:"price" binding:"required"`
		Type   string `json:"type" binding:"required"`
	}{}
	if err := c.BindJSON(&payload); err != nil {
		return
	}

	var message string
	user := db.GetUserById(payload.UserId)
	if user == nil {
		message = fmt.Sprintf(":eyes: %s viewed the upgrade dialog (%s), type: %s", payload.UserId, payload.Price, payload.Type)
	} else {
		message = fmt.Sprintf(":eyes: %s %s (%s) viewed the upgrade dialog (%s), type: %s", user.FirstName, user.LastName, user.Email, payload.Price, payload.Type)
	}

	slackbot.SendTextMessageWithType(
		message,
		slackbot.MONETIZATION,
	)

	c.JSON(http.StatusOK, gin.H{})
}

// @Summary Gets the user by email
// @Tags analytics
// @Accept json
// @Produce json
// @Param email path string true "User email"
// @Success 200 {object} models.User
// @Router /analytics/user/{email} [get]
func getUserByEmail(c *gin.Context) {
	email := c.Param("email")
	user := db.GetUserByEmail(email)
	c.JSON(http.StatusOK, user)
}
