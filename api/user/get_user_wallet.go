package user

import (
	"context"
	"database/sql"
	"net/http"
	"strconv"
	"sushi-mart/common"
)

func (v *Validator) GetUserWallet(ctx context.Context, Id int) (*GetWalletRes, *common.ErrorResponse) {
	return v.UsersService.GetUserWallet(ctx, Id)
}

func (u *UsersServiceImpl) GetUserWallet(ctx context.Context, Id int) (*GetWalletRes, *common.ErrorResponse) {

	resp, err := u.Queries.GetWallet(ctx, int32(Id))

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, &common.ErrorResponse{
				Status:  http.StatusServiceUnavailable,
				Message: "server unavailable",
			}
		}

		return nil, &common.ErrorResponse{
			Status:  http.StatusInternalServerError,
			Message: "internal server error",
		}
	}
	bal, _ := strconv.ParseFloat(resp.Balance, 64)
	return &GetWalletRes{
		Username:   resp.Username,
		Balance:    bal,
		WalletType: resp.WalletType,
	}, nil
}
