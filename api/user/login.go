package user

import (
	"context"
	"sushi-mart/common"
	"sushi-mart/internal/database"

	"github.com/sirupsen/logrus"
)

// Validator acts like a wrapper around UsersServiceImpl.GetUser
// It does validations for incoming requests
func (v *Validator) GetUser(ctx context.Context, emailAddres string) (*database.Customer, error) {
	return v.UsersService.GetUser(ctx, emailAddres)
}

func (u *UsersServiceImpl) GetUser(ctx context.Context, emailAddres string) (*database.Customer, error) {
	logger := common.ExtractLoggerUnsafe(ctx).WithFields(logrus.Fields{"method": "GetUser", "request": emailAddres})

	resp, err := u.Queries.GetCustomer(ctx, emailAddres)
	if err != nil {
		logger.WithError(err).Error("error in getting the user")
		return nil, err
	}

	return &resp, nil
}
