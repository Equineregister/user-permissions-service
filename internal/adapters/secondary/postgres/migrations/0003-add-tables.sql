-- roles are the supported roles.
CREATE TABLE roles (
    tenant_id UUID NOT NULL,
    role_id UUID PRIMARY KEY,
    role_name TEXT NOT NULL
);
CREATE INDEX idx_roles_tenant_id_role_name ON roles (tenant_id, role_name);

-- Cluster the roles table on the tenant_id and role_name index
CLUSTER roles USING idx_roles_tenant_id_role_name;

-- role_hierarchy is the hierarchy of roles.
-- A role (parent_role_id) can inherit permissions from another role (child_role_id).
-- A role can inherit permissions from multiple roles.
CREATE TABLE role_hierarchy (
    tenant_id UUID NOT NULL,
    parent_role_id UUID NOT NULL,
    child_role_id UUID NOT NULL,
    FOREIGN KEY (parent_role_id) REFERENCES roles(role_id) ON DELETE CASCADE,
    FOREIGN KEY (child_role_id) REFERENCES roles(role_id) ON DELETE CASCADE,
    PRIMARY KEY (parent_role_id, child_role_id),
    CONSTRAINT chk_parent_child_different CHECK (parent_role_id <> child_role_id)
);
CREATE INDEX idx_role_hierarchy_tenant_id_parent_role_id ON role_hierarchy (tenant_id, parent_role_id);
CREATE INDEX idx_role_hierarchy_tenant_id_child_role_id ON role_hierarchy (tenant_id, child_role_id);

-- Cluster the role_hierarchy table on the tenant_id and parent_role_id index
CLUSTER role_hierarchy USING idx_role_hierarchy_tenant_id_parent_role_id;

-- permissions are the supported permissions.
CREATE TABLE permissions (
    tenant_id UUID NOT NULL,
    permission_id UUID PRIMARY KEY,
    permission_name TEXT NOT NULL
);
CREATE INDEX idx_permissions_tenant_id_permission_name ON permissions (tenant_id, permission_name);

-- Cluster the permissions table on the tenant_id and permission_name index
CLUSTER permissions USING idx_permissions_tenant_id_permission_name;

-- tenant_permissions are the permissions that are active for the Tenant.
CREATE TABLE tenant_permissions (
    tenant_id UUID NOT NULL,
    permission_id UUID PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE,
    FOREIGN KEY (permission_id) REFERENCES permissions(permission_id) ON DELETE CASCADE
);
CREATE INDEX idx_tenant_permissions_tenant_id_permission_id ON tenant_permissions (tenant_id, permission_id);

-- Cluster the tenant_permissions table on the tenant_id and permission_id index
CLUSTER tenant_permissions USING idx_tenant_permissions_tenant_id_permission_id;

-- role_permissions are the permissions that are assigned to a Role.
CREATE TABLE role_permissions (
    tenant_id UUID NOT NULL,
    role_id UUID NOT NULL,
    permission_id UUID NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE,
    PRIMARY KEY (role_id, permission_id),
    FOREIGN KEY (role_id) REFERENCES roles(role_id) ON DELETE CASCADE,
    FOREIGN KEY (permission_id) REFERENCES permissions(permission_id) ON DELETE CASCADE
);
CREATE INDEX idx_role_permissions_tenant_id_role_id ON role_permissions (tenant_id, role_id);
CREATE INDEX idx_role_permissions_tenant_id_permission_id ON role_permissions (tenant_id, permission_id);

-- Cluster the role_permissions table on the tenant_id and role_id index
CLUSTER role_permissions USING idx_role_permissions_tenant_id_role_id;

-- resource_types are the supported resource types.
CREATE TABLE resource_types (
    tenant_id UUID NOT NULL,
    resource_type_id BIGINT PRIMARY KEY,
    resource_type_name TEXT NOT NULL UNIQUE
);
CREATE INDEX idx_resource_types_tenant_id_resource_type_name ON resource_types (tenant_id, resource_type_name);

-- Cluster the resource_types table on the tenant_id and resource_type_name index
CLUSTER resource_types USING idx_resource_types_tenant_id_resource_type_name;

-- user_roles are the roles that are assigned to a User.
CREATE TABLE user_roles (
    tenant_id UUID NOT NULL,
    user_roles_id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    user_id UUID NOT NULL,
    role_id UUID NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE,
    FOREIGN KEY (role_id) REFERENCES roles(role_id) ON DELETE CASCADE
);
CREATE INDEX idx_user_roles_tenant_id_user_id ON user_roles (tenant_id, user_id);
CREATE INDEX idx_user_roles_tenant_id_role_id ON user_roles (tenant_id, role_id);

-- Cluster the user_roles table on the tenant_id and user_id index
CLUSTER user_roles USING idx_user_roles_tenant_id_user_id;

CREATE TABLE user_permissions (
    tenant_id UUID NOT NULL,
    user_id UUID NOT NULL,
    permission_id UUID NOT NULL,
    permission_type TEXT NOT NULL CHECK (permission_type IN ('extra', 'revoked')),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE,
    PRIMARY KEY (user_id, permission_id),
    FOREIGN KEY (permission_id) REFERENCES permissions(permission_id) ON DELETE CASCADE
);
CREATE INDEX idx_user_permissions_tenant_id_user_id ON user_permissions (tenant_id, user_id);
CREATE INDEX idx_user_permissions_tenant_id_permission_id ON user_permissions (tenant_id, permission_id);
CREATE INDEX idx_user_permissions_tenant_id_permission_type ON user_permissions (tenant_id, permission_type);

-- Cluster the user_permissions table on the tenant_id and user_id index
CLUSTER user_permissions USING idx_user_permissions_tenant_id_user_id;

-- user_resources are the resources that are assigned to a User and the permission assigned to the User on the resource.
CREATE TABLE user_resources (
    tenant_id UUID NOT NULL,
    user_resources_id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    user_id UUID NOT NULL,
    resource_type_id BIGINT NOT NULL,       
    resource_id UUID NOT NULL,
    permission_id UUID NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE,
    FOREIGN KEY (permission_id) REFERENCES permissions(permission_id) ON DELETE CASCADE,
    FOREIGN KEY (resource_type_id) REFERENCES resource_types(resource_type_id) ON DELETE CASCADE
);
CREATE INDEX idx_user_resources_tenant_id_user_id ON user_resources (tenant_id, user_id);
CREATE INDEX idx_user_resources_tenant_id_resource_id ON user_resources (tenant_id, resource_id);
CREATE INDEX idx_user_resources_tenant_id_permission_id ON user_resources (tenant_id, permission_id);
CREATE INDEX idx_user_resources_tenant_id_resource_type_id ON user_resources (tenant_id, resource_type_id);

-- Cluster the user_resources table on the tenant_id and user_id index
CLUSTER user_resources USING idx_user_resources_tenant_id_user_id;