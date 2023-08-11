package user

type SignUpReq struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Phone    string `json:"phone" binding:"required"`
	Address  string `json:"address" binding:"required"`
}

type LoginReq struct {
	Password string `json:"password" binding:"required"`
	Email    string `json:"email" binding:"required"`
}

type LoginResp struct {
	Username string `json:"username"`
	JWTToken string `json:"token"`
}

type CreateWalletReq struct {
	Balance    float64 `json:"balance" binding:"required"`
	WalletType string  `json:"wallet_type" binding:"required"`
}

type GetWalletRes struct {
	Username   string  `json:"username"`
	Balance    float64 `json:"balance"`
	WalletType string  `json:"wallet_type"`
}

type UpdateWalletReq struct {
	Balance    float64 `json:"balance,omitempty" binding:"omitempty"`
	WalletType string  `json:"wallet_type,omitempty" binding:"omitempty"`
}

type ProductResp struct {
	Name         string  `json:"name"`
	Quantity     int32   `json:"quantity"`
	Category     string  `json:"category"`
	UnitPrice    float64 `json:"unit_price"`
	DateAdded    string  `json:"date_added"`
	DateModified string  `json:"date_modified"`
}

type GetAllProductsResp struct {
	Products []ProductResp `json:"products"`
}
