package api

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	db "github.com/sayedppqq/banking-backend/db/sqlc"
	"github.com/sayedppqq/banking-backend/util"
	"net/http"
)

type createTransferRequest struct {
	FromAccountID int64  `json:"from_account_id" binding:"required,min=1"`
	ToAccountID   int64  `json:"to_account_id" binding:"required,min=1"`
	Amount        int64  `json:"amount" binding:"required,gt=0"`
	Currency      string `json:"currency" binding:"required,validCurrency"`
}

func (server *Server) createTransfer(ctx *gin.Context) {
	var req createTransferRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if err := server.validTransfer(ctx, req.FromAccountID, req.Currency); err != nil {
		return
	}
	if err := server.validTransfer(ctx, req.ToAccountID, req.Currency); err != nil {
		return
	}

	arg := db.TransferTxParams{
		FromAccountID: req.FromAccountID,
		ToAccountID:   req.ToAccountID,
		Amount:        req.Amount,
	}

	result, err := server.store.TransferTx(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, result)
}
func (server *Server) validTransfer(ctx *gin.Context, accountID int64, currency string) error {
	account, err := server.store.GetAccount(ctx, accountID)
	if err != nil {
		if errors.Is(err, util.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return err
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return err
	}
	if currency != account.Currency {
		err := fmt.Errorf("account [%d] currency mismatch: %s vs %s", account.ID, account.Currency, currency)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return err
	}
	return nil
}
