package api

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kanatsanan6/go-test/utils"
)

const (
	authorizationPayloadKey = "authorization_payload"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.GetHeader("authorization")

		if len(token) == 0 {
			err := errors.New("authorization header is empty")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"errors": err.Error()})
			return
		}

		result, err := utils.ValidateJwtToken(token)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"errors": err.Error()})
			return
		}

		ctx.Set(authorizationPayloadKey, result)
		ctx.Next()
	}

}
