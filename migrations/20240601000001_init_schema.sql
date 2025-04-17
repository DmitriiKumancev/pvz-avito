-- +goose Up
-- +goose StatementBegin

----------------------------------------
-- Таблица пунктов выдачи заказов (ПВЗ)
----------------------------------------
-- Содержит информацию о пунктах выдачи в разрешенных городах.
-- Примеры использования: создание нового ПВЗ, проверка доступности города.
CREATE TABLE IF NOT EXISTS pvz (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    registration_date TIMESTAMP NOT NULL DEFAULT NOW(),
    city VARCHAR(255) NOT NULL CHECK (city IN ('Москва', 'Санкт-Петербург', 'Казань')) -- Ограничение на поддерживаемые города
);

----------------------------------------
-- Таблица приемок товаров
----------------------------------------
-- Отслеживает статусы приемки товаров в ПВЗ.
-- Статусы:
--   - in_progress: приемка открыта для добавления товаров
--   - close: приемка завершена, изменения невозможны
CREATE TABLE IF NOT EXISTS reception (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    date_time TIMESTAMP NOT NULL DEFAULT NOW(),
    pvz_id UUID NOT NULL REFERENCES pvz(id) ON DELETE CASCADE, -- При удалении ПВЗ удаляются все связанные приемки
    status VARCHAR(20) NOT NULL CHECK (status IN ('in_progress', 'close'))
);

-- Ускоряет поиск активных приемок (in_progress) для конкретного ПВЗ
CREATE INDEX IF NOT EXISTS idx_reception_pvz_status ON reception(pvz_id, status);

----------------------------------------
-- Таблица товаров
----------------------------------------
-- Хранит информацию о товарах в приемках.
-- Типы товаров: электроника, одежда, обувь.
-- Каждый товар привязан к одной приемке.
CREATE TABLE IF NOT EXISTS product (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    date_time TIMESTAMP NOT NULL DEFAULT NOW(),
    type VARCHAR(50) NOT NULL CHECK (type IN ('электроника', 'одежда', 'обувь')),
    reception_id UUID NOT NULL REFERENCES reception(id) ON DELETE CASCADE -- При удалении приемки удаляются все ее товары
);

-- Ускоряет фильтрацию товаров по принадлежности к приемке
CREATE INDEX IF NOT EXISTS idx_product_reception ON product(reception_id);

----------------------------------------
-- Управление порядком товаров (LIFO)
----------------------------------------
-- Определяет последовательность извлечения товаров:
-- Последний добавленный товар (с максимальным id) извлекается первым.
-- Используется для реализации логики "последним пришел - первым ушел".
CREATE TABLE IF NOT EXISTS product_sequence (
    id SERIAL PRIMARY KEY,
    reception_id UUID NOT NULL REFERENCES reception(id) ON DELETE CASCADE,
    product_id UUID NOT NULL UNIQUE REFERENCES product(id) ON DELETE CASCADE -- Гарантирует уникальность товара в последовательностях
);

-- Оптимизирует поиск последовательности товаров для конкретной приемки
CREATE INDEX IF NOT EXISTS idx_product_sequence_reception ON product_sequence(reception_id);

----------------------------------------
-- Управление пользователями
----------------------------------------
-- Роли:
--   - employee: сотрудник ПВЗ (работа с приемками и товарами)
--   - moderator: расширенные права (управление ПВЗ и пользователями)
CREATE TYPE user_role AS ENUM ('employee', 'moderator');

-- Хранит учетные данные и роли пользователей системы
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) NOT NULL UNIQUE, -- Логин пользователя
    password_hash VARCHAR(255) NOT NULL, -- Хэш пароля
    role user_role NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
/* 
ВАЖНО! Удаление таблиц производится в порядке, обратном созданию,
чтобы избежать ошибок зависимостей между внешними ключами.
*/
DROP TABLE IF EXISTS users;
DROP TYPE IF EXISTS user_role;
DROP TABLE IF EXISTS product_sequence;
DROP TABLE IF EXISTS product;
DROP TABLE IF EXISTS reception;
DROP TABLE IF EXISTS pvz;

-- +goose StatementEnd