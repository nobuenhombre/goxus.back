SELECT
    p.id, p.name, p.slug, p.created_at, p.updated_at
FROM
    public.rbac_permissions p
    JOIN public.rbac_role_permissions rp ON p.id = rp.permission_id
    JOIN public.rbac_user_roles ur ON rp.role_id = ur.role_id
WHERE
    ur.user_id = %%userID int64%%
  AND p.slug = %%permSlug string%%
ORDER BY
    p.id ASC