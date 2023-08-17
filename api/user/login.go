package user

import (
	"context"
	"net/http"
	"sushi-mart/common"

	"github.com/jackc/pgx/v5"
	"github.com/sirupsen/logrus"
)

// Validator acts like a wrapper around UsersServiceImpl.GetUser
// It does validations for incoming requests
func (v *Validator) GetUser(ctx context.Context, req *LoginReq) (*CustomerInfo, *common.ErrorResponse) {
	return v.UsersService.GetUser(ctx, req)
}

func (u *UsersServiceImpl) GetUser(ctx context.Context, req *LoginReq) (*CustomerInfo, *common.ErrorResponse) {
	logger := common.ExtractLoggerUnsafe(ctx).WithFields(logrus.Fields{"method": "GetUser", "request": req})

	resp, err := u.Queries.GetCustomer(ctx, req.Email)

	if err != nil {
		if err == pgx.ErrNoRows {
			logger.WithError(err).Error("customer not found")
			return nil, &common.ErrorResponse{
				Status:  http.StatusBadRequest,
				Message: "need to signup first",
			}
		}
		logger.WithError(err).Error("failed to get the customer")
		return nil, &common.ErrorResponse{
			Status:  http.StatusInternalServerError,
			Message: "internal server error",
		}
	}

	return &CustomerInfo{
		CustId:   int(resp.ID),
		Username: resp.Username,
		Password: resp.Password,
		Email:    resp.Email,
		Phone:    resp.Phone.String,
		Address:  resp.Address.String,
	}, nil
}
