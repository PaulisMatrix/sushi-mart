package user

import (
	"context"
	"net/http"
	"strconv"
	"sushi-mart/common"
	"sushi-mart/internal/database"
	"time"
)

func (v *Validator) UpdateUserWallet(ctx context.Context, req *UpdateWalletReq, Id int) *common.ErrorResponse {
	return v.UsersService.UpdateUserWallet(ctx, req, Id)
}

func (u *UsersServiceImpl) UpdateUserWallet(ctx context.Context, req *UpdateWalletReq, Id int) *common.ErrorResponse {

	dbParams := database.UpdateBalanceParams{ID: int32(Id), UpdateDateModified: true, DateModified: time.Now().Local()}

	if req.Balance != 0 {
		dbParams.UpdateBalance = true
		dbParams.Balance = strconv.FormatFloat(req.Balance, 'E', -1, 64)
	}

	if req.WalletType != "" {
		dbParams.UpdateWalletType = true
	}

	err := u.Queries.UpdateBalance(ctx, dbParams)

	if err != nil {
		return &common.ErrorResponse{
			Status:  http.StatusInternalServerError,
			Message: "internal server error",
		}
	}

	return nil
}
