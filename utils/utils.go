package utils

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/gin-gonic/gin"
)

func HandleError(c *gin.Context, statusCode int, message string, err error) {
	if err != nil {
		c.JSON(statusCode, gin.H{"success": false, "status code": statusCode, "error": message, "details": err.Error()})
	} else {
		fmt.Println(message)
		c.JSON(statusCode, gin.H{"success": false, "status code": statusCode, "error": message})
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
