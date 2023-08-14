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

               


4. ## TRIGGERS in the DB:

**Trigger to check for product quantity in productItems before customer places an order**:

```
CREATE OR REPLACE FUNCTION update_quantity_fun() RETURNS trigger AS $update_quantity_fun$
     DECLARE q NUMERIC;
     BEGIN
         SELECT quantity into q FROM productItems where id=NEW.product_id;
         IF  q < NEW.units THEN
             RAISE EXCEPTION 'no more product items left to purchase';
         END IF;

     	 UPDATE productItems set quantity=quantity - NEW.units where id=NEW.product_id;
	 RETURN NEW;
     END;
$update_quantity_fun$ LANGUAGE plpgsql;

CREATE TRIGGER update_quantity BEFORE INSERT ON orders
    FOR EACH ROW EXECUTE FUNCTION update_quantity_fun();
```

**Trigger to check for account balance before customer places an order**:

```
CREATE OR REPLACE FUNCTION update_balance_fun() RETURNS trigger AS $update_balance_fun$
     DECLARE total_amt NUMERIC;
     DECLARE bal NUMERIC;
     BEGIN
         SELECT (NEW.total_amt) INTO total_amt;
         SELECT w.balance INTO bal from wallet w where w.customer_id=NEW.customer_id;
         IF total_amt > bal THEN
                RAISE EXCEPTION 'account balance insufficient for the transaction';
         END IF;
       
         UPDATE wallet set balance=balance - total_amt;
         RETURN NEW; 
     END;
$update_balance_fun$ LANGUAGE plpgsql;

CREATE TRIGGER update_balance BEFORE INSERT ON orders
    FOR EACH ROW EXECUTE FUNCTION update_balance_fun();
```

**Trigger to revert balance after customer cancels an order**:

```
CREATE OR REPLACE FUNCTION revert_balance_fun() RETURNS trigger AS $revert_balance_fun$
     BEGIN
         UPDATE wallet SET balance = balance + NEW.total_amt WHERE NEW.order_status='CANCELLED';
         RETURN NEW;
     END;
$revert_balance_fun$ LANGUAGE plpgsql;

CREATE TRIGGER revert_balance AFTER UPDATE ON orders
     FOR EACH ROW EXECUTE FUNCTION revert_balance_fun();
```

**Trigger to mark productReview in_active as FALSE IF a productItem is deleted**:

```
CREATE OR REPLACE FUNCTION mark_prodreview_inactive_fun() RETURNS trigger AS $mark_prodreview_inactive_fun$
     BEGIN
         UPDATE productReviews pr SET pr.is_active=FALSE WHERE pr.product_id = OLD.id;
     END;
$mark_prodreview_inactive_fun$ LANGUAGE plpgsql;

CREATE TRIGGER mark_prodreview_inactive AFTER DELETE ON productItems
     FOR EACH ROW EXECUTE FUNCTION mark_prodreview_inactive_fun();
```

**Trigger to mark an order in_active as FALSE IF a productItem or a customer is deleted**:

```
CREATE OR REPLACE FUNCTION mark_order_inactive_fun() RETURNS trigger AS $mark_order_inactive_fun$
     BEGIN
         UPDATE orders o SET o.is_active=FALSE WHERE o.product_id = OLD.id OR o.customer_id=OLD.id;
     END;
$mark_order_inactive_fun$ LANGUAGE plpgsql;


CREATE TRIGGER mark_order_inactive_prodItems AFTER DELETE ON productItems
     FOR EACH ROW EXECUTE FUNCTION mark_order_inactive_fun();

CREATE TRIGGER mark_order_inactive_cust AFTER DELETE ON customers
     FOR EACH ROW EXECUTE FUNCTION mark_order_inactive_fun();
```

5. Orders are not deleted when a productItem or customer is deleted so as to maintain history. Rather they are marked as `is_active=false`