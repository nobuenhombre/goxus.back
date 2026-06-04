SELECT
    r.id, r.name, r.slug, r.created_at, r.updated_at
FROM
    public.rbac_roles r
    JOIN public.rbac_user_roles ur ON r.id = ur.role_id
WHERE
    ur.user_id = %%userID int64%%
ORDER BY
    r.id ASC