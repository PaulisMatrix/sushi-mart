## sushi-mart
Order Management System. 

1. ## Access the documentation here:
     `http://localhost:8080/api/v2/swagger/index.html`

2. ## Different modes to run the project

     a. Main server: `./main serve`

     b. Start the background consumer which consumes the orders placed: `./main consume`

3. ## Different routes:

     a.   Admin routes (Restricted to admins with BasicAuthentication):

          *Analytics Routes:

               i.   GET: /api/v1/admin/analytics/avg-cust-ratings - to check average customer ratings.

               ii.  GET: /api/v1/admin/analytics/top-orders-placed?limit={limit} - to get top orders placed with limit.
          
          *Inventory Routes:

               i.   GET: /api/v1/admin/inventory/all-products - to get all products in the inventory.

               ii.  POST: /api/v1/admin/inventory/add-product - to add a product in the inventory.

               iii. PATCH: /api/v1/admin/inventory/update-product - to update a product in the inventory.

               iv.  DELETE: /api/v1/admin/inventory/delete-product/{id} - to delete a particular product in the inventory.
     
     b.   Orders routes (Protected with JWT Authentication):

          i.   POST: /api/v1/orders/place-order - to place an order.

          ii.  POST: /api/v1/orders/cancel-order - to cancel an order.

          iii. GET: /api/v1/orders/get-orders - to get all customer orders.
     
     c.   User routes (Protected with JWT Authentication except for login/signup routes)

          i.   POST: /api/v1/users/signup - customer/user create account.

          ii.  POST: /api/v1/users/login - customer login and generate a JWT token.

          iii. POST: /api/v1/users/create-wallet - create customer wallet.

          iv.  GET: /api/v1/users/get-wallet - get customer wallet.

          v.   PATCH: /api/v1/users/update-wallet - update customer wallet.

          vi.  GET: /api/v1/users/all-products - products list from which customer can choose to place an order.

          vii. POST: /api/v1/users/add-review - add a review for the purchased product.

4. Orders are not deleted when a productItem or customer is deleted so as to maintain history. Rather they are marked as `is_active=false`

## TODO:

- Use pgx as the database driver instead of the default one.
- Use docker for easy setup.
- Add more tests.
- An analytics dashboard displaying order trends and product popularity.(use websockets here)
- A frontend with good UI. (use react?)
- Use a much reliant rabbitmq as the msg broker instead of redis.