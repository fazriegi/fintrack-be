CREATE TABLE IF NOT EXISTS assets (
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT NOT NULL,
    name VARCHAR(100) NOT NULL,
    category_id INT NOT NULL,
    amount DECIMAL(15,2),
    purchase_price DECIMAL(18,2),
    status ENUM('active', 'sold', 'archived'),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    modified_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (category_id) REFERENCES user_asset_categories(id)
);

CREATE INDEX idx_assets_name ON assets(name);
CREATE INDEX idx_assets_status ON assets(status);

CREATE TABLE IF NOT EXISTS asset_histories (
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT NOT NULL,
    name VARCHAR(100) NOT NULL,
    category_id INT NOT NULL,
    amount DECIMAL(15,2),
    purchase_price DECIMAL(18,2),
    status ENUM('active', 'sold', 'archived'),
    created_at TIMESTAMP,
    modified_at TIMESTAMP,
    record_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (category_id) REFERENCES user_asset_categories(id)
);

CREATE INDEX idx_asset_histories_name ON asset_histories(name);
CREATE INDEX idx_asset_histories_status ON asset_histories(status);