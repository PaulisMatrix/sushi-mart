package user_test

import (
	"errors"
	"sushi-mart/api/user"
	"sushi-mart/internal/database"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/shopspring/decimal"
)

func (u *UsersServiceSuite) TestCreateWalletOkReq() {
	Id := 10

	req := &user.CreateWalletReq{
		Balance:    121212.1212,
		WalletType: "PAYTM",
	}

	dbParams := database.CreateWalletParams{
		Balance:      decimal.NewFromFloat(req.Balance),
		WalletType:   req.WalletType,
		DateAdded:    pgtype.Timestamp{Time: time.Now().Local(), Valid: true},
		DateModified: pgtype.Timestamp{Time: time.Now().Local(), Valid: true},
		CustomerID:   pgtype.Int4{Int32: int32(Id), Valid: true},
	}

	u.queriesMock.EXPECT().CreateWallet(gomock.Any(), dbParams).Return(nil)
	errResp := u.UsersService.CreateUserWallet(u.context, req, Id)
	u.Nil(errResp)
}

func (u *UsersServiceSuite) TestCreateWalletBadReq() {
	Id := 10

	req := &user.CreateWalletReq{
		Balance:    float64(121212.1212),
		WalletType: "PAYTM",
	}

	dbParams := database.CreateWalletParams{
		Balance:      decimal.NewFromFloat(req.Balance),
		WalletType:   req.WalletType,
		DateAdded:    pgtype.Timestamp{Time: time.Now().Local(), Valid: true},
		DateModified: pgtype.Timestamp{Time: time.Now().Local(), Valid: true},
		CustomerID:   pgtype.Int4{Int32: int32(Id), Valid: true},
	}

	expectedErr := errors.New("internal server error")
	u.queriesMock.EXPECT().CreateWallet(gomock.Any(), dbParams).Return(expectedErr)
	errResp := u.UsersService.CreateUserWallet(u.context, req, Id)
	u.NotNil(errResp)
	u.EqualValues(expectedErr.Error(), errResp.Message)
}
