package user

import (
	"context"
	"database/sql"
	"net/http"
	"strconv"
	"sushi-mart/common"

	"github.com/sirupsen/logrus"
)

func (v *Validator) GetUserWallet(ctx context.Context, Id int) (*GetWalletRes, *common.ErrorResponse) {
	return v.UsersService.GetUserWallet(ctx, Id)
}

func (u *UsersServiceImpl) GetUserWallet(ctx context.Context, Id int) (*GetWalletRes, *common.ErrorResponse) {
	logger := common.ExtractLoggerUnsafe(ctx).WithFields(logrus.Fields{"method": "GetUserWallet", "request": Id})

	resp, err := u.Queries.GetWallet(ctx, int32(Id))

	if err != nil {
		if err == sql.ErrNoRows {
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
	bal, _ := strconv.ParseFloat(resp.Balance, 64)
	return &GetWalletRes{
		Username:    resp.Username,
		Balance:     bal,
		WalletType:  resp.WalletType,
		WalletAdded: resp.DateAdded.String(),
	}, nil
}
