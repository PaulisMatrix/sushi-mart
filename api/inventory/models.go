package inventory

type AddProductReq struct {
	Name      string  `json:"name" binding:"required"`
	Quantity  int32   `json:"quantity" binding:"required"`
	Category  string  `json:"category" binding:"required"`
	UnitPrice float64 `json:"unit_price" binding:"required"`
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

type UpdateProductReq struct {
	Name      string  `json:"name,omitempty" binding:"omitempty"`
	Quantity  int32   `json:"quantity,omitempty" binding:"omitempty"`
	Category  string  `json:"category,omitempty" binding:"omitempty"`
	UnitPrice float64 `json:"unit_price,omitempty" binding:"omitempty"`
}

type pathParameters struct {
	ID int64 `uri:"id" binding:"required"`
}
