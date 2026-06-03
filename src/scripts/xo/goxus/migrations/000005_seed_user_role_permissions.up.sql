-- Seed user-role management permissions

-- 1. Create permissions
INSERT INTO public.rbac_permissions (name, slug)
VALUES ('User Role Add', 'user_role_add'),
       ('User Role View', 'user_role_view'),
       ('User Role Delete', 'user_role_delete')
ON CONFLICT (slug) DO NOTHING;

-- 2. Link new permissions to admin role
INSERT INTO public.rbac_role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM public.rbac_roles r
         CROSS JOIN public.rbac_permissions p
WHERE r.slug = 'admin'
  AND p.slug IN ('user_role_add', 'user_role_view', 'user_role_delete')
ON CONFLICT (role_id, permission_id) DO NOTHING;