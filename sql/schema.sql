-- Define the Customers table

CREATE TABLE IF NOT EXISTS customers (
    id SERIAL PRIMARY KEY,
    username VARCHAR(100) NOT NULL,
    password VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    phone VARCHAR(20),
    address TEXT
);

-- Define the Products table
CREATE TABLE IF NOT EXISTS productItems( 
	id SERIAL PRIMARY KEY,
	name VARCHAR(50) NOT NULL,
    quantity INT NOT NULL,
	category VARCHAR(50) NOT NULL,
    unit_price DECIMAL(5,2) NOT NULL,
    date_added TIMESTAMP NOT NULL,
    date_modified TIMESTAMP NOT NULL 
);

-- Define the Wallet table
CREATE TABLE wallet(
	id SERIAL PRIMARY KEY,
	balance DECIMAL(10,3) NOT NULL,
	wallet_type VARCHAR(20) NOT NULL,
	date_added TIMESTAMP NOT NULL,
    date_modified TIMESTAMP NOT NULL,
	customer_id INT REFERENCES customers(id)
)
