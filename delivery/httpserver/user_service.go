package httpserver

import (
	"Game-Application/repository/mongo"
	"Game-Application/service/authservice"
	"Game-Application/service/user"
	"github.com/gin-gonic/gin"
	"net/http"
)

func UserRegister(c *gin.Context) {

	var request user.RegisterRequest
	err := c.Bind(&request)
	if err != nil {
		return
	}
	repo, Merr := mongo.New("mongodb://localhost:27017", "game")
	if Merr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": Merr.Error()})
	}
	authService := authservice.New()
	UserSvc := user.New(authService, repo)
	_, RErr := UserSvc.Register(request)
	if RErr != nil {

		c.JSON(http.StatusInternalServerError, gin.H{"error": RErr.Error()})

	}
	c.JSON(http.StatusOK, gin.H{"success": "True"})
	return
}
