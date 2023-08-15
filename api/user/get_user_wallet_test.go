package user_test

import (
	"database/sql"
	"sushi-mart/internal/database"
	"time"

	"github.com/golang/mock/gomock"
)

func (u *UsersServiceSuite) TestGetWalletOkReq() {
	Id := 1

	expectedRow := database.GetWalletRow{
		Username:   "testing",
		Balance:    "212313.12",
		WalletType: "PAYTM",
		DateAdded:  time.Now().Local(),
	}

	u.queriesMock.EXPECT().GetWallet(gomock.Any(), int32(Id)).Return(expectedRow, nil)
	resp, errResp := u.UsersService.GetUserWallet(u.context, Id)
	u.Nil(errResp)
	u.NotNil(resp)
	u.EqualValues(resp.Username, expectedRow.Username)
}

func (u *UsersServiceSuite) TestGetWalletNoCustBadReq() {
	Id := 1
	expectedRow := database.GetWalletRow{}
	expectedErr := sql.ErrNoRows
	expectedErrMsg := "wallet not found. create a wallet first"

	u.queriesMock.EXPECT().GetWallet(gomock.Any(), int32(Id)).Return(expectedRow, expectedErr)
	resp, errResp := u.UsersService.GetUserWallet(u.context, Id)
	u.NotNil(errResp)
	u.Nil(resp)
	u.EqualValues(expectedErrMsg, errResp.Message)
}
