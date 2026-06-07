CREATE UNIQUE INDEX IF NOT EXISTS users_settings_user_id_settings_id_idx
    ON public.users_settings (user_id, settings_id);

COMMENT ON INDEX public.users_settings_user_id_settings_id_idx IS 'Index for fast lookup of user settings by user and setting';
