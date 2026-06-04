SELECT
    p.id, p.name, p.slug, p.created_at, p.updated_at
FROM
    public.rbac_permissions p
    JOIN public.rbac_role_permissions rp ON p.id = rp.permission_id
    JOIN public.rbac_roles r ON r.id = rp.role_id
WHERE
    r.slug = %%roleSlug string%%
ORDER BY
    p.id ASC