package api

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/ruhancs/bank-go/db/sqlc"
)

//body do http
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

	if !server.validAccount(ctx, req.FromAccountID,req.Currency) {
		return
	}

	if !server.validAccount(ctx, req.ToAccountID,req.Currency) {
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

func(server *Server) validAccount(ctx *gin.Context, accountID int64, currency string) bool {
	account,err := server.store.GetAccount(ctx,accountID)
	if err != nil {
		//se o id nao existe
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return false
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return false
	}

	if account.Currency != currency {
		err := fmt.Errorf("account [%d] currency mismatch %s vs %s",account.ID,account.Currency, currency)
		//enviar resposta com o erro
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return false
	}

	return true
}