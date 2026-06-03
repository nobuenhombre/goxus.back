-- Rollback RBAC seed: remove admin role from user, remove links, delete role and permissions

-- 1. Remove admin role from nobuenhombre@yandex.ru
DELETE FROM public.rbac_user_roles
WHERE user_id = (SELECT id FROM public.users WHERE email = 'nobuenhombre@yandex.ru')
  AND role_id = (SELECT id FROM public.rbac_roles WHERE slug = 'admin');

-- 2. Remove permission links from admin role
DELETE FROM public.rbac_role_permissions
WHERE role_id = (SELECT id FROM public.rbac_roles WHERE slug = 'admin');

-- 3. Delete admin role
DELETE FROM public.rbac_roles WHERE slug = 'admin';

-- 4. Delete permissions
DELETE FROM public.rbac_permissions WHERE slug IN ('user_add', 'user_edit', 'user_delete', 'user_view');