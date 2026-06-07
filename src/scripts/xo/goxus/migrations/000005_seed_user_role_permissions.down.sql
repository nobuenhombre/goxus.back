-- Rollback user-role permission seed

-- 1. Remove permission links from admin role
DELETE
FROM public.rbac_role_permissions
WHERE permission_id IN (SELECT id
                        FROM public.rbac_permissions
                        WHERE slug IN ('user_role_add', 'user_role_view', 'user_role_delete'));

-- 2. Delete the permissions
DELETE
FROM public.rbac_permissions
WHERE slug IN ('user_role_add', 'user_role_view', 'user_role_delete');
