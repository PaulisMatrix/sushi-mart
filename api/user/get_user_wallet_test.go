package user_test

import (
	"sushi-mart/internal/database"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/shopspring/decimal"
)

func (u *UsersServiceSuite) TestGetWalletOkReq() {
	Id := 1

	expectedRow := database.GetWalletRow{
		Username:   "testing",
		Balance:    decimal.NewFromFloat(212313.12),
		WalletType: "PAYTM",
		DateAdded:  pgtype.Timestamp{Time: time.Now().Local(), Valid: true},
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
	expectedErr := pgx.ErrNoRows

	u.queriesMock.EXPECT().GetWallet(gomock.Any(), int32(Id)).Return(expectedRow, expectedErr)
	resp, errResp := u.UsersService.GetUserWallet(u.context, Id)
	u.NotNil(errResp)
	u.Nil(resp)
	u.EqualValues(expectedErr.Error(), "no rows in result set")
}
