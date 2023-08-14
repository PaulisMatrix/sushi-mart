package orders

type OrderStatus string

const (
	PROCESSING OrderStatus = "PROCESSING"
	DELIVERED  OrderStatus = "DELIVERED"
	SHIPPED    OrderStatus = "SHIPPED"
	CANCELLED  OrderStatus = "CANCELLED"
)

type PlaceOrderReq struct {
	CustomerID  int    `json:"cust_id,omitempty"`
	Units       int    `json:"units" binding:"required"`
	ProductId   int    `json:"product_id" binding:"required"`
	PaymentType string `json:"payment_type" binding:"required"`
}

type UpdateOrderReq struct {
	OrderId int `json:"order_id" binding:"required"`
}

type GetAllOrdersResp struct {
	Orders []GetAllOrders `json:"orders"`
}

type GetAllOrders struct {
	OrderDate   string  `json:"order_date"`
	OrderStatus string  `json:"order_status"`
	TotalAmount float64 `json:"total_amount"`
	Username    string  `json:"username"`
	ProductName string  `json:"product_name"`
}
