package user_test

import (
	"database/sql"
	"sushi-mart/api/user"
	"sushi-mart/internal/database"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/shopspring/decimal"
)

func (u *UsersServiceSuite) TestAddReviewOkReq() {
	custId := 1
	req := &user.AddReviewReq{
		Rating:     4,
		ReviewText: "testing review",
		ProductId:  1,
	}

	expectedResp := database.Order{
		ID:          1,
		OrderStatus: "DELIVERED",
		TotalAmt:    decimal.NewFromFloat(1212.212),
		Units:       int32(10),
		PaymentType: "GPAY",
		OrderDate:   pgtype.Timestamp{Time: time.Now().Local(), Valid: true},
		CustomerID:  pgtype.Int4{Int32: int32(custId), Valid: true},
		ProductID:   pgtype.Int4{Int32: int32(1), Valid: true},
		IsActive:    true,
	}

	dbParams := database.AddReviewParams{
		Rating:     int32(req.Rating),
		ReviewText: req.ReviewText,
		CustomerID: pgtype.Int4{Int32: int32(custId), Valid: true},
		ProductID:  pgtype.Int4{Int32: int32(1), Valid: true},
	}

	u.queriesMock.EXPECT().ValidateProductOrderReview(gomock.Any(), sql.NullInt32{Int32: int32(req.ProductId), Valid: true}).Return(expectedResp, nil)
	u.queriesMock.EXPECT().AddReview(gomock.Any(), dbParams).Return(nil)
	errResp := u.UsersService.AddReview(u.context, req, custId)
	u.Nil(errResp)
}

func (u *UsersServiceSuite) TestAddReviewOrderNABadReq() {
	custId := 1
	req := &user.AddReviewReq{
		Rating:     4,
		ReviewText: "testing review",
		ProductId:  1,
	}

	expectedErr := pgx.ErrNoRows

	u.queriesMock.EXPECT().ValidateProductOrderReview(gomock.Any(), sql.NullInt32{Int32: int32(req.ProductId), Valid: true}).Return(database.Order{}, expectedErr)
	errResp := u.UsersService.AddReview(u.context, req, custId)
	u.NotNil(errResp)
	u.EqualValues(expectedErr.Error(), "no rows in result set")
}
