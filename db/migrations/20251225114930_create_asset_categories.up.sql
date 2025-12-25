CREATE TABLE IF NOT EXISTS asset_categories (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL
);

INSERT INTO asset_categories (name) VALUES
('Bond'),
('Cash'),
('Crypto'),
('Gold'),
('Mutual Fund'),
('Property'),
('Savings'),
('Stock'),
('Other');
