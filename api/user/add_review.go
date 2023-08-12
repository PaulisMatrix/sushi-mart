package user

import (
	"context"
	"database/sql"
	"net/http"
	"sushi-mart/common"
	"sushi-mart/internal/database"
	"time"

	"github.com/sirupsen/logrus"
)

func (v *Validator) AddReview(ctx context.Context, req *AddReviewReq, custId int) *common.ErrorResponse {
	return v.UsersService.AddReview(ctx, req, custId)
}

func (u *UsersServiceImpl) AddReview(ctx context.Context, req *AddReviewReq, custId int) *common.ErrorResponse {
	logger := common.ExtractLoggerUnsafe(ctx).WithFields(logrus.Fields{"method": "AddReview", "request": req})
	dbParams := database.AddReviewParams{
		Rating:     int32(req.Rating),
		ReviewText: req.ReviewText,
		ReviewDate: time.Now().Local(),
		CustomerID: sql.NullInt32{Int32: int32(custId), Valid: true},
		ProductID:  sql.NullInt32{Int32: int32(req.ProductId), Valid: true},
	}

	err := u.Queries.AddReview(ctx, dbParams)
	if err != nil {
		logger.WithError(err).Error("error in adding a new customer review")
		return &common.ErrorResponse{
			Status:  http.StatusInternalServerError,
			Message: "internal server error",
		}
	}
	return nil
}
