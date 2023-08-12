package user

import (
	"context"
	"sushi-mart/common"
	"sushi-mart/internal/database"
)

type UsersService interface {
	CreateUser(context.Context, database.CreateCustomerParams) (*database.Customer, error)
	GetUser(context.Context, string) (*database.Customer, error)

	CreateUserWallet(context.Context, *CreateWalletReq, int) *common.ErrorResponse
	GetUserWallet(context.Context, int) (*GetWalletRes, *common.ErrorResponse)
	UpdateUserWallet(context.Context, *UpdateWalletReq, int) *common.ErrorResponse

	GetAllProducts(context.Context) (*GetAllProductsResp, *common.ErrorResponse)

	AddReview(context.Context, *AddReviewReq, int) *common.ErrorResponse
}

type UsersServiceImpl struct {
	Queries *database.Queries
}

type Cache struct {
	UsersService
}
type Validator struct {
	UsersService
}

type RoutesWrapper struct {
	UsersService
}

func New(Queries *database.Queries) *RoutesWrapper {
	return &RoutesWrapper{
		UsersService: &Validator{
			UsersService: &Cache{
				UsersService: &UsersServiceImpl{
					Queries: Queries,
				},
			},
		},
	}
}
