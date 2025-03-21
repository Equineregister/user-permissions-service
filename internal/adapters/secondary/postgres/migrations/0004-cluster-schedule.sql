-- SKIP_WHEN_TESTING

-- Enable the pg_cron extension
CREATE EXTENSION IF NOT EXISTS pg_cron;

-- Schedule a job to cluster the tables every day at midnight
SELECT cron.schedule('cluster_roles', '0 0 * * *', $$CLUSTER roles USING idx_roles_tenant_id_role_name$$);
SELECT cron.schedule('cluster_role_hierarchy', '0 0 * * *', $$CLUSTER role_hierarchy USING idx_role_hierarchy_tenant_id_parent_role_id$$);
SELECT cron.schedule('cluster_permissions', '0 0 * * *', $$CLUSTER permissions USING idx_permissions_tenant_id_permission_name$$);
SELECT cron.schedule('cluster_tenant_permissions', '0 0 * * *', $$CLUSTER tenant_permissions USING idx_tenant_permissions_tenant_id_permission_id$$);
SELECT cron.schedule('cluster_role_permissions', '0 0 * * *', $$CLUSTER role_permissions USING idx_role_permissions_tenant_id_role_id$$);
SELECT cron.schedule('cluster_resource_types', '0 0 * * *', $$CLUSTER resource_types USING idx_resource_types_tenant_id_resource_type_name$$);
SELECT cron.schedule('cluster_user_roles', '0 0 * * *', $$CLUSTER user_roles USING idx_user_roles_tenant_id_user_id$$);
SELECT cron.schedule('cluster_user_permissions', '0 0 * * *', $$CLUSTER user_permissions USING idx_user_permissions_tenant_id_user_id$$);
SELECT cron.schedule('cluster_user_resources', '0 0 * * *', $$CLUSTER user_resources USING idx_user_resources_tenant_id_user_id$$);
