package api

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/ruhancs/bank-go/token"
)

const (
	authorizationHeaderKey = "authorization"
	authorizationBearerType = "bearer"
	authorizationPayloadKey = "authorization_payload"
)

//retorna a middleware de authenticacao
func authMiddleware(tokenMaker token.Maker) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		//pegar o header
		authorizationHeader := ctx.GetHeader(authorizationHeaderKey)
		if len(authorizationHeader) == 0 {
			err := errors.New("authorization header is not provided")
			//abortar a requisicao e enviar o erro de unauthorized
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		//separar a string do bearer token
		fields := strings.Fields(authorizationHeader)
		if len(fields) < 2 {
			err := errors.New("invalid authorization header format")
			//abortar a requisicao e enviar o erro de unauthorized
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		//pegar o bearer em minusculo
		authorizationType := strings.ToLower(fields[0])
		if authorizationType != authorizationBearerType {
			err := errors.New("invalid authorization type")
			//abortar a requisicao e enviar o erro de unauthorized
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		accessToken := fields[1]
		payload,err := tokenMaker.VerifyToken(accessToken)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		//setar o token no contexto com nome de authorization_payload
		ctx.Set(authorizationPayloadKey, payload)
		ctx.Next()//passa para a proxima funcao
	}
}