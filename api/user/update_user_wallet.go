package user

import (
	"context"
	"database/sql"
	"net/http"
	"strconv"
	"sushi-mart/common"
	"sushi-mart/internal/database"
	"time"

	"github.com/sirupsen/logrus"
)

func (v *Validator) UpdateUserWallet(ctx context.Context, req *UpdateWalletReq, Id int) *common.ErrorResponse {
	return v.UsersService.UpdateUserWallet(ctx, req, Id)
}

func (u *UsersServiceImpl) UpdateUserWallet(ctx context.Context, req *UpdateWalletReq, Id int) *common.ErrorResponse {
	logger := common.ExtractLoggerUnsafe(ctx).WithFields(logrus.Fields{"method": "UpdateUserWallet", "request": req, "user_id": Id})

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
		if err == sql.ErrNoRows {
			logger.WithError(err).Info("invalid user id to update the wallet")
			return &common.ErrorResponse{
				Status:  http.StatusOK,
				Message: "invalid user id. pass correct user id to update the record",
			}
		}
		return &common.ErrorResponse{
			Status:  http.StatusInternalServerError,
			Message: "internal server error",
		}
	}

	return nil
}
