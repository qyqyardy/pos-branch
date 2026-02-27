CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name TEXT NOT NULL,
    email TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    role TEXT NOT NULL
);

CREATE TABLE products (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name TEXT NOT NULL,
    price BIGINT NOT NULL
);

CREATE TABLE orders (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    cashier_id UUID REFERENCES users(id),
    total BIGINT NOT NULL,
    created_at TIMESTAMP DEFAULT now()
);

CREATE TABLE order_items (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    order_id UUID REFERENCES orders(id) ON DELETE CASCADE,
    product_id UUID REFERENCES products(id),
    qty INT NOT NULL,
    price BIGINT NOT NULL
);

-- password: 123456
INSERT INTO users (name,email,password_hash,role)
VALUES (
    'Admin',
    'admin@pos.com',
    '$2a$10$m8nlpInZRx1xlCTir.wyNOQ9YvJ3CpP9AZbw8PM46Z8V61/TJL6KS',
    'admin'
);

-- ======================
-- DUMMY PRODUCTS (WARKOP)
-- ======================
INSERT INTO products (name, price) VALUES
    ('Kopi Hitam', 8000),
    ('Kopi Susu', 10000),
    ('Kopi Susu Gula Aren', 14000),
    ('Teh Manis', 7000),
    ('Es Teh Manis', 7000),
    ('Jeruk Hangat', 9000),
    ('Es Jeruk', 9000),
    ('Air Mineral', 5000),
    ('Soda Gembira', 14000),
    ('Indomie Goreng', 12000),
    ('Indomie Kuah', 12000),
    ('Indomie Goreng Double', 16000),
    ('Mie Dok Dok', 20000),
    ('Nasi Goreng Warkop', 22000),
    ('Roti Bakar Coklat', 15000),
    ('Roti Bakar Keju', 16000),
    ('Pisang Goreng', 12000),
    ('Tahu Goreng', 10000),
    ('Tempe Mendoan', 11000),
    ('Sate Taichan', 25000);
