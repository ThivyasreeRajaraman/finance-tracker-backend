package controllers

import (
	"net/http"
	"os"

	"github.com/Thivyasree-Rajaraman/finance-tracker/initializers"
	"github.com/Thivyasree-Rajaraman/finance-tracker/models"
	authservices "github.com/Thivyasree-Rajaraman/finance-tracker/services/auth"
	"github.com/Thivyasree-Rajaraman/finance-tracker/utils"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var googleOauthConfig = &oauth2.Config{}

//nolint:gochecknoinits
func init() {
	initializers.LoadEnvVariables(false)
	initializers.InitDB(os.Getenv("DATABASE_URL"))
	initializers.SyncDatabase()
	googleOauthConfig = &oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("GOOGLE_REDIRECT_URL"),
		Scopes:       []string{"profile", "email"},
		Endpoint:     google.Endpoint,
	}
}

func GoogleCallbackHandler(c *gin.Context) {
	userInfo := authservices.ExchangeCodeForTokenAndGetUserInfo(c, googleOauthConfig)
	var user models.User

	// extract user info and check if user already exists
	err := authservices.CheckAndCreateUser(&user, userInfo)
	if err != nil {
		utils.HandleError(c, http.StatusBadRequest, "Error creating user", err)
		return
	}

	// sign user info using secret key
	signedToken, err := authservices.SignJWT(&user)
	if err != nil {
		utils.HandleError(c, http.StatusBadRequest, "Error signing token", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "user": user, "token": signedToken})
}
