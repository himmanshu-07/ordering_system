PRAGMA foreign_keys=OFF;
BEGIN TRANSACTION;
CREATE TABLE tables (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		table_number INTEGER NOT NULL UNIQUE,
		qr_code TEXT NOT NULL UNIQUE
	);
CREATE TABLE categories (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		display_order INTEGER DEFAULT 0
	);
INSERT INTO categories VALUES(1,'Burger',1);
INSERT INTO categories VALUES(2,'Coke',3);
INSERT INTO categories VALUES(3,'Pizza',2);
CREATE TABLE products (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		category_id INTEGER NOT NULL,
		name TEXT NOT NULL,
		description TEXT,
		price REAL NOT NULL,
		image_url TEXT,
		is_available BOOLEAN DEFAULT 1,
		FOREIGN KEY (category_id) REFERENCES categories(id)
	);
INSERT INTO products VALUES(1,1,'Burger','',40.0,'',1);
INSERT INTO products VALUES(2,2,'Pizza','',80.0,'',1);
INSERT INTO products VALUES(3,3,'Coke ','',30.0,'',1);
INSERT INTO products VALUES(4,2,'Dal Fry with Steamed Rice','North Indian Simple Thali',150.0,'',1);
INSERT INTO products VALUES(6,3,'Kadhi Chawal','',170.0,'',1);
INSERT INTO products VALUES(7,1,'Test Case','',3.0,'',1);
CREATE TABLE orders (
		id TEXT PRIMARY KEY,
		table_id INTEGER NOT NULL,
		status TEXT NOT NULL DEFAULT 'pending',
		total REAL NOT NULL DEFAULT 0,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP, payment_method TEXT NOT NULL DEFAULT 'cash',
		FOREIGN KEY (table_id) REFERENCES tables(id)
	);
INSERT INTO orders VALUES('6823a20a-9089-447c-a4ed-5b013b9a2fbe',1,'served',120.0,'2026-08-31 10:37:29','cash');
INSERT INTO orders VALUES('92c6b05a-09ba-4716-9256-b5e759fac739',1,'served',139.99999999999999999,'2026-08-31 10:38:35','upi');
INSERT INTO orders VALUES('0a2cc1f2-1ef3-4c86-bd44-202123ffc1fb',2,'served',160.0,'2026-08-31 10:39:10','cash');
INSERT INTO orders VALUES('f72a8113-d927-491b-b17c-d2c5fb88a6c0',1,'served',70.0,'2026-08-31 14:51:00','cash');
INSERT INTO orders VALUES('76aa767c-c729-49e9-adaa-1990face2d5d',1,'served',70.0,'2026-08-31 15:26:37','cash');
INSERT INTO orders VALUES('3404cd7d-92eb-4ee0-897b-ad2eba248f5b',1,'pending',3.0,'2026-08-31 16:30:53','upi');
INSERT INTO orders VALUES('81b59615-79d6-45bc-9d20-24dbd08a545d',2,'pending',3.0,'2026-08-31 16:38:02','cash');
CREATE TABLE order_items (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		order_id TEXT NOT NULL,
		product_id INTEGER NOT NULL,
		quantity INTEGER NOT NULL,
		unit_price REAL NOT NULL,
		notes TEXT,
		FOREIGN KEY (order_id) REFERENCES orders(id),
		FOREIGN KEY (product_id) REFERENCES products(id)
	);
INSERT INTO order_items VALUES(1,'6823a20a-9089-447c-a4ed-5b013b9a2fbe',1,1,40.0,'');
INSERT INTO order_items VALUES(2,'6823a20a-9089-447c-a4ed-5b013b9a2fbe',2,1,80.0,'');
INSERT INTO order_items VALUES(3,'92c6b05a-09ba-4716-9256-b5e759fac739',2,1,80.0,'');
INSERT INTO order_items VALUES(4,'92c6b05a-09ba-4716-9256-b5e759fac739',3,2,30.0,'');
INSERT INTO order_items VALUES(5,'0a2cc1f2-1ef3-4c86-bd44-202123ffc1fb',1,4,40.0,'');
INSERT INTO order_items VALUES(6,'f72a8113-d927-491b-b17c-d2c5fb88a6c0',1,1,40.0,'');
INSERT INTO order_items VALUES(7,'f72a8113-d927-491b-b17c-d2c5fb88a6c0',3,1,30.0,'');
INSERT INTO order_items VALUES(8,'76aa767c-c729-49e9-adaa-1990face2d5d',1,1,40.0,'');
INSERT INTO order_items VALUES(9,'76aa767c-c729-49e9-adaa-1990face2d5d',3,1,30.0,'');
INSERT INTO order_items VALUES(10,'3404cd7d-92eb-4ee0-897b-ad2eba248f5b',7,1,3.0,'');
INSERT INTO order_items VALUES(11,'81b59615-79d6-45bc-9d20-24dbd08a545d',7,1,3.0,'');
CREATE TABLE users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT NOT NULL UNIQUE,
		password_hash TEXT NOT NULL,
		role TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
INSERT INTO users VALUES(1,'orange','$2a$10$VQeMStGgtNJPhsp2kFsgkuEcuq15sFdnBVxL7e2CzMhvQUkBR7dBW','admin','2026-08-12 11:16:36');
INSERT INTO users VALUES(2,'mango','$2a$10$ZBTQ9QgQR/.LFHNWVHTh2uQ3w5Ig0Ou6HUQhwEKYFbUdqeBH9nQg.','kitchen','2026-08-12 11:17:17');
CREATE TABLE settings (
		key TEXT PRIMARY KEY,
		value TEXT NOT NULL
	);
INSERT INTO settings VALUES('total_tables','5');
DELETE FROM sqlite_sequence;
INSERT INTO sqlite_sequence VALUES('users',2);
INSERT INTO sqlite_sequence VALUES('products',7);
INSERT INTO sqlite_sequence VALUES('categories',3);
INSERT INTO sqlite_sequence VALUES('order_items',11);
COMMIT;
