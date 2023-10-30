package api

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/timly278/simplebank/token"
)

const (
	AUTHORIZATION_HEADER_KEY  = "authorization"
	AUTHORIZATION_TYPE_BEARER = "bearer"
	AUTHORIZATION_PAYLOAD_KEY = "authorization_payload"
)

func authMiddleware(tokenMaker token.Maker) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authorizationHeader := ctx.GetHeader(AUTHORIZATION_HEADER_KEY)
		if len(authorizationHeader) == 0 {
			err := errors.New("authorization header is not provided")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		fields := strings.Fields(authorizationHeader)
		if len(fields) < 2 {
			err := errors.New("invalid authorization header format")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		authorizationType := strings.ToLower(fields[0])
		if authorizationType != AUTHORIZATION_TYPE_BEARER {
			err := fmt.Errorf("unsupported authorization type %s", authorizationType)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}
		accessToken := fields[1]
		payload, err := tokenMaker.VerifyToken(accessToken)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}
		ctx.Set(AUTHORIZATION_PAYLOAD_KEY, payload)
		ctx.Next() // forward the request to the next handler
	}
}
// authMiddleware just provides access to real handler
// it doesn't care the user who owes the token has permisssion to 
// perform the request or not. 