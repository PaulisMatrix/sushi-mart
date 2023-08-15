package user

import (
	"context"
	"sushi-mart/common"
	"sushi-mart/internal/database"
)

type UsersService interface {
	CreateUser(context.Context, *SignUpReq) *common.ErrorResponse
	GetUser(context.Context, *LoginReq) (*CustomerInfo, *common.ErrorResponse)

	CreateUserWallet(context.Context, *CreateWalletReq, int) *common.ErrorResponse
	GetUserWallet(context.Context, int) (*GetWalletRes, *common.ErrorResponse)
	UpdateUserWallet(context.Context, *UpdateWalletReq, int) *common.ErrorResponse

	GetAllProducts(context.Context) (*GetAllProductsResp, *common.ErrorResponse)

	AddReview(context.Context, *AddReviewReq, int) *common.ErrorResponse
}

type UsersServiceImpl struct {
	Queries database.Querier
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
