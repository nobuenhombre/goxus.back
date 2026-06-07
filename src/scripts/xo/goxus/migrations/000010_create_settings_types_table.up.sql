CREATE TABLE public.settings_types
(
    id          BIGINT GENERATED ALWAYS AS IDENTITY
        CONSTRAINT settings_types_pk
            PRIMARY KEY,
    name        VARCHAR(255) NOT NULL
        CONSTRAINT settings_types_name_uindex
            UNIQUE,
    description TEXT
);

COMMENT ON TABLE public.settings_types IS 'Types of user settings (input types like text, number, switch, select, etc.)';

-- Seed all setting types from settings.go
INSERT INTO public.settings_types (name, description)
VALUES ('inputTextField', 'Single-line text input field'),
       ('inputPasswordField', 'Password input field'),
       ('inputIntNumberField', 'Integer number input field'),
       ('inputFloatNumberField', 'Floating point number input field'),
       ('textareaField', 'Multi-line text area'),
       ('inputIntSlider', 'Integer slider'),
       ('inputIntSliderRange', 'Integer slider with range (start/end)'),
       ('switch', 'On/off toggle switch'),
       ('listChecks', 'List of checkboxes (multiple selection)'),
       ('listRadios', 'List of radio buttons (single selection)'),
       ('selectSimple', 'Simple dropdown select'),
       ('selectWithSearch', 'Dropdown select with search')
ON CONFLICT (name) DO NOTHING;
