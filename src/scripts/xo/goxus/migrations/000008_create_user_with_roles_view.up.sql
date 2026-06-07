-- create a view that extends users with concatenated role names
CREATE OR REPLACE VIEW public.user_with_roles AS
SELECT u.id,
       u.name,
       u.email,
       u.password,
       u.email_verified_at,
       u.created_at,
       u.updated_at,
       u.deleted_at,
       COALESCE(
               (SELECT STRING_AGG(r.name, ', ' ORDER BY r.name)
                FROM public.rbac_roles r
                         JOIN public.rbac_user_roles ur ON r.id = ur.role_id
                WHERE ur.user_id = u.id),
               ''
       ) AS roles
FROM public.users u
ORDER BY u.id ASC;
