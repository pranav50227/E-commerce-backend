-- Seed data
INSERT INTO users (id, username, password, name, email) 
VALUES ('user1', 'demouser', 'demopassword', 'Demo User', 'demo@example.com')
ON CONFLICT DO NOTHING;

INSERT INTO products (id, name, description, price, category) 
VALUES ('prod1', 'Wireless Headset', 'High-end wireless headphones', 150.00, 'Electronics')
ON CONFLICT DO NOTHING;

INSERT INTO inventory (product_id, quantity) 
VALUES ('prod1', 50)
ON CONFLICT DO NOTHING;
