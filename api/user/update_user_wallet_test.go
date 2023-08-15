package user_test

import (
	"database/sql"
	"strconv"
	"sushi-mart/api/user"
	"sushi-mart/internal/database"

	"github.com/golang/mock/gomock"
)

func (u *UsersServiceSuite) TestUpdateWalletOkReq() {
	Id := 4
	req := &user.UpdateWalletReq{
		Balance: float64(21221.21),
	}
	dbParams := database.UpdateBalanceParams{ID: int32(Id)}

	if req.Balance > 0 {
		dbParams.UpdateBalance = true
		dbParams.Balance = strconv.FormatFloat(req.Balance, 'E', -1, 64)
	} else {
		dbParams.UpdateBalance = false
		dbParams.Balance = strconv.FormatFloat(req.Balance, 'E', -1, 64)
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
		dbParams.Balance = strconv.FormatFloat(req.Balance, 'E', -1, 64)
	} else {
		dbParams.UpdateBalance = false
		dbParams.Balance = strconv.FormatFloat(req.Balance, 'E', -1, 64)
	}

	if req.WalletType == "" {
		dbParams.UpdateWalletType = false
		dbParams.WalletType = req.WalletType
	} else {
		dbParams.UpdateWalletType = true
		dbParams.WalletType = req.WalletType
	}

	expectedErr := sql.ErrNoRows
	expectedErrMsg := "invalid user id. pass correct user id to update the record"

	u.queriesMock.EXPECT().UpdateBalance(gomock.Any(), dbParams).Return(expectedErr)
	errResp := u.UsersService.UpdateUserWallet(u.context, req, Id)
	u.NotNil(errResp)
	u.EqualValues(expectedErrMsg, errResp.Message)
}
