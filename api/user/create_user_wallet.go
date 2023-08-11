package user

import (
	"context"
	"database/sql"
	"net/http"
	"strconv"
	"sushi-mart/common"
	"sushi-mart/internal/database"
	"time"
)

func (v *Validator) CreateUserWallet(ctx context.Context, req *CreateWalletReq, Id int) *common.ErrorResponse {
	return v.UsersService.CreateUserWallet(ctx, req, Id)
}

func (i *UsersServiceImpl) CreateUserWallet(ctx context.Context, req *CreateWalletReq, Id int) *common.ErrorResponse {
	dbParams := database.CreateWalletParams{
		Balance:      strconv.FormatFloat(req.Balance, 'E', -1, 64),
		WalletType:   req.WalletType,
		DateAdded:    time.Now().Local(),
		DateModified: time.Now().Local(),
		CustomerID:   sql.NullInt32{Int32: int32(Id), Valid: true},
	}

	err := i.Queries.CreateWallet(ctx, dbParams)

	if err != nil {
		return &common.ErrorResponse{
			Status: http.StatusInternalServerError,
		}
	}
	return nil
}
