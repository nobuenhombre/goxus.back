CREATE TABLE public.users
(
    id                BIGINT GENERATED ALWAYS AS IDENTITY
        CONSTRAINT users_pk
            PRIMARY KEY,
    name              VARCHAR(255)            NOT NULL,
    email             VARCHAR(255)            NOT NULL,
    password          VARCHAR(255)            NOT NULL,
    email_verified_at TIMESTAMP,
    created_at        TIMESTAMP DEFAULT NOW() NOT NULL,
    updated_at        TIMESTAMP DEFAULT NOW() NOT NULL,
    deleted_at        TIMESTAMP
);

CREATE UNIQUE INDEX users_email_uindex
    ON public.users (email);
