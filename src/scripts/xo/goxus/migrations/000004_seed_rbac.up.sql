-- Seed RBAC: permissions, admin role, assignment to nobuenhombre@yandex.ru

-- 1. Create permissions
INSERT INTO public.rbac_permissions (name, slug)
VALUES ('User Add', 'user_add'),
       ('User Edit', 'user_edit'),
       ('User Delete', 'user_delete'),
       ('User View', 'user_view')
ON CONFLICT (slug) DO NOTHING;

-- 2. Create admin role
INSERT INTO public.rbac_roles (name, slug)
VALUES ('Admin', 'admin')
ON CONFLICT (slug) DO NOTHING;

-- 3. Link all permissions to admin role
INSERT INTO public.rbac_role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM public.rbac_roles r
         CROSS JOIN public.rbac_permissions p
WHERE r.slug = 'admin'
  AND p.slug IN ('user_add', 'user_edit', 'user_delete', 'user_view')
ON CONFLICT (role_id, permission_id) DO NOTHING;

-- 4. Assign admin role to nobuenhombre@yandex.ru
INSERT INTO public.rbac_user_roles (user_id, role_id)
SELECT u.id, r.id
FROM public.users u
         CROSS JOIN public.rbac_roles r
WHERE u.email = 'nobuenhombre@yandex.ru'
  AND r.slug = 'admin'
ON CONFLICT (user_id, role_id) DO NOTHING;