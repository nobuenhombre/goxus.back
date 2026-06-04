SELECT
    id, role_id, permission_id, created_at, updated_at
FROM
    public.rbac_role_permissions
WHERE
    permission_id = %%permissionID int64%%
ORDER BY
    id ASC