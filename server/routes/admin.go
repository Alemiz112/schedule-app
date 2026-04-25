package routes

import (
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"schej.it/server/db"
	"schej.it/server/errs"
	"schej.it/server/middleware"
	"schej.it/server/models"
	"schej.it/server/responses"
)

func InitAdmin(router *gin.RouterGroup) {
	r := router.Group("/admin")
	r.Use(middleware.AuthRequired(), middleware.AdminRequired())

	r.GET("/users", listUsers)
	r.POST("/users", createAdminUser)
	r.PATCH("/users/:userId/role", setUserRole)
	r.DELETE("/users/:userId", deleteAdminUser)
	r.GET("/settings", getAdminSettings)
	r.PATCH("/settings", updateAdminSettings)
}

type adminUserView struct {
	Id        primitive.ObjectID `json:"_id"`
	Email     string             `json:"email"`
	FirstName string             `json:"firstName"`
	LastName  string             `json:"lastName"`
	Role      models.UserRole    `json:"role"`
}

func toAdminUserView(u models.User) adminUserView {
	return adminUserView{
		Id:        u.Id,
		Email:     u.Email,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Role:      u.Role,
	}
}

func listUsers(c *gin.Context) {
	users := db.GetAllUsers()
	views := make([]adminUserView, len(users))
	for i, u := range users {
		views[i] = toAdminUserView(u)
	}
	c.JSON(http.StatusOK, views)
}

func createAdminUser(c *gin.Context) {
	payload := struct {
		Email     string `json:"email" binding:"required"`
		FirstName string `json:"firstName"`
		LastName  string `json:"lastName"`
	}{}
	if err := c.Bind(&payload); err != nil {
		return
	}

	email := strings.ToLower(strings.TrimSpace(payload.Email))
	if existing := db.GetUserByEmail(email); existing != nil {
		c.JSON(http.StatusConflict, responses.Error{Error: errs.UserAlreadyExists})
		return
	}

	user := models.User{
		Id:        primitive.NewObjectID(),
		Email:     email,
		FirstName: strings.TrimSpace(payload.FirstName),
		LastName:  strings.TrimSpace(payload.LastName),
		Role:      models.RoleUser,
	}
	if _, err := db.UsersCollection.InsertOne(context.Background(), user); err != nil {
		c.JSON(http.StatusInternalServerError, responses.Error{Error: "insert-failed"})
		return
	}

	c.JSON(http.StatusOK, toAdminUserView(user))
}

func setUserRole(c *gin.Context) {
	payload := struct {
		Role models.UserRole `json:"role" binding:"required"`
	}{}
	if err := c.Bind(&payload); err != nil {
		return
	}

	if payload.Role != models.RoleUser && payload.Role != models.RoleAdmin {
		c.JSON(http.StatusBadRequest, responses.Error{Error: "invalid-role"})
		return
	}

	userInterface, _ := c.Get("authUser")
	self := userInterface.(*models.User)
	targetId := c.Param("userId")

	if self.Id.Hex() == targetId {
		c.JSON(http.StatusBadRequest, responses.Error{Error: "cannot-change-own-role"})
		return
	}

	target := db.GetUserById(targetId)
	if target == nil {
		c.JSON(http.StatusNotFound, responses.Error{Error: errs.UserDoesNotExist})
		return
	}

	if err := db.UpdateUserRole(targetId, payload.Role); err != nil {
		c.JSON(http.StatusInternalServerError, responses.Error{Error: "update-failed"})
		return
	}

	c.Status(http.StatusOK)
}

func deleteAdminUser(c *gin.Context) {
	userInterface, _ := c.Get("authUser")
	self := userInterface.(*models.User)
	targetId := c.Param("userId")

	if self.Id.Hex() == targetId {
		c.JSON(http.StatusBadRequest, responses.Error{Error: "cannot-delete-self"})
		return
	}

	target := db.GetUserById(targetId)
	if target == nil {
		c.JSON(http.StatusNotFound, responses.Error{Error: errs.UserDoesNotExist})
		return
	}

	if err := db.DeleteUserById(targetId); err != nil {
		c.JSON(http.StatusInternalServerError, responses.Error{Error: "delete-failed"})
		return
	}

	c.Status(http.StatusOK)
}

func getAdminSettings(c *gin.Context) {
	settings := db.GetInstanceSettings()
	c.JSON(http.StatusOK, settings)
}

func updateAdminSettings(c *gin.Context) {
	payload := struct {
		AllowRegistration *bool `json:"allowRegistration" binding:"required"`
	}{}
	if err := c.Bind(&payload); err != nil {
		return
	}

	settings := models.InstanceSettings{
		AllowRegistration: *payload.AllowRegistration,
	}
	if err := db.UpsertInstanceSettings(settings); err != nil {
		c.JSON(http.StatusInternalServerError, responses.Error{Error: "update-failed"})
		return
	}

	c.Status(http.StatusOK)
}
