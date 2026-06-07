CREATE TABLE public.users_settings
(
    id          BIGINT GENERATED ALWAYS AS IDENTITY
        CONSTRAINT users_settings_pk
            PRIMARY KEY,
    user_id     BIGINT NOT NULL
        CONSTRAINT users_settings_users_id_fk
            REFERENCES public.users
            ON DELETE CASCADE,
    settings_id BIGINT NOT NULL,
    value       JSON
);

COMMENT ON TABLE public.users_settings IS 'User-specific setting values';

-- Seed Theme setting for both users with default value (Light)
INSERT INTO public.users_settings (user_id, settings_id, value)
SELECT u.id, s.id, '{
  "value": 1
}'::JSON
FROM public.users u
         CROSS JOIN public.settings s
         JOIN public.settings_groups sg ON s.group_id = sg.id
         JOIN public.settings_types st ON s.type_id = st.id
WHERE u.email IN ('nobuenhombre@yandex.ru', 'data.worker.nobuenhombre@yandex.ru')
  AND st.name = 'listRadios'
  AND sg.name = 'Appearance'
  AND s.name = 'Theme'
ON CONFLICT DO NOTHING;
