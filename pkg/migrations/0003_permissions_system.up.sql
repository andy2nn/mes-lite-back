-- =========================
-- РАЗРЕШЕНИЯ (ПРАВА)
-- =========================
CREATE TABLE permissions (
    id BIGSERIAL PRIMARY KEY,
    code VARCHAR(100) NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    category VARCHAR(100),
    created_at TIMESTAMP DEFAULT NOW()
);

-- =========================
-- РОЛИ → РАЗРЕШЕНИЯ (многие-ко-многим)
-- =========================
CREATE TABLE role_permissions (
    id BIGSERIAL PRIMARY KEY,
    role_id BIGINT NOT NULL,
    permission_id BIGINT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    CONSTRAINT fk_role_permissions_role FOREIGN KEY(role_id) REFERENCES roles(id) ON DELETE CASCADE,
    CONSTRAINT fk_role_permissions_permission FOREIGN KEY(permission_id) REFERENCES permissions(id) ON DELETE CASCADE,
    UNIQUE(role_id, permission_id)
);

-- =========================
-- ОСНОВНЫЕ РАЗРЕШЕНИЯ ДЛЯ MES
-- =========================
INSERT INTO permissions (code, name, description, category) VALUES
-- Заказы
('order.view', 'Просмотр заказов', 'Просмотр рабочих заказов', 'Заказы'),
('order.create', 'Создание заказов', 'Создание новых рабочих заказов', 'Заказы'),
('order.edit', 'Редактирование заказов', 'Изменение данных заказов', 'Заказы'),
('order.delete', 'Удаление заказов', 'Удаление рабочих заказов', 'Заказы'),
('order.status', 'Изменение статуса', 'Изменение статуса заказа', 'Заказы'),
('order.comment', 'Комментирование', 'Добавление комментариев к заказам', 'Заказы'),

-- Оборудование
('machine.view', 'Просмотр оборудования', 'Просмотр машин и линий', 'Оборудование'),
('machine.edit', 'Управление оборудованием', 'Добавление и изменение оборудования', 'Оборудование'),
('machine.status', 'Управление статусом', 'Изменение статуса оборудования', 'Оборудование'),

-- Продукция
('product.view', 'Просмотр продукции', 'Просмотр списка продукции', 'Продукция'),
('product.edit', 'Управление продукцией', 'Создание и изменение продукции', 'Продукция'),
('product.delete', 'Удаление продукции', 'Удаление видов продукции', 'Продукция'),
('product.instance.view', 'Просмотр экземпляров', 'Просмотр экземпляров продукции', 'Продукция'),
('product.instance.create', 'Создание экземпляров', 'Создание штрихкодов продукции', 'Продукция'),

-- Этапы производства
('stage.view', 'Просмотр этапов', 'Просмотр этапов производства', 'Производство'),
('stage.edit', 'Управление этапами', 'Настройка этапов продукции', 'Производство'),
('stage.execute', 'Выполнение этапов', 'Отметка выполнения производственных этапов', 'Производство'),

-- Производственные операции
('production.start', 'Запуск производства', 'Запуск выполнения заказа', 'Производство'),
('production.complete', 'Завершение производства', 'Завершение выполнения заказа', 'Производство'),

-- Инциденты
('incident.view', 'Просмотр инцидентов', 'Просмотр всех инцидентов', 'Инциденты'),
('incident.create', 'Создание инцидентов', 'Регистрация новых инцидентов', 'Инциденты'),
('incident.edit', 'Редактирование инцидентов', 'Изменение данных инцидентов', 'Инциденты'),
('incident.resolve', 'Закрытие инцидентов', 'Разрешение и закрытие инцидентов', 'Инциденты'),

-- Качество
('quality.view', 'Просмотр контроля', 'Просмотр контроля качества', 'Качество'),
('quality.inspect', 'Контроль качества', 'Проведение контроля качества продукции', 'Качество'),
('quality.edit', 'Управление качеством', 'Изменение статусов качества', 'Качество'),

