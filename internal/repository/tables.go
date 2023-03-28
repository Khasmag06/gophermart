package repository

var (
	createUsersTable = `CREATE TABLE IF NOT EXISTS users(
            				id serial PRIMARY KEY,
		    				login VARCHAR(255) NOT NULL UNIQUE,
							password text NOT NULL)`

	createOrdersTable = `CREATE TABLE IF NOT EXISTS orders(
    						order_id serial PRIMARY KEY,
            				order_num VARCHAR(25) NOT NULL UNIQUE,
            				user_id INTEGER NOT NULL,
		    				status VARCHAR(25) NOT NULL,
							accrual INTEGER DEFAULT 0,
            				uploaded_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
            				FOREIGN KEY (user_id) REFERENCES users (id))`

	createWithdrawsTable = `CREATE TABLE IF NOT EXISTS withdrawals(
               				   id serial PRIMARY KEY,
    						   user_id INTEGER NOT NULL,
							   order_num VARCHAR(25) NOT NULL UNIQUE,
							   sum INTEGER DEFAULT 0,
							   processed_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
							   FOREIGN KEY (user_id) REFERENCES users (id))`

	tables = []string{createUsersTable, createOrdersTable, createWithdrawsTable}
)
