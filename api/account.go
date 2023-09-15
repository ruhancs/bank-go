package api

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	db "github.com/ruhancs/bank-go/db/sqlc"
	"github.com/ruhancs/bank-go/token"
)

//body do http
//VALIDATOR PLAYGROUND VALIDATOR
type createAccountRequest struct {
	//currency criado validator customizado em validator.go e em server
	Currency string `json:"currency" binding:"required,currency"`
}

//gin.Context para pegar os parametros de entrada
func (server *Server) createAccount(ctx *gin.Context) {
	var req createAccountRequest
	//verificar o request
	if err := ctx.ShouldBindJSON(&req); err != nil {
		//enviar json com erro
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	//pegar o usuario pelo payload do token no header, e inserir no token payload
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	arg := db.CreateAccountParams{
		Owner: authPayload.Username,
		Currency: req.Currency,
		Balance: 0,
	}

	account,err := server.store.CreateAccount(ctx,arg)
	if err != nil {
		//erro ao criar account, verificar o tipo do db
		if pqErr,ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "foreign_key_violation", "unique_violation":
				ctx.JSON(http.StatusForbidden, errorResponse(err))
			}
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	ctx.JSON(http.StatusOK, account)
}

type getAccountRequest struct {
	//uri indica que o parametro vem da url
	ID    int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getAccount(ctx *gin.Context) {
	var req getAccountRequest
	//verificar o request
	if err := ctx.ShouldBindUri(&req); err != nil {
		//enviar json com erro
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	account,err := server.store.GetAccount(ctx,req.ID)
	if err != nil {
		//se o id nao existe
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	//pegar o usuario pelo payload do token no header, e inserir no token payload
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	if account.Owner != authPayload.Username {
		err := errors.New("account doesn`t belong to the authenticated user")
		ctx.JSON(http.StatusUnauthorized,errorResponse(err))
	}

	ctx.JSON(http.StatusOK,account)
}

type listAccountRequest struct {
	//form indica que o parametro é um query param
	Page_ID    int32 `form:"page_id" binding:"required,min=1"`
	PageSize    int32 `form:"page_size" binding:"required,min=5,max=10"`
}

func (server *Server) listAccount(ctx *gin.Context) {
	var req listAccountRequest
	//verificar o request
	if err := ctx.ShouldBindQuery(&req); err != nil {
		//enviar json com erro
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	//pegar o usuario pelo payload do token no header, e inserir no token payload
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	arg := db.ListAccountsParams{
		Owner: authPayload.Username,
		Limit: req.PageSize,
		Offset: (req.Page_ID - 1) * req.PageSize,
	}
	accounts,err := server.store.ListAccounts(ctx,arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK,accounts)
}