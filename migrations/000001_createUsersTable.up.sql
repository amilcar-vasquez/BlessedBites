CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    email TEXT UNIQUE NOT NULL,
    full_name TEXT NOT NULL,
    phone_no TEXT,
    password_hash TEXT NOT NULL,
    role TEXT NOT NULL, -- 'admin' or 'user'
    created_at TIMESTAMP DEFAULT NOW()
);