package analytics

type AvgCustomerRatingsResp struct {
	AvgRatings []AvgCustomerRatings `json:"avg_ratings"`
}

type AvgCustomerRatings struct {
	ProductName     string  `json:"product_name"`
	ProductCategory string  `json:"product_category"`
	AvgRating       float64 `json:"avg_rating"`
}

type MostOrdersPlacedResp struct {
	OrdersPlaced []OrdersPlaced `json:"orders_placed"`
}

type OrdersPlaced struct {
	Username   string `json:"username"`
	Email      string `json:"email"`
	OrderCount int    `json:"order_count"`
}
