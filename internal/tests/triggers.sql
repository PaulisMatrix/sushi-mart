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

CREATE OR REPLACE FUNCTION revert_balance_fun() RETURNS trigger AS $revert_balance_fun$
     BEGIN
         UPDATE wallet SET balance = balance + NEW.total_amt WHERE NEW.order_status='CANCELLED';
         RETURN NEW;
     END;
$revert_balance_fun$ LANGUAGE plpgsql;

CREATE TRIGGER revert_balance AFTER UPDATE ON orders
     FOR EACH ROW EXECUTE FUNCTION revert_balance_fun();

CREATE OR REPLACE FUNCTION mark_prodreview_inactive_fun() RETURNS trigger AS $mark_prodreview_inactive_fun$
     BEGIN
         UPDATE productReviews pr SET pr.is_active=FALSE WHERE pr.product_id = OLD.id;
     END;
$mark_prodreview_inactive_fun$ LANGUAGE plpgsql;

CREATE TRIGGER mark_prodreview_inactive AFTER DELETE ON productItems
     FOR EACH ROW EXECUTE FUNCTION mark_prodreview_inactive_fun();

CREATE OR REPLACE FUNCTION mark_order_inactive_fun() RETURNS trigger AS $mark_order_inactive_fun$
     BEGIN
         UPDATE orders o SET o.is_active=FALSE WHERE o.product_id = OLD.id OR o.customer_id=OLD.id;
     END;
$mark_order_inactive_fun$ LANGUAGE plpgsql;


CREATE TRIGGER mark_order_inactive_prodItems AFTER DELETE ON productItems
     FOR EACH ROW EXECUTE FUNCTION mark_order_inactive_fun();

CREATE TRIGGER mark_order_inactive_cust AFTER DELETE ON customers
     FOR EACH ROW EXECUTE FUNCTION mark_order_inactive_fun();