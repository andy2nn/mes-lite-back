-- УДАЛЕНИЕ В ОБРАТНОМ ПОРЯДКЕ

-- Удаляем представления
DROP VIEW IF EXISTS v_user_permissions;
DROP VIEW IF EXISTS v_role_permissions;

-- Удаляем функцию
DROP FUNCTION IF EXISTS has_permission(BIGINT, VARCHAR);
DROP FUNCTION IF EXISTS assign_role_permissions(VARCHAR, VARCHAR[]);

-- Удаляем связи ролей с разрешениями
DELETE FROM role_permissions;

-- Удаляем таблицы
DROP TABLE IF EXISTS role_permissions;
DROP TABLE IF EXISTS permissions;

-- Удаляем добавленные роли (если они были добавлены этой миграцией)
-- ВАЖНО: Эта операция может удалить данные, если роли уже использовались!
-- DELETE FROM roles WHERE name IN ('Контролер качества', 'Планировщик');