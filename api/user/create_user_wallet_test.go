package user_test

import (
	"database/sql"
	"errors"
	"strconv"
	"sushi-mart/api/user"
	"sushi-mart/internal/database"

	"github.com/golang/mock/gomock"
)

func (u *UsersServiceSuite) TestCreateWalletOkReq() {
	Id := 10

	req := &user.CreateWalletReq{
		Balance:    float64(121212.1212),
		WalletType: "PAYTM",
	}

	dbParams := database.CreateWalletParams{
		Balance:    strconv.FormatFloat(req.Balance, 'E', -1, 64),
		WalletType: req.WalletType,
		//DateAdded:    time.Now().Local(),
		//DateModified: time.Now().Local(),
		CustomerID: sql.NullInt32{Int32: int32(Id), Valid: true},
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
		Balance:    strconv.FormatFloat(req.Balance, 'E', -1, 64),
		WalletType: req.WalletType,
		//DateAdded:    time.Now().Local(),
		//DateModified: time.Now().Local(),
		CustomerID: sql.NullInt32{Int32: int32(Id), Valid: true},
	}

	expectedErr := errors.New("internal server error")
	u.queriesMock.EXPECT().CreateWallet(gomock.Any(), dbParams).Return(expectedErr)
	errResp := u.UsersService.CreateUserWallet(u.context, req, Id)
	u.NotNil(errResp)
	u.EqualValues(expectedErr.Error(), errResp.Message)
}
