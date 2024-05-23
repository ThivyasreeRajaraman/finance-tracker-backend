package authservices

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	userhelper "github.com/Thivyasree-Rajaraman/finance-tracker/helpers/query/user"
	"github.com/Thivyasree-Rajaraman/finance-tracker/models"
	"github.com/Thivyasree-Rajaraman/finance-tracker/utils"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"gorm.io/gorm"
)

func getUserInfo(accessToken string) (map[string]interface{}, error) {
	userInfoEndpoint := os.Getenv("GOOGLE_USER_INFO_ENDPOINT")
	//nolint:noctx
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s?access_token=%s", userInfoEndpoint, accessToken), nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var userInfo map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, err
	}
	return userInfo, nil
}

func SignJWT(user *models.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id": user.ID,
		//nolint:add-constant
		utils.EmailKey: user.Email,
		"iss":          "oauth-app-golang",
		"exp":          time.Now().Add(time.Hour * 24 * utils.TokenExpirationDays).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(os.Getenv("SECRET_KEY")))
	if err != nil {
		return utils.EmptyString, err
	}

	return signedToken, nil
}

func ExchangeCodeForTokenAndGetUserInfo(c *gin.Context, googleOauthConfig *oauth2.Config) map[string]interface{} {
	jwtToken := c.Query("code")
	token, err := googleOauthConfig.Exchange(context.Background(), jwtToken)

	if err != nil {
		utils.HandleError(c, http.StatusBadRequest, "Error exchanging jwt token", err)
		return nil
	}

	userInfo, err := getUserInfo(token.AccessToken)
	if err != nil {
		utils.HandleError(c, http.StatusBadRequest, "Error getting user info", err)
		return nil
	}
	return userInfo
}

func extractUserInfo(userInfo map[string]interface{}) (string, string, error) {
	name, nameOk := userInfo["name"].(string)
	email, emailOk := userInfo["email"].(string)
	if !nameOk || !emailOk {
		return utils.EmptyString, utils.EmptyString, fmt.Errorf("failed to extract user info from map")
	}
	return name, email, nil
}

func CheckAndCreateUser(user *models.User, userInfo map[string]interface{}) error {
	// extract user info
	name, email, err := extractUserInfo(userInfo)
	if err != nil {
		return err
	}

	user.Name = name
	user.Email = email

	// check if user already exists
	if err := userhelper.SearchByEmail(user); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// create a new user
			err = userhelper.CreateUser(user)
			if err != nil {
				return fmt.Errorf("failed to create user")
			}
		} else {
			return fmt.Errorf("error retrieving user")
		}
	}
	return nil
}
