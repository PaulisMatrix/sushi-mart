package user

import (
	"context"
	"net/http"
	"sushi-mart/common"
	"sushi-mart/internal/database"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
)

func (v *Validator) CreateUserWallet(ctx context.Context, req *CreateWalletReq, Id int) *common.ErrorResponse {
	return v.UsersService.CreateUserWallet(ctx, req, Id)
}

func (i *UsersServiceImpl) CreateUserWallet(ctx context.Context, req *CreateWalletReq, Id int) *common.ErrorResponse {
	logger := common.ExtractLoggerUnsafe(ctx).WithFields(logrus.Fields{"method": "CreateUserWallet", "request": req})

	dbParams := database.CreateWalletParams{
		Balance:      decimal.NewFromFloat(req.Balance),
		WalletType:   req.WalletType,
		DateAdded:    pgtype.Timestamp{Time: time.Now().Local(), Valid: true},
		DateModified: pgtype.Timestamp{Time: time.Now().Local(), Valid: true},
		CustomerID:   pgtype.Int4{Int32: int32(Id), Valid: true},
	}

	err := i.Queries.CreateWallet(ctx, dbParams)

	if err != nil {
		logger.WithError(err).Error("error in creating a new user wallet")
		return &common.ErrorResponse{
			Status:  http.StatusInternalServerError,
			Message: "internal server error",
		}
	}
	return nil
}
