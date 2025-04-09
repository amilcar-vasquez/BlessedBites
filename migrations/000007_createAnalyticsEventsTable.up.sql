CREATE TABLE analytics_events (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id),
    action TEXT NOT NULL,
    menu_item_id INT REFERENCES menu_items(id),
    created_at TIMESTAMP DEFAULT NOW()
);