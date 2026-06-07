CREATE TABLE public.settings
(
    id               BIGINT GENERATED ALWAYS AS IDENTITY
        CONSTRAINT settings_pk
            PRIMARY KEY,
    type_id          BIGINT       NOT NULL
        CONSTRAINT settings_settings_types_id_fk
            REFERENCES public.settings_types
            ON DELETE RESTRICT,
    group_id         BIGINT       NOT NULL
        CONSTRAINT settings_settings_groups_id_fk
            REFERENCES public.settings_groups
            ON DELETE CASCADE,
    name             VARCHAR(255) NOT NULL,
    description      TEXT,
    available_values JSON,
    default_value    JSON
);

COMMENT ON TABLE public.settings IS 'User setting definitions with type, group, available values and defaults';

-- Seed Theme setting
INSERT INTO public.settings (type_id, group_id, name, description, available_values, default_value)
SELECT st.id,
       sg.id,
       'Theme',
       'Select the theme for the dashboard.',
       '{
         "value": {
           "1": "Light",
           "2": "Dark"
         }
       }'::JSON,
       '{
         "value": 1
       }'::JSON
FROM public.settings_types st
         CROSS JOIN public.settings_groups sg
WHERE st.name = 'listRadios'
  AND sg.name = 'Appearance'
ON CONFLICT DO NOTHING;
