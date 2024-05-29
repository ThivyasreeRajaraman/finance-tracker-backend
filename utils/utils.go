package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/Thivyasree-Rajaraman/finance-tracker/models"
	"github.com/gin-gonic/gin"
)

func HandleError(c *gin.Context, statusCode int, message string, err error) {
	if err != nil {
		c.JSON(statusCode, gin.H{
			"success":     false,
			"status code": statusCode,
			"error":       message,
			"details":     err.Error(),
		})
	} else {
		fmt.Println(message)
		c.JSON(statusCode, gin.H{
			"success":     false,
			"status code": statusCode,
			"error":       message,
		})
	}
}

func UnmarshalData(c *gin.Context, model interface{}) error {
	data, err := io.ReadAll(c.Request.Body)
	if err != nil {
		return err
	}
	if len(data) == Zero {
		return ErrEmptyRequestBody
	}
	return json.Unmarshal(data, &model)
}

func SendResponse(c *gin.Context, message string, key string, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": message,
		key:       data,
	})
}

func CreateError(message string) error {
	return errors.New(message)
}

func GetUserID(c *gin.Context) (uint, error) {
	userInterface, _ := c.Get("currentUser")
	user, ok := userInterface.(models.User)
	if !ok {
		err := CreateError("invalid user data")
		HandleError(c, http.StatusBadRequest, err.Error(), nil)
		return Zero, err
	}
	return user.ID, nil
}

func ParseUintParam(c *gin.Context, paramName string) (uint, error) {
	paramValue := c.Param(paramName)
	paramUint, err := strconv.ParseUint(paramValue, 10, 64)
	if err != nil {
		HandleError(c, http.StatusBadRequest, fmt.Errorf("invalid %s", paramName).Error(), err)
		return 0, err
	}
	return uint(paramUint), nil
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
