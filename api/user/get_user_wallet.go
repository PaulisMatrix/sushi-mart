package user

import (
	"context"
	"net/http"
	"sushi-mart/common"

	"github.com/jackc/pgx/v5"
	"github.com/sirupsen/logrus"
)

func (v *Validator) GetUserWallet(ctx context.Context, Id int) (*GetWalletRes, *common.ErrorResponse) {
	return v.UsersService.GetUserWallet(ctx, Id)
}

func (u *UsersServiceImpl) GetUserWallet(ctx context.Context, Id int) (*GetWalletRes, *common.ErrorResponse) {
	logger := common.ExtractLoggerUnsafe(ctx).WithFields(logrus.Fields{"method": "GetUserWallet", "request": Id})

	resp, err := u.Queries.GetWallet(ctx, int32(Id))

	if err != nil {
		if err == pgx.ErrNoRows {
			logger.WithError(err).Info("wallet not found for this user")
			return nil, &common.ErrorResponse{
				Status:  http.StatusOK,
				Message: "wallet not found. create a wallet first",
			}
		}

		logger.WithError(err).Error("error in getting the user wallet")
		return nil, &common.ErrorResponse{
			Status:  http.StatusInternalServerError,
			Message: "internal server error",
		}
	}
	return &GetWalletRes{
		Username:    resp.Username,
		Balance:     resp.Balance.Abs().InexactFloat64(),
		WalletType:  resp.WalletType,
		WalletAdded: resp.DateAdded.Time.Local().String(),
	}, nil
}
