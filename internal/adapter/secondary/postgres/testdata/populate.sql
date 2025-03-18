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
    ('ffffffff-ffff-ffff-ffff-ffffffffffff', 'products:delete');

-- Insert test data into tenant_permissions
INSERT INTO tenant_permissions (permission_id) VALUES
    -- The Tenant has all permissions.
    ('aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa'),
    ('bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb'),
    ('cccccccc-cccc-cccc-cccc-cccccccccccc'),
    ('dddddddd-dddd-dddd-dddd-dddddddddddd'),
    ('eeeeeeee-eeee-eeee-eeee-eeeeeeeeeeee'),
    ('ffffffff-ffff-ffff-ffff-ffffffffffff');

-- Insert test data into role_permissions
INSERT INTO role_permissions (role_id, permission_id) VALUES
    -- The Admin role can do everything.
    ('11111111-1111-1111-1111-111111111111', 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa'),
    ('11111111-1111-1111-1111-111111111111', 'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb'),
    ('11111111-1111-1111-1111-111111111111', 'cccccccc-cccc-cccc-cccc-cccccccccccc'),
    ('11111111-1111-1111-1111-111111111111', 'dddddddd-dddd-dddd-dddd-dddddddddddd'),
    ('11111111-1111-1111-1111-111111111111', 'eeeeeeee-eeee-eeee-eeee-eeeeeeeeeeee'),
    ('11111111-1111-1111-1111-111111111111', 'ffffffff-ffff-ffff-ffff-ffffffffffff'),
    -- The Sales Person can create invoices.
    ('22222222-2222-2222-2222-222222222222', 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa'),
    -- The Sales Manager can read and delete invoices and create and read products.
    ('33333333-3333-3333-3333-333333333333', 'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb'),
    ('33333333-3333-3333-3333-333333333333', 'cccccccc-cccc-cccc-cccc-cccccccccccc'),
    ('33333333-3333-3333-3333-333333333333', 'eeeeeeee-eeee-eeee-eeee-eeeeeeeeeeee'),
    ('33333333-3333-3333-3333-333333333333', 'dddddddd-dddd-dddd-dddd-dddddddddddd');
    
-- Insert test data into resource_types
INSERT INTO resource_types (resource_type_id, resource_type_name) VALUES
    (1, 'invoice'),
    (2, 'product');

-- Insert test data into user_roles
INSERT INTO user_roles (user_id, role_id) VALUES
    ('032fb302-4aee-4a68-b426-0c6faf12081e', '11111111-1111-1111-1111-111111111111'), -- A User with the Admin role.
    ('2133479c-35a8-4a49-a682-2952d4772ecc', '22222222-2222-2222-2222-222222222222'), -- A User with the Sales Person role.
    ('652f4d18-dd3d-40c0-874e-cbe3566abccf', '33333333-3333-3333-3333-333333333333'); -- A User with the Sales Manager Role.

-- Insert test data into user_resources
INSERT INTO user_resources (user_id, resource_type_id, resource_id, permission_id) VALUES
    ('2133479c-35a8-4a49-a682-2952d4772ecc', 1, '2f9606d8-4bff-46e7-bd8f-ae9e476d3995', 'cccccccc-cccc-cccc-cccc-cccccccccccc'), -- The Sales Person can delete this Invoice.
    ('2133479c-35a8-4a49-a682-2952d4772ecc', 1, '568104df-6ff3-40be-b660-91e3160aa7e6', 'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb'), -- The Sales Person can read this Invoice.
    ('652f4d18-dd3d-40c0-874e-cbe3566abccf', 2, '75248bd5-73a2-4507-9ab3-5418abd33a3c', 'ffffffff-ffff-ffff-ffff-ffffffffffff'); -- The Sales Manager can delete this Product.