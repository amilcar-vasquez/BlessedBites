CREATE TABLE recommendations (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id),
    menu_item_id INT REFERENCES menu_items(id),
    reason TEXT,
    created_at TIMESTAMP DEFAULT NOW()
);
