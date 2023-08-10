package user

import (
	"context"
	"sushi-mart/internal/database"
)

// Validator acts like a wrapper around UsersServiceImpl.GetUser
// It does validations for incoming requests
func (v *Validator) GetUser(ctx context.Context, emailAddres string) (*database.Customer, error) {
	return v.UsersService.GetUser(ctx, emailAddres)
}

func (u *UsersServiceImpl) GetUser(ctx context.Context, emailAddres string) (*database.Customer, error) {
	resp, err := u.Queries.GetCustomer(ctx, emailAddres)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}