-- Расписание
('schedule.view', 'Просмотр расписания', 'Просмотр производственного расписания', 'Планирование'),
('schedule.edit', 'Управление расписанием', 'Создание и изменение расписания', 'Планирование'),

-- Пользователи
('user.view', 'Просмотр пользователей', 'Просмотр списка пользователей', 'Администрирование'),
('user.edit', 'Управление пользователями', 'Создание и изменение пользователей', 'Администрирование'),
('user.delete', 'Удаление пользователей', 'Удаление пользователей из системы', 'Администрирование'),

-- Роли и права
('role.view', 'Просмотр ролей', 'Просмотр списка ролей', 'Администрирование'),
('role.edit', 'Управление ролями', 'Создание и изменение ролей', 'Администрирование'),
('permission.view', 'Просмотр разрешений', 'Просмотр списка разрешений', 'Администрирование'),
('permission.assign', 'Назначение разрешений', 'Назначение разрешений ролям', 'Администрирование'),

-- Логи и мониторинг
('log.view', 'Просмотр логов', 'Просмотр системных логов', 'Мониторинг'),
('report.view', 'Просмотр отчетов', 'Доступ к отчетам и аналитике', 'Отчеты'),
('dashboard.view', 'Просмотр дашборда', 'Доступ к панели управления', 'Отчеты');

-- =========================
-- ОСНОВНЫЕ РОЛИ (если ещё не существуют)
-- =========================
INSERT INTO roles (name) VALUES 
('Администратор'),
('Менеджер'),
('Рабочий'),
('Контролер качества'),
('Планировщик')
ON CONFLICT (name) DO NOTHING;

-- =========================
-- ФУНКЦИЯ ДЛЯ НАЗНАЧЕНИЯ ПРАВ
-- =========================
CREATE OR REPLACE FUNCTION assign_role_permissions(role_name VARCHAR, permission_codes VARCHAR[])
RETURNS void AS $$
DECLARE
    r_id BIGINT;
    p_id BIGINT;
    perm_code VARCHAR;
BEGIN
    SELECT id INTO r_id FROM roles WHERE name = role_name;
    
    FOREACH perm_code IN ARRAY permission_codes
    LOOP
        SELECT id INTO p_id FROM permissions WHERE code = perm_code;
        
        IF r_id IS NOT NULL AND p_id IS NOT NULL THEN
            INSERT INTO role_permissions (role_id, permission_id)
            VALUES (r_id, p_id)
            ON CONFLICT (role_id, permission_id) DO NOTHING;
        END IF;
    END LOOP;
END;
$$ LANGUAGE plpgsql;

-- =========================
-- НАЗНАЧЕНИЕ ПРАВ РОЛЯМ
-- =========================

-- 1. АДМИНИСТРАТОР (все права)
SELECT assign_role_permissions('Администратор', ARRAY[
    'order.view', 'order.create', 'order.edit', 'order.delete', 'order.status', 'order.comment',
    'machine.view', 'machine.edit', 'machine.status',
    'product.view', 'product.edit', 'product.delete', 'product.instance.view', 'product.instance.create',
    'stage.view', 'stage.edit', 'stage.execute',
    'production.start', 'production.complete',
    'incident.view', 'incident.create', 'incident.edit', 'incident.resolve',
    'quality.view', 'quality.inspect', 'quality.edit',
    'schedule.view', 'schedule.edit',
    'user.view', 'user.edit', 'user.delete',
    'role.view', 'role.edit',
    'permission.view', 'permission.assign',
    'log.view',
    'report.view',
    'dashboard.view'
]);

-- 2. МЕНЕДЖЕР (оперативное управление)
SELECT assign_role_permissions('Менеджер', ARRAY[
    'order.view', 'order.create', 'order.edit', 'order.status', 'order.comment',
    'machine.view',
    'product.view', 'product.edit',
    'stage.view',
    'production.start', 'production.complete',
    'incident.view', 'incident.create', 'incident.edit',
    'quality.view',
    'schedule.view', 'schedule.edit',
    'user.view',
    'report.view',
    'dashboard.view'
]);

-- 3. РАБОЧИЙ (исполнение)
SELECT assign_role_permissions('Рабочий', ARRAY[
    'order.view',
    'machine.view',
    'product.view',
    'stage.view', 'stage.execute',
    'production.start', 'production.complete',
    'incident.create',
    'schedule.view'
]);

-- 4. КОНТРОЛЕР КАЧЕСТВА
SELECT assign_role_permissions('Контролер качества', ARRAY[
    'order.view',
    'product.view', 'product.instance.view',
    'stage.view',
    'incident.create',
    'quality.view', 'quality.inspect', 'quality.edit',
    'report.view'
]);

-- 5. ПЛАНИРОВЩИК
SELECT assign_role_permissions('Планировщик', ARRAY[
    'order.view', 'order.create', 'order.edit',
    'machine.view',
    'product.view',
    'schedule.view', 'schedule.edit',
    'report.view',
    'dashboard.view'
]);

-- =========================
-- ФУНКЦИЯ ПРОВЕРКИ ПРАВ
-- =========================
CREATE OR REPLACE FUNCTION has_permission(user_id BIGINT, permission_code VARCHAR)
RETURNS BOOLEAN AS $$
BEGIN
    RETURN EXISTS (
        SELECT 1
        FROM users u
        JOIN roles r ON u.role_id = r.id
        JOIN role_permissions rp ON r.id = rp.role_id
        JOIN permissions p ON rp.permission_id = p.id
        WHERE u.id = user_id 
          AND p.code = permission_code
    );
END;
$$ LANGUAGE plpgsql;

-- =========================
-- ВСПОМОГАТЕЛЬНЫЕ ПРЕДСТАВЛЕНИЯ
-- =========================
CREATE VIEW v_role_permissions AS
SELECT 
    r.id as role_id,
    r.name as role_name,
    p.id as permission_id,
    p.code as permission_code,
    p.name as permission_name,
    p.category
FROM roles r
JOIN role_permissions rp ON r.id = rp.role_id
JOIN permissions p ON rp.permission_id = p.id
ORDER BY r.name, p.category, p.name;

CREATE VIEW v_user_permissions AS
SELECT 
    u.id as user_id,
    u.username,
    u.full_name,
    r.name as role_name,
    p.code as permission_code,
    p.name as permission_name,
    p.category
FROM users u
JOIN roles r ON u.role_id = r.id
JOIN role_permissions rp ON r.id = rp.role_id
JOIN permissions p ON rp.permission_id = p.id
ORDER BY u.username, p.category, p.name;

-- =========================
-- ИНДЕКСЫ ДЛЯ ПРОИЗВОДИТЕЛЬНОСТИ
-- =========================
CREATE INDEX idx_permissions_code ON permissions(code);
CREATE INDEX idx_role_permissions_composite ON role_permissions(role_id, permission_id);

-- =========================
-- НАЗНАЧЕНИЕ РОЛИ ПО УМОЛЧАНИЮ СУЩЕСТВУЮЩИМ ПОЛЬЗОВАТЕЛЯМ
-- =========================
DO $$
DECLARE
    default_role_id BIGINT;
BEGIN
    SELECT id INTO default_role_id FROM roles WHERE name = 'Рабочий';
    
    IF default_role_id IS NOT NULL THEN
        UPDATE users 
        SET role_id = default_role_id 
        WHERE role_id IS NULL;
        
        INSERT INTO system_logs (level, source, message, created_at)
        VALUES ('INFO', 'Migration', 
                'Назначена роль по умолчанию пользователям без роли', 
                NOW());
    END IF;
END $$;