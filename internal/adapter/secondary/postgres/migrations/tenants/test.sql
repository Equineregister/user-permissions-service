-- Insert test data into roles
INSERT INTO roles (role_id, role_name) VALUES
    ('11111111-1111-1111-1111-111111111111', 'admin'),
    ('22222222-2222-2222-2222-222222222222', 'sales person'),
    ('33333333-3333-3333-3333-333333333333', 'sales manager');

-- Insert test data into permissions
INSERT INTO permissions (permission_id, permission_name) VALUES
    ('aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 'invoices:create'),
    ('bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb', 'invoices:read'),
    ('cccccccc-cccc-cccc-cccc-cccccccccccc', 'invoices:delete'),
    ('dddddddd-dddd-dddd-dddd-dddddddddddd', 'products:create'),
    ('eeeeeeee-eeee-eeee-eeee-eeeeeeeeeeee', 'products:read'),
    ('ffffffff-ffff-ffff-ffff-ffffffffffff', 'products:delete'),
    ('gggggggg-gggg-gggg-gggg-gggggggggggg', 'products:update'),
    ('hhhhhhhh-hhhh-hhhh-hhhh-hhhhhhhhhhhh', 'products:disable');


-- Insert test data into tenant_permissions
INSERT INTO tenant_permissions (permission_id, created_at) VALUES
    -- The Tenant has all permissions.
    ('aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', NOW()),
    ('bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb', NOW()),
    ('cccccccc-cccc-cccc-cccc-cccccccccccc', NOW()),
    ('dddddddd-dddd-dddd-dddd-dddddddddddd', NOW()),
    ('eeeeeeee-eeee-eeee-eeee-eeeeeeeeeeee', NOW()),
    ('ffffffff-ffff-ffff-ffff-ffffffffffff', NOW()),
    ('gggggggg-gggg-gggg-gggg-gggggggggggg', NOW()),
    ('hhhhhhhh-hhhh-hhhh-hhhh-hhhhhhhhhhhh', NOW());

-- Insert test data into role_permissions
INSERT INTO role_permissions (role_id, permission_id, created_at) VALUES
    -- The Admin role can do everything.
    ('11111111-1111-1111-1111-111111111111', 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', NOW()),
    ('11111111-1111-1111-1111-111111111111', 'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb', NOW()),
    ('11111111-1111-1111-1111-111111111111', 'cccccccc-cccc-cccc-cccc-cccccccccccc', NOW()),
    ('11111111-1111-1111-1111-111111111111', 'dddddddd-dddd-dddd-dddd-dddddddddddd', NOW()),
    ('11111111-1111-1111-1111-111111111111', 'eeeeeeee-eeee-eeee-eeee-eeeeeeeeeeee', NOW()),
    ('11111111-1111-1111-1111-111111111111', 'ffffffff-ffff-ffff-ffff-ffffffffffff', NOW()),
    -- The Sales Person can create invoices.
    ('22222222-2222-2222-2222-222222222222', 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', NOW()),
    -- The Sales Manager can read and delete invoices and create, read and disable products.
    ('33333333-3333-3333-3333-333333333333', 'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb', NOW()),
    ('33333333-3333-3333-3333-333333333333', 'cccccccc-cccc-cccc-cccc-cccccccccccc', NOW()),
    ('33333333-3333-3333-3333-333333333333', 'eeeeeeee-eeee-eeee-eeee-eeeeeeeeeeee', NOW()),
    ('33333333-3333-3333-3333-333333333333', 'dddddddd-dddd-dddd-dddd-dddddddddddd', NOW()),
    ('33333333-3333-3333-3333-333333333333', 'hhhhhhhh-hhhh-hhhh-hhhh-hhhhhhhhhhhh', NOW());
    
-- Insert test data into resource_types
INSERT INTO resource_types (resource_type_id, resource_type_name) VALUES
    (1, 'invoice'),
    (2, 'product');

-- Insert test data into user_roles
INSERT INTO user_roles (user_id, role_id, created_at) VALUES
    ('032fb302-4aee-4a68-b426-0c6faf12081e', '11111111-1111-1111-1111-111111111111', NOW()), -- A User with the Admin role.
    ('2133479c-35a8-4a49-a682-2952d4772ecc', '22222222-2222-2222-2222-222222222222', NOW()), -- A User with the Sales Person role.
    ('652f4d18-dd3d-40c0-874e-cbe3566abccf', '33333333-3333-3333-3333-333333333333', NOW()); -- A User with the Sales Manager Role.

-- Insert test data into user_resources
INSERT INTO user_resources (user_id, resource_type_id, resource_id, permission_id, created_at) VALUES
    ('2133479c-35a8-4a49-a682-2952d4772ecc', 1, '2f9606d8-4bff-46e7-bd8f-ae9e476d3995', 'cccccccc-cccc-cccc-cccc-cccccccccccc', NOW()), -- The Sales Person can delete this Invoice.
    ('2133479c-35a8-4a49-a682-2952d4772ecc', 1, '568104df-6ff3-40be-b660-91e3160aa7e6', 'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb', NOW()), -- The Sales Person can read this Invoice.
    ('652f4d18-dd3d-40c0-874e-cbe3566abccf', 2, '75248bd5-73a2-4507-9ab3-5418abd33a3c', 'ffffffff-ffff-ffff-ffff-ffffffffffff', NOW()); -- The Sales Manager can delete this Product.

-- Insert test data into user_permissions
INSERT INTO user_permissions (user_id, permission_id, permission_type, created_at) VALUES
    -- The Sales Manager has an extra permission for products:update
    ('652f4d18-dd3d-40c0-874e-cbe3566abccf', 'gggggggg-gggg-gggg-gggg-gggggggggggg', 'extra', NOW()),
    -- The Sales Manager has a revoked permission for products:disable
    ('652f4d18-dd3d-40c0-874e-cbe3566abccf', 'hhhhhhhh-hhhh-hhhh-hhhh-hhhhhhhhhhhh', 'revoked', NOW());
