-- 社区生鲜团购数据库初始化脚本
-- PostgreSQL 12+

-- 创建数据库
-- CREATE DATABASE fresh_groupbuy;
-- \c fresh_groupbuy;

-- 用户表
CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    phone VARCHAR(20) UNIQUE NOT NULL,
    nickname VARCHAR(50),
    avatar_url VARCHAR(255),
    password_hash VARCHAR(255) NOT NULL,
    role VARCHAR(20) NOT NULL DEFAULT 'member', -- member, leader, admin
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- 团长表
CREATE TABLE IF NOT EXISTS leaders (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT UNIQUE NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    community_name VARCHAR(100) NOT NULL,
    address VARCHAR(255) NOT NULL,
    contact_phone VARCHAR(20) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'pending', -- pending, active, suspended
    audit_note VARCHAR(255),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- 商品表
CREATE TABLE IF NOT EXISTS products (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    image_url VARCHAR(255),
    original_price NUMERIC(10,2) NOT NULL,
    group_price NUMERIC(10,2) NOT NULL,
    stock INT NOT NULL DEFAULT 0,
    group_threshold INT NOT NULL DEFAULT 2, -- 成团所需人数
    category VARCHAR(50),
    leader_id BIGINT REFERENCES leaders(id) ON DELETE SET NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'active', -- active, inactive, sold_out
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- 购物车表
CREATE TABLE IF NOT EXISTS carts (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    product_id BIGINT NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    quantity INT NOT NULL DEFAULT 1,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, product_id)
);

-- 拼团表
CREATE TABLE IF NOT EXISTS group_buys (
    id BIGSERIAL PRIMARY KEY,
    group_code VARCHAR(32) UNIQUE NOT NULL, -- 拼团ID/邀请码
    product_id BIGINT NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    initiator_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    leader_id BIGINT REFERENCES leaders(id) ON DELETE SET NULL,
    current_count INT NOT NULL DEFAULT 1,
    target_count INT NOT NULL, -- 目标成团人数
    status VARCHAR(20) NOT NULL DEFAULT 'pending', -- pending(待成团), success(已成团), failed(失败), cancelled
    expire_at TIMESTAMP NOT NULL, -- 拼团过期时间
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- 拼团成员表
CREATE TABLE IF NOT EXISTS group_buy_members (
    id BIGSERIAL PRIMARY KEY,
    group_buy_id BIGINT NOT NULL REFERENCES group_buys(id) ON DELETE CASCADE,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    order_id BIGINT,
    is_initiator BOOLEAN NOT NULL DEFAULT false,
    joined_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(group_buy_id, user_id)
);

-- 订单表
CREATE TABLE IF NOT EXISTS orders (
    id BIGSERIAL PRIMARY KEY,
    order_no VARCHAR(32) UNIQUE NOT NULL,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    product_id BIGINT NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    group_buy_id BIGINT REFERENCES group_buys(id) ON DELETE SET NULL,
    leader_id BIGINT REFERENCES leaders(id) ON DELETE SET NULL,
    unit_price NUMERIC(10,2) NOT NULL,
    quantity INT NOT NULL DEFAULT 1,
    total_amount NUMERIC(10,2) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'pending_group', -- pending_group(待成团), grouped(已成团), shipped(已发货), delivered(已送达), cancelled
    payment_status VARCHAR(20) NOT NULL DEFAULT 'unpaid', -- unpaid, paid, refunded
    delivery_address VARCHAR(255),
    delivery_phone VARCHAR(20),
    delivery_name VARCHAR(50),
    tracking_no VARCHAR(50),
    remark VARCHAR(255),
    paid_at TIMESTAMP,
    shipped_at TIMESTAMP,
    delivered_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- 索引
CREATE INDEX IF NOT EXISTS idx_orders_user_id ON orders(user_id);
CREATE INDEX IF NOT EXISTS idx_orders_group_buy_id ON orders(group_buy_id);
CREATE INDEX IF NOT EXISTS idx_orders_status ON orders(status);
CREATE INDEX IF NOT EXISTS idx_group_buys_group_code ON group_buys(group_code);
CREATE INDEX IF NOT EXISTS idx_group_buys_status ON group_buys(status);
CREATE INDEX IF NOT EXISTS idx_group_buys_expire_at ON group_buys(expire_at);
CREATE INDEX IF NOT EXISTS idx_products_category ON products(category);
CREATE INDEX IF NOT EXISTS idx_carts_user_id ON carts(user_id);

-- 插入测试数据
INSERT INTO users (phone, nickname, password_hash, role) VALUES
('13800138001', '张团长', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'leader'),
('13800138002', '李用户', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'member'),
('13800138003', '王用户', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'member'),
('13800138004', '赵用户', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'member');

INSERT INTO leaders (user_id, community_name, address, contact_phone, status) VALUES
(1, '阳光小区团', '北京市朝阳区阳光小区1号楼', '13800138001', 'active');

INSERT INTO products (name, description, image_url, original_price, group_price, stock, group_threshold, category, leader_id, status) VALUES
('新鲜草莓', '云南高山草莓，香甜可口', 'https://example.com/strawberry.jpg', 39.90, 29.90, 100, 3, '水果', 1, 'active'),
('有机鸡蛋', '农家散养土鸡蛋，30枚装', 'https://example.com/eggs.jpg', 49.90, 35.90, 200, 2, '禽蛋', 1, 'active'),
('新鲜牛奶', '牧场直供鲜牛奶，2L装', 'https://example.com/milk.jpg', 25.90, 19.90, 150, 2, '乳品', 1, 'active'),
('精选五花肉', '农家散养黑猪五花肉，500g', 'https://example.com/pork.jpg', 45.90, 38.90, 80, 3, '肉类', 1, 'active');
