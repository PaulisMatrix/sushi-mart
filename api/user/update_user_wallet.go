package user

import (
	"context"
	"net/http"
	"sushi-mart/common"
	"sushi-mart/internal/database"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
)

func (v *Validator) UpdateUserWallet(ctx context.Context, req *UpdateWalletReq, Id int) *common.ErrorResponse {
	return v.UsersService.UpdateUserWallet(ctx, req, Id)
}

func (u *UsersServiceImpl) UpdateUserWallet(ctx context.Context, req *UpdateWalletReq, Id int) *common.ErrorResponse {
	logger := common.ExtractLoggerUnsafe(ctx).WithFields(logrus.Fields{"method": "UpdateUserWallet", "request": req, "user_id": Id})

	dbParams := database.UpdateBalanceParams{ID: int32(Id), UpdateDateModified: true, DateModified: pgtype.Timestamp{Time: time.Now().Local(), Valid: true}}

	if req.Balance > 0 {
		dbParams.UpdateBalance = true
		dbParams.Balance = decimal.NewFromFloat(req.Balance)
	} else {
		dbParams.UpdateBalance = false
		dbParams.Balance = decimal.NewFromFloat(req.Balance)
	}

	if req.WalletType == "" {
		dbParams.UpdateWalletType = false
		dbParams.WalletType = req.WalletType
	} else {
		dbParams.UpdateWalletType = true
		dbParams.WalletType = req.WalletType
	}

	err := u.Queries.UpdateBalance(ctx, dbParams)

	if err != nil {
		if err == pgx.ErrNoRows {
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
