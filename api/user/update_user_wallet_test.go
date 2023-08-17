package user_test

import (
	"sushi-mart/api/user"
	"sushi-mart/internal/database"

	"github.com/golang/mock/gomock"
	"github.com/jackc/pgx/v5"
	"github.com/shopspring/decimal"
)

func (u *UsersServiceSuite) TestUpdateWalletOkReq() {
	Id := 4
	req := &user.UpdateWalletReq{
		Balance: float64(21221.21),
	}
	dbParams := database.UpdateBalanceParams{ID: int32(Id)}

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

	u.queriesMock.EXPECT().UpdateBalance(gomock.Any(), dbParams).Return(nil)
	errResp := u.UsersService.UpdateUserWallet(u.context, req, Id)
	u.Nil(errResp)
}

func (u *UsersServiceSuite) TestUpdateWalletBadReq() {
	Id := 15
	req := &user.UpdateWalletReq{
		Balance: float64(21221.21),
	}
	dbParams := database.UpdateBalanceParams{ID: int32(Id)}

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

	expectedErr := pgx.ErrNoRows

	u.queriesMock.EXPECT().UpdateBalance(gomock.Any(), dbParams).Return(expectedErr)
	errResp := u.UsersService.UpdateUserWallet(u.context, req, Id)
	u.NotNil(errResp)
	u.EqualValues(expectedErr.Error(), "no rows in result set")
}
