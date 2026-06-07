-- Seed data permissions, roles, user

-- 1. Create permissions
INSERT INTO public.rbac_permissions (name, slug)
VALUES ('Data View', 'data_view'),
       ('Data Add', 'data_add'),
       ('Data Edit', 'data_edit'),
       ('Data Delete', 'data_delete')
ON CONFLICT (slug) DO NOTHING;

-- 2. Create roles
INSERT INTO public.rbac_roles (name, slug)
VALUES ('Data Analytics', 'data_analytics'),
       ('Data Operator', 'data_operator')
ON CONFLICT (slug) DO NOTHING;

-- 3. Link data permissions to admin role
INSERT INTO public.rbac_role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM public.rbac_roles r
         CROSS JOIN public.rbac_permissions p
WHERE r.slug = 'admin'
  AND p.slug IN ('data_view', 'data_add', 'data_edit', 'data_delete')
ON CONFLICT (role_id, permission_id) DO NOTHING;

-- 4. Link data_view to data_analytics role
INSERT INTO public.rbac_role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM public.rbac_roles r
         CROSS JOIN public.rbac_permissions p
WHERE r.slug = 'data_analytics'
  AND p.slug IN ('data_view')
ON CONFLICT (role_id, permission_id) DO NOTHING;

-- 5. Link data_add, data_edit, data_delete to data_operator role
INSERT INTO public.rbac_role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM public.rbac_roles r
         CROSS JOIN public.rbac_permissions p
WHERE r.slug = 'data_operator'
  AND p.slug IN ('data_add', 'data_edit', 'data_delete')
ON CONFLICT (role_id, permission_id) DO NOTHING;

-- 6. Create user Ivan Data Worker
INSERT INTO public.users (name, email, password, email_verified_at)
VALUES ('Ivan Data Worker', 'data.worker.nobuenhombre@yandex.ru', '123', NOW())
ON CONFLICT (email) DO NOTHING;

-- 7. Assign data_operator and data_analytics roles to Ivan Data Worker
INSERT INTO public.rbac_user_roles (user_id, role_id)
SELECT u.id, r.id
FROM public.users u
         CROSS JOIN public.rbac_roles r
WHERE u.email = 'data.worker.nobuenhombre@yandex.ru'
  AND r.slug IN ('data_operator', 'data_analytics')
ON CONFLICT (user_id, role_id) DO NOTHING;
