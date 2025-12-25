-- =========================
-- СПРАВОЧНИКИ
-- =========================
CREATE TABLE order_priorities (
    id SERIAL PRIMARY KEY,
    name VARCHAR NOT NULL UNIQUE
);

CREATE TABLE order_statuses (
    id SERIAL PRIMARY KEY,
    name VARCHAR NOT NULL UNIQUE
);

CREATE TABLE machine_statuses (
    id SERIAL PRIMARY KEY,
    name VARCHAR NOT NULL UNIQUE
);

CREATE TABLE incident_severities (
    id SERIAL PRIMARY KEY,
    name VARCHAR NOT NULL UNIQUE
);

CREATE TABLE quality_statuses (
    id SERIAL PRIMARY KEY,
    name VARCHAR NOT NULL UNIQUE
);

-- =========================
-- РОЛИ
-- =========================
CREATE TABLE roles (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR NOT NULL UNIQUE
);

-- =========================
-- ПОЛЬЗОВАТЕЛИ
-- =========================
CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    username VARCHAR NOT NULL UNIQUE,
    password TEXT NOT NULL,
    full_name VARCHAR,
    role_id BIGINT,
    created_at TIMESTAMP DEFAULT NOW(),
    CONSTRAINT fk_users_roles FOREIGN KEY(role_id) REFERENCES roles(id)
);

-- =========================
-- ЛИНИИ ПРОИЗВОДСТВА
-- =========================
CREATE TABLE lines (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR NOT NULL
);

-- =========================
-- МАШИНЫ
-- =========================
CREATE TABLE machines (
    id BIGSERIAL PRIMARY KEY,
    code VARCHAR NOT NULL UNIQUE,
    name VARCHAR NOT NULL,
    line_id BIGINT,
    status_id INT DEFAULT 1,
    CONSTRAINT fk_machines_lines FOREIGN KEY(line_id) REFERENCES lines(id),
    CONSTRAINT fk_machines_status FOREIGN KEY(status_id) REFERENCES machine_statuses(id)
);

-- =========================
-- ПРОДУКТЫ
-- =========================
CREATE TABLE products (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR NOT NULL,
    sku VARCHAR NOT NULL UNIQUE,
    tech_cycle_min INTEGER,
    created_by BIGINT,
    CONSTRAINT fk_products_users FOREIGN KEY(created_by) REFERENCES users(id)
);

-- =========================
-- ЭТАПЫ ПРОДУКТА
-- =========================
CREATE TABLE stages (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR NOT NULL,
    description TEXT,
    is_strict_sequence BOOLEAN DEFAULT TRUE
);

CREATE TABLE product_stages (
    id BIGSERIAL PRIMARY KEY,
    product_id BIGINT,
    stage_id BIGINT,
    stage_order INTEGER,
    CONSTRAINT fk_product_stages_products FOREIGN KEY(product_id) REFERENCES products(id),
    CONSTRAINT fk_product_stages_stages FOREIGN KEY(stage_id) REFERENCES stages(id),
    UNIQUE(product_id, stage_id)
);

-- =========================
-- ИНСТАНСЫ ПРОДУКТОВ (штрихкоды)
-- =========================
CREATE TABLE product_instances (
    id BIGSERIAL PRIMARY KEY,
    product_id BIGINT,
    barcode VARCHAR NOT NULL UNIQUE,
    created_at TIMESTAMP DEFAULT NOW(),
    CONSTRAINT fk_product_instances_products FOREIGN KEY(product_id) REFERENCES products(id)
);

-- =========================
-- ПОЛЬЗОВАТЕЛИ → ЭТАПЫ
-- =========================
CREATE TABLE user_stages (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT,
    stage_id BIGINT,
    CONSTRAINT fk_user_stages_users FOREIGN KEY(user_id) REFERENCES users(id),
    CONSTRAINT fk_user_stages_stages FOREIGN KEY(stage_id) REFERENCES stages(id)
);

-- =========================
-- РАБОЧИЕ ЗАКАЗЫ
-- =========================
CREATE TABLE work_orders (
    id BIGSERIAL PRIMARY KEY,
    wo_number VARCHAR NOT NULL UNIQUE,
    product_id BIGINT,
    machine_id BIGINT,
    quantity INTEGER NOT NULL,
    priority_id INT DEFAULT 2,
    status_id INT DEFAULT 1,
    deadline DATE,
    created_by BIGINT,
    created_at TIMESTAMP DEFAULT NOW(),
    CONSTRAINT fk_work_orders_products FOREIGN KEY(product_id) REFERENCES products(id),
    CONSTRAINT fk_work_orders_machines FOREIGN KEY(machine_id) REFERENCES machines(id),
    CONSTRAINT fk_work_orders_users FOREIGN KEY(created_by) REFERENCES users(id),
    CONSTRAINT fk_work_orders_priority FOREIGN KEY(priority_id) REFERENCES order_priorities(id),
    CONSTRAINT fk_work_orders_status FOREIGN KEY(status_id) REFERENCES order_statuses(id)
);

