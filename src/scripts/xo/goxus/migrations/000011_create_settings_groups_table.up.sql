CREATE TABLE public.settings_groups
(
    id          BIGINT GENERATED ALWAYS AS IDENTITY
        CONSTRAINT settings_groups_pk
            PRIMARY KEY,
    name        VARCHAR(255) NOT NULL,
    description TEXT,
    order_pos   BIGINT       NOT NULL DEFAULT 0
);

COMMENT ON TABLE public.settings_groups IS 'Groups of user settings for UI organization';

-- Seed default appearance group
INSERT INTO public.settings_groups (name, description, order_pos)
VALUES ('Appearance', 'Customize the appearance of the app. Automatically switch between day and night themes.', 1)
ON CONFLICT DO NOTHING;
