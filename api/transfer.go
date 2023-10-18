package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/timly278/simplebank/db/sqlc"
)

type createListTransferRequest struct {
	PageID    int32 `form:"page_id" binding:"required,min=1"`
	PageSize  int32 `form:"page_size" binding:"required,min=5,max=20"`
	AccountID int64 `form:"account_id" binding:"required,min=1"`
}

func (server *Server) listTransfer(ctx *gin.Context) {
	var req createListTransferRequest

	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.ListTransfersParams{
		FromAccountID: req.AccountID,
		ToAccountID:   req.AccountID,
		Limit:         req.PageSize,
		Offset:        (req.PageID - 1) * req.PageSize,
	}

	transfers, err := server.store.ListTransfers(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, transfers)
}