-- =========================
-- КОММЕНТАРИИ К ЗАКАЗАМ
-- =========================
CREATE TABLE work_order_comments (
    id BIGSERIAL PRIMARY KEY,
    work_order_id BIGINT,
    user_id BIGINT,
    message TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    CONSTRAINT fk_comments_work_orders FOREIGN KEY(work_order_id) REFERENCES work_orders(id),
    CONSTRAINT fk_comments_users FOREIGN KEY(user_id) REFERENCES users(id)
);

-- =========================
-- РАСПИСАНИЕ
-- =========================
CREATE TABLE schedule (
    id BIGSERIAL PRIMARY KEY,
    work_order_id BIGINT,
    machine_id BIGINT,
    start_time TIMESTAMP NOT NULL,
    end_time TIMESTAMP NOT NULL,
    CONSTRAINT fk_schedule_work_orders FOREIGN KEY(work_order_id) REFERENCES work_orders(id),
    CONSTRAINT fk_schedule_machines FOREIGN KEY(machine_id) REFERENCES machines(id)
);

-- =========================
-- ИНЦИДЕНТЫ
-- =========================
CREATE TABLE incidents (
    id BIGSERIAL PRIMARY KEY,
    machine_id BIGINT,
    work_order_id BIGINT,
    severity_id INT,
    description TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    resolved_at TIMESTAMP,
    CONSTRAINT fk_incidents_machines FOREIGN KEY(machine_id) REFERENCES machines(id),
    CONSTRAINT fk_incidents_work_orders FOREIGN KEY(work_order_id) REFERENCES work_orders(id),
    CONSTRAINT fk_incidents_severity FOREIGN KEY(severity_id) REFERENCES incident_severities(id)
);

-- =========================
-- ЛОГИ ПОЛЬЗОВАТЕЛЕЙ И СИСТЕМЫ
-- =========================
CREATE TABLE user_logs (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT,
    action TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    CONSTRAINT fk_user_logs_users FOREIGN KEY(user_id) REFERENCES users(id)
);

CREATE TABLE system_logs (
    id BIGSERIAL PRIMARY KEY,
    level VARCHAR,
    source VARCHAR,
    message TEXT,
    created_at TIMESTAMP DEFAULT NOW()
);

-- =========================
-- ИСПОЛНЕНИЕ ЭТАПОВ
-- =========================
CREATE TABLE stage_execution (
    id BIGSERIAL PRIMARY KEY,
    product_instance_id BIGINT,
    stage_id BIGINT,
    user_id BIGINT,
    start_time TIMESTAMP NOT NULL,
    end_time TIMESTAMP,
    CONSTRAINT fk_stage_execution_instances FOREIGN KEY(product_instance_id) REFERENCES product_instances(id),
    CONSTRAINT fk_stage_execution_stages FOREIGN KEY(stage_id) REFERENCES stages(id),
    CONSTRAINT fk_stage_execution_users FOREIGN KEY(user_id) REFERENCES users(id)
);

-- =========================
-- ИСТОРИЯ СТАТУСОВ ЗАКАЗОВ
-- =========================
CREATE TABLE work_order_status_history (
    id BIGSERIAL PRIMARY KEY,
    work_order_id BIGINT,
    old_status_id INT,
    new_status_id INT,
    changed_by BIGINT,
    changed_at TIMESTAMP DEFAULT NOW(),
    CONSTRAINT fk_status_history_work_orders FOREIGN KEY(work_order_id) REFERENCES work_orders(id),
    CONSTRAINT fk_status_history_users FOREIGN KEY(changed_by) REFERENCES users(id),
    CONSTRAINT fk_status_history_old_status FOREIGN KEY(old_status_id) REFERENCES order_statuses(id),
    CONSTRAINT fk_status_history_new_status FOREIGN KEY(new_status_id) REFERENCES order_statuses(id)
);

-- =========================
-- КАЧЕСТВО ПРОДУКЦИИ
-- =========================
CREATE TABLE product_quality (
    id BIGSERIAL PRIMARY KEY,
    product_instance_id BIGINT,
    inspected_by BIGINT,
    quality_status_id INT,
    description TEXT,
    inspected_at TIMESTAMP DEFAULT NOW(),
    CONSTRAINT fk_product_quality_instances FOREIGN KEY(product_instance_id) REFERENCES product_instances(id),
    CONSTRAINT fk_product_quality_users FOREIGN KEY(inspected_by) REFERENCES users(id),
    CONSTRAINT fk_product_quality_status FOREIGN KEY(quality_status_id) REFERENCES quality_statuses(id)
);
