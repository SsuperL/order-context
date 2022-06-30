package common

import (
	"order-context/ohs/local/pl/errors"

	"github.com/gin-gonic/gin"
)

// ValidateBody ...
func ValidateBody(c *gin.Context, body interface{}) (interface{}, error) {
	err := c.ShouldBindJSON(&body)
	if err != nil {
		return nil, errors.BadRequest("invalid body")
	}

	return body, nil
}
