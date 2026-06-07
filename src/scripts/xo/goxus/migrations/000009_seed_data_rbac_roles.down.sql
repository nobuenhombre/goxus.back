-- Rollback data RBAC seed: remove roles from user, remove user, remove links, remove roles and permissions

-- 1. Remove roles from Ivan Data Worker
DELETE
FROM public.rbac_user_roles
WHERE user_id = (SELECT id FROM public.users WHERE email = 'data.worker.nobuenhombre@yandex.ru')
  AND role_id IN (SELECT id FROM public.rbac_roles WHERE slug IN ('data_operator', 'data_analytics'));

-- 2. Delete Ivan Data Worker
DELETE
FROM public.users
WHERE email = 'data.worker.nobuenhombre@yandex.ru';

-- 3. Remove permission links from data_analytics role
DELETE
FROM public.rbac_role_permissions
WHERE role_id = (SELECT id FROM public.rbac_roles WHERE slug = 'data_analytics');

-- 4. Remove permission links from data_operator role
DELETE
FROM public.rbac_role_permissions
WHERE role_id = (SELECT id FROM public.rbac_roles WHERE slug = 'data_operator');

-- 5. Remove data permission links from admin role
DELETE
FROM public.rbac_role_permissions
WHERE permission_id IN (SELECT id
                        FROM public.rbac_permissions
                        WHERE slug IN ('data_view', 'data_add', 'data_edit', 'data_delete'));

-- 6. Delete roles
DELETE
FROM public.rbac_roles
WHERE slug IN ('data_analytics', 'data_operator');

-- 7. Delete permissions
DELETE
FROM public.rbac_permissions
WHERE slug IN ('data_view', 'data_add', 'data_edit', 'data_delete');
