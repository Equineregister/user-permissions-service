-- roles are the supported roles.
CREATE TABLE roles (
    role_id UUID PRIMARY KEY,
    role_name TEXT NOT NULL
);
CREATE INDEX idx_roles_role_name ON roles (role_name);

-- role_hierarchy is the hierarchy of roles.
-- A role can have zero or more multiple child roles.
-- A role can have zero or more parent roles.
CREATE TABLE role_hierarchy (
    parent_role_id UUID,
    child_role_id UUID,
    FOREIGN KEY (parent_role_id) REFERENCES roles(role_id) ON DELETE CASCADE,
    FOREIGN KEY (child_role_id) REFERENCES roles(role_id) ON DELETE CASCADE,
    PRIMARY KEY (parent_role_id, child_role_id),
    CONSTRAINT chk_parent_child_different CHECK (parent_role_id <> child_role_id)
);
CREATE INDEX idx_role_hierarchy_parent_role_id ON role_hierarchy (parent_role_id);
CREATE INDEX idx_role_hierarchy_child_role_id ON role_hierarchy (child_role_id);

-- permissions are the supported permissions.
CREATE TABLE permissions (
    permission_id UUID PRIMARY KEY,
    permission_name TEXT NOT NULL
);
CREATE INDEX idx_permissions_permission_name ON permissions (permission_name);

-- tenant_permissions are the permissions that are active for the Tenant.
CREATE TABLE tenant_permissions (
    permission_id UUID PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE,
    FOREIGN KEY (permission_id) REFERENCES permissions(permission_id) ON DELETE CASCADE
);

-- role_permissions are the permissions that are assigned to a Role.
CREATE TABLE role_permissions (
    role_id UUID NOT NULL,
    permission_id UUID NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE,
    PRIMARY KEY (role_id, permission_id),
    FOREIGN KEY (role_id) REFERENCES roles(role_id) ON DELETE CASCADE,
    FOREIGN KEY (permission_id) REFERENCES permissions(permission_id) ON DELETE CASCADE
);
CREATE INDEX idx_role_permissions_role_id ON role_permissions (role_id);
CREATE INDEX idx_role_permissions_permission_id ON role_permissions (permission_id);

-- resource_types are the supported resource types.
CREATE TABLE resource_types (
    resource_type_id BIGINT PRIMARY KEY,
    resource_type_name TEXT NOT NULL UNIQUE
);

-- user_roles are the roles that are assigned to a User.
CREATE TABLE user_roles (
    user_roles_id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    user_id UUID NOT NULL,
    role_id UUID NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE,
    FOREIGN KEY (role_id) REFERENCES roles(role_id) ON DELETE CASCADE
);
CREATE INDEX idx_user_roles_user_id ON user_roles (user_id);
CREATE INDEX idx_user_roles_role_id ON user_roles (role_id);

CREATE TABLE user_permissions (
    user_id UUID NOT NULL,
    permission_id UUID NOT NULL,
    permission_type TEXT NOT NULL CHECK (permission_type IN ('extra', 'revoked')),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE,
    PRIMARY KEY (user_id, permission_id),
    FOREIGN KEY (permission_id) REFERENCES permissions(permission_id) ON DELETE CASCADE
);
CREATE INDEX idx_user_permissions_user_id ON user_permissions (user_id);
CREATE INDEX idx_user_permissions_permission_id ON user_permissions (permission_id); 
CREATE INDEX idx_user_permissions_permission_type ON user_permissions (permission_type);

-- user_resources are the resources that are assigned to a User.
CREATE TABLE user_resources (
    user_resources_id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    user_id UUID  NOT NULL,
    resource_type_id BIGINT NOT NULL,       
    resource_id UUID  NOT NULL,             -- ID of the resource associated with the user, externally defined.
    permission_id UUID  NOT NULL,           -- What permission the user has on this resource.
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE,
    FOREIGN KEY (permission_id) REFERENCES permissions(permission_id) ON DELETE CASCADE,
    FOREIGN KEY (resource_type_id) REFERENCES resource_types(resource_type_id) ON DELETE CASCADE
);
CREATE INDEX idx_user_resources_user_id ON user_resources (user_id);
CREATE INDEX idx_user_resources_resource_id ON user_resources (resource_id);
CREATE INDEX idx_user_resources_permission_id ON user_resources (permission_id);
CREATE INDEX idx_user_resources_resource_type_id ON user_resources (resource_type_id);