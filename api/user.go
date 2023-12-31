package api

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/lib/pq"
	db "github.com/ruhancs/bank-go/db/sqlc"
	"github.com/ruhancs/bank-go/util"
)

//body do http
//VALIDATOR PLAYGROUND VALIDATOR
type createUserRequest struct {
	Username    string `json:"username" binding:"required,alphanum"`
	Password    string `json:"password" binding:"required,min=6"`
	FullName string `json:"full_name" binding:"required"`
	Email string `json:"email" binding:"required"`
}

type userResponse struct {
	Username          string    `json:"username"`
	FullName          string    `json:"full_name"`
	Email             string    `json:"email"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
	CreatedAt         time.Time `json:"created_at"`
}

func newUserResponse(user db.User) userResponse {
	return userResponse{
		Username: user.Username,
		FullName: user.FullName,
		Email: user.Email,
		PasswordChangedAt: user.PasswordChangedAt,
		CreatedAt: user.CreatedAt,
	}
}

//gin.Context para pegar os parametros de entrada
func (server *Server) createUser(ctx *gin.Context) {
	var req createUserRequest
	//verificar o request
	if err := ctx.ShouldBindJSON(&req); err != nil {
		//enviar json com erro
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	hashedPassword,err := util.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	arg := db.CreateUserParams{
		Username: req.Username,
		HashedPassword: hashedPassword,
		FullName: req.FullName,
		Email: req.Email,
	}

	user,err := server.store.CreateUser(ctx,arg)
	if err != nil {
		//erro ao criar user, verificar o tipo do db
		if pqErr,ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation": //verificar erro de username ja existe
				ctx.JSON(http.StatusForbidden, errorResponse(err))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	resp := newUserResponse(user)

	ctx.JSON(http.StatusOK, resp)
}

type loginUserRequest struct {
	Username    string `json:"username" binding:"required,alphanum"`
	Password    string `json:"password" binding:"required,min=6"`
}

type loginUserResponse struct {
	SessionID uuid.UUID `json:"session_id"`
	AccessToken string `json:"access_token"`
	AccessTokenExpiresAt time.Time `json:"access_token_expires_at"`
	RefreshToken string `json:"refresh_token"`
	RefreshTokenExpiresAt time.Time `json:"refresh_token_expires_at"`
	User userResponse `json:"user" binding:"required"`
}

func (server *Server) loginUser(ctx *gin.Context) {
	err := godotenv.Load()
	if err != nil {
		fmt.Println(err)
		fmt.Println("Error loading .env")
	}
	
	var req loginUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
	}

	user,err := server.store.GetUser(ctx,req.Username)
	if err !=  nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	
	err = util.CheckPassword(req.Password, user.HashedPassword)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	tokenDuration,_ := strconv.Atoi(os.Getenv("ACCESS_TOKEN_DURATION")) 
	duration := time.Duration(tokenDuration * int(time.Minute))
	accessToken,accessTokenPayload,err := server.tokenMaker.CreateToken(req.Username,duration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError,errorResponse(err))
		return
	}

	//gerar refresh token
	refreshTokenDuration,_ := strconv.Atoi(os.Getenv("REFRESH_TOKEN_DURATION")) 
	refreshDuration := time.Duration(refreshTokenDuration * int(time.Hour))
	refreshToken,refreshTokenPayload,err:= server.tokenMaker.CreateToken(user.Username,refreshDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError,errorResponse(err))
		return
	}

	session,err := server.store.CreateSession(ctx,db.CreateSessionParams{
		ID: refreshTokenPayload.ID,
		Username: user.Username,
		RefreshToken: refreshToken,
		UserAgent: ctx.Request.UserAgent(),
		ClientIp: ctx.ClientIP(),
		IsBlocked: false,
		ExpiresAt: refreshTokenPayload.ExpiredAt,
	})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError,errorResponse(err))
		return
	}

	resp := loginUserResponse{
		SessionID: session.ID,
		AccessToken: accessToken,
		AccessTokenExpiresAt: accessTokenPayload.ExpiredAt,
		RefreshToken: refreshToken,
		RefreshTokenExpiresAt: refreshTokenPayload.ExpiredAt,
		User: newUserResponse(user),
	}
	ctx.JSON(http.StatusOK, resp)
}
