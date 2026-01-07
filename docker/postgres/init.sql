-- Создаем таблицу пользователей
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    email VARCHAR(100),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Вставляем тестовых пользователей
INSERT INTO users (username, password_hash, email) VALUES
    ('admin', '$2a$10$N9qo8uLOickgx2ZMRZoMy.MrqK0YpY7lB6sY/4CJj7q7G', 'admin@example.com'),
    ('user1', '$2a$10$N9qo8uLOickgx2ZMRZoMy.MrqK0YpY7lB6sY/4CJj7q7G', 'user1@example.com')
ON CONFLICT (username) DO NOTHING;

-- Создаем таблицу токенов (если нужно)
CREATE TABLE IF NOT EXISTS tokens (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id),
    token VARCHAR(500) NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    is_revoked BOOLEAN DEFAULT FALSE
);

-- Создаем индекс для быстрого поиска по токену
CREATE INDEX IF NOT EXISTS idx_tokens_token ON tokens(token);
CREATE INDEX IF NOT EXISTS idx_tokens_user_id ON tokens(user_id);