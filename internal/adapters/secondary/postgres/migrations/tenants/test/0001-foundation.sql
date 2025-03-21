-- Insert test data into roles
INSERT INTO roles (tenant_id, role_id, role_name) VALUES
    ('639a4003-b342-4e17-8aac-a8d1bdd2c8e3', '550e8400-e29b-41d4-a716-446655440000', 'admin'),
    ('639a4003-b342-4e17-8aac-a8d1bdd2c8e3', 'da244750-f014-415c-b7b9-43ead3d8fa25', 'sales auditor'),
    ('639a4003-b342-4e17-8aac-a8d1bdd2c8e3', '123e4567-e89b-12d3-a456-426614174000', 'sales person'),
    ('639a4003-b342-4e17-8aac-a8d1bdd2c8e3', 'f47ac10b-58cc-4372-a567-0e02b2c3d479', 'sales manager');

-- Insert test data into role_hierarchy
INSERT INTO role_hierarchy (tenant_id, parent_role_id, child_role_id) VALUES
    -- Sales Manager inherits from Sales Person
    ('639a4003-b342-4e17-8aac-a8d1bdd2c8e3', 'f47ac10b-58cc-4372-a567-0e02b2c3d479', '123e4567-e89b-12d3-a456-426614174000'),
    -- Sales Person inherits from Sales Auditor
    ('639a4003-b342-4e17-8aac-a8d1bdd2c8e3', '123e4567-e89b-12d3-a456-426614174000', 'da244750-f014-415c-b7b9-43ead3d8fa25');

-- Insert test data into permissions
INSERT INTO permissions (tenant_id, permission_id, permission_name) VALUES
    ('639a4003-b342-4e17-8aac-a8d1bdd2c8e3', '2f9606d8-4bff-46e7-bd8f-ae9e476d3995', 'invoices:create'),
    ('639a4003-b342-4e17-8aac-a8d1bdd2c8e3', '8f20eca6-9859-4532-babb-65a528e1611e', 'invoices:read'),
    ('639a4003-b342-4e17-8aac-a8d1bdd2c8e3', '41c21275-b7d5-4031-b551-b5e293b85319', 'invoices:delete'),
    ('639a4003-b342-4e17-8aac-a8d1bdd2c8e3', 'df6ae9bc-e957-41c1-a683-3773667c7628', 'products:create'),
    ('639a4003-b342-4e17-8aac-a8d1bdd2c8e3', '62752f21-fbe2-4301-a72d-7dc8963e08e2', 'products:read'),
    ('639a4003-b342-4e17-8aac-a8d1bdd2c8e3', 'acecdadf-f527-45bf-8123-353b7ee8dc6a', 'products:delete'),
    ('639a4003-b342-4e17-8aac-a8d1bdd2c8e3', 'e12d692b-3a96-43aa-a966-dd3add99d312', 'products:update'),
    ('639a4003-b342-4e17-8aac-a8d1bdd2c8e3', 'cf7dc325-6bc9-44f5-aafb-fcdc694b111d', 'products:disable');

-- Insert test data into tenant_permissions
INSERT INTO tenant_permissions (tenant_id, permission_id, created_at) VALUES
    ('639a4003-b342-4e17-8aac-a8d1bdd2c8e3', '2f9606d8-4bff-46e7-bd8f-ae9e476d3995', NOW()),
    ('639a4003-b342-4e17-8aac-a8d1bdd2c8e3', '8f20eca6-9859-4532-babb-65a528e1611e', NOW()),
    ('639a4003-b342-4e17-8aac-a8d1bdd2c8e3', '41c21275-b7d5-4031-b551-b5e293b85319', NOW()),
    ('639a4003-b342-4e17-8aac-a8d1bdd2c8e3', 'df6ae9bc-e957-41c1-a683-3773667c7628', NOW()),
    ('639a4003-b342-4e17-8aac-a8d1bdd2c8e3', '62752f21-fbe2-4301-a72d-7dc8963e08e2', NOW()),
    ('639a4003-b342-4e17-8aac-a8d1bdd2c8e3', 'acecdadf-f527-45bf-8123-353b7ee8dc6a', NOW()),
    ('639a4003-b342-4e17-8aac-a8d1bdd2c8e3', 'e12d692b-3a96-43aa-a966-dd3add99d312', NOW()),
    ('639a4003-b342-4e17-8aac-a8d1bdd2c8e3', 'cf7dc325-6bc9-44f5-aafb-fcdc694b111d', NOW());

-- Insert test data into role_permissions
INSERT INTO role_permissions (tenant_id, role_id, permission_id, created_at) VALUES
    -- The Admin role can do everything.
    ('639a4003-b342-4e17-8aac-a8d1bdd2c8e3', '550e8400-e29b-41d4-a716-446655440000', '2f9606d8-4bff-46e7-bd8f-ae9e476d3995', NOW()),
    ('639a4003-b342-4e17-8aac-a8d1bdd2c8e3', '550e8400-e29b-41d4-a716-446655440000', '8f20eca6-9859-4532-babb-65a528e1611e', NOW()),
    ('639a4003-b342-4e17-8aac-a8d1bdd2c8e3', '550e8400-e29b-41d4-a716-446655440000', '41c21275-b7d5-4031-b551-b5e293b85319', NOW()),
    ('639a4003-b342-4e17-8aac-a8d1bdd2c8e3', '550e8400-e29b-41d4-a716-446655440000', 'df6ae9bc-e957-41c1-a683-3773667c7628', NOW()),
    ('639a4003-b342-4e17-8aac-a8d1bdd2c8e3', '550e8400-e29b-41d4-a716-446655440000', '62752f21-fbe2-4301-a72d-7dc8963e08e2', NOW()),
    ('639a4003-b342-4e17-8aac-a8d1bdd2c8e3', '550e8400-e29b-41d4-a716-446655440000', 'acecdadf-f527-45bf-8123-353b7ee8dc6a', NOW()),
    -- The Sales Auditor can read invoices.
    ('639a4003-b342-4e17-8aac-a8d1bdd2c8e3', 'da244750-f014-415c-b7b9-43ead3d8fa25', '8f20eca6-9859-4532-babb-65a528e1611e', NOW()),
    -- The Sales Person can create invoices.
    -- They can also read invoices, due to role inheritance from Sales Auditor.
    ('639a4003-b342-4e17-8aac-a8d1bdd2c8e3', '123e4567-e89b-12d3-a456-426614174000', '2f9606d8-4bff-46e7-bd8f-ae9e476d3995', NOW()),
    -- The Sales Manager can delete invoices and create, read and disable products. 
    -- They can also create invoices, due to role inheritance from Sales Person.
    -- They can also read invoices, due to role inheritance from Sales Auditor -> Sales Person.
    ('639a4003-b342-4e17-8aac-a8d1bdd2c8e3', 'f47ac10b-58cc-4372-a567-0e02b2c3d479', '41c21275-b7d5-4031-b551-b5e293b85319', NOW()),
    ('639a4003-b342-4e17-8aac-a8d1bdd2c8e3', 'f47ac10b-58cc-4372-a567-0e02b2c3d479', '62752f21-fbe2-4301-a72d-7dc8963e08e2', NOW()),
    ('639a4003-b342-4e17-8aac-a8d1bdd2c8e3', 'f47ac10b-58cc-4372-a567-0e02b2c3d479', 'df6ae9bc-e957-41c1-a683-3773667c7628', NOW()),
    ('639a4003-b342-4e17-8aac-a8d1bdd2c8e3', 'f47ac10b-58cc-4372-a567-0e02b2c3d479', 'cf7dc325-6bc9-44f5-aafb-fcdc694b111d', NOW());
    
-- Insert test data into resource_types
INSERT INTO resource_types (tenant_id, resource_type_id, resource_type_name) VALUES
    ('639a4003-b342-4e17-8aac-a8d1bdd2c8e3', 1, 'invoices'),
    ('639a4003-b342-4e17-8aac-a8d1bdd2c8e3', 2, 'products');

-- Insert test data into user_roles
INSERT INTO user_roles (tenant_id, user_id, role_id, created_at) VALUES
    ('639a4003-b342-4e17-8aac-a8d1bdd2c8e3', '032fb302-4aee-4a68-b426-0c6faf12081e', '550e8400-e29b-41d4-a716-446655440000', NOW()),
    ('639a4003-b342-4e17-8aac-a8d1bdd2c8e3', '2133479c-35a8-4a49-a682-2952d4772ecc', '123e4567-e89b-12d3-a456-426614174000', NOW()),
    ('639a4003-b342-4e17-8aac-a8d1bdd2c8e3', '652f4d18-dd3d-40c0-874e-cbe3566abccf', 'f47ac10b-58cc-4372-a567-0e02b2c3d479', NOW());

-- Insert test data into user_resources
INSERT INTO user_resources (tenant_id, user_id, resource_type_id, resource_id, permission_id, created_at) VALUES
    ('639a4003-b342-4e17-8aac-a8d1bdd2c8e3', '2133479c-35a8-4a49-a682-2952d4772ecc', 1, '6b63b489-61cb-4087-8636-f10716bd724e', '41c21275-b7d5-4031-b551-b5e293b85319', NOW()),
    ('639a4003-b342-4e17-8aac-a8d1bdd2c8e3', '2133479c-35a8-4a49-a682-2952d4772ecc', 1, '568104df-6ff3-40be-b660-91e3160aa7e6', '8f20eca6-9859-4532-babb-65a528e1611e', NOW()),
    ('639a4003-b342-4e17-8aac-a8d1bdd2c8e3', '652f4d18-dd3d-40c0-874e-cbe3566abccf', 2, '75248bd5-73a2-4507-9ab3-5418abd33a3c', 'acecdadf-f527-45bf-8123-353b7ee8dc6a', NOW());

-- Insert test data into user_permissions
INSERT INTO user_permissions (tenant_id, user_id, permission_id, permission_type, created_at) VALUES
    -- The Sales Manager has an extra permission for products:update
    ('639a4003-b342-4e17-8aac-a8d1bdd2c8e3', '652f4d18-dd3d-40c0-874e-cbe3566abccf', 'e12d692b-3a96-43aa-a966-dd3add99d312', 'extra', NOW()),
    -- The Sales Manager has a revoked permission for products:disable
    ('639a4003-b342-4e17-8aac-a8d1bdd2c8e3', '652f4d18-dd3d-40c0-874e-cbe3566abccf', 'cf7dc325-6bc9-44f5-aafb-fcdc694b111d', 'revoked', NOW());
