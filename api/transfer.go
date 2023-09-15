package api

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/ruhancs/bank-go/db/sqlc"
	"github.com/ruhancs/bank-go/token"
)

//body do http
//VALIDATOR PLAYGROUND VALIDATOR
type createTranferRequest struct {
	FromAccountID    int64 `json:"from_account_id" binding:"required,min=1"`
	ToAccountID int64 `json:"to_account_id" binding:"required,min=1"`
	Amount int64 `json:"amount" binding:"required,gt=1"`
	Currency string `json:"currency" binding:"required,currency"`
}

//gin.Context para pegar os parametros de entrada
func (server *Server) createTransfer(ctx *gin.Context) {
	var req createTranferRequest
	//verificar o request
	if err := ctx.ShouldBindJSON(&req); err != nil {
		//enviar json com erro
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	//verificar se a conta que ira tranferir é a msm logada
	fromAccount,valid := server.validAccount(ctx,req.FromAccountID,req.Currency)
	if !valid {
		return
	}
	
	//pegar o usuario pelo payload do token no header, e inserir no token payload
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	if fromAccount.Owner != authPayload.Username {
		err := errors.New("from account doesn`t belong to the authenticated user")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}
	
	_,valid = server.validAccount(ctx,req.ToAccountID,req.Currency)
	if !valid {
		return
	}


	arg := db.TranferTxParams{
		FromAccountID: req.FromAccountID,
		ToAccountID: req.ToAccountID,
		Amount: req.Amount,
	}

	result,err := server.store.TranferTx(ctx,arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	ctx.JSON(http.StatusOK, result)
}

func(server *Server) validAccount(ctx *gin.Context, accountID int64, currency string) (db.Account, bool) {
	account,err := server.store.GetAccount(ctx,accountID)
	if err != nil {
		//se o id nao existe
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return account,false
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return account,false
	}

	if account.Currency != currency {
		err := fmt.Errorf("account [%d] currency mismatch %s vs %s",account.ID,account.Currency, currency)
		//enviar resposta com o erro
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return account,false
	}

	return account,true
}