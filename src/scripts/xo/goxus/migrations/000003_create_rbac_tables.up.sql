CREATE TABLE public.rbac_roles
(
    id         BIGINT GENERATED ALWAYS AS IDENTITY
        CONSTRAINT rbac_roles_pk
            PRIMARY KEY,
    name       VARCHAR(255)            NOT NULL,
    slug       VARCHAR(255)            NOT NULL,
    created_at TIMESTAMP DEFAULT NOW() NOT NULL,
    updated_at TIMESTAMP DEFAULT NOW() NOT NULL
);

CREATE UNIQUE INDEX rbac_roles_slug_uindex
    ON public.rbac_roles (slug);

CREATE TABLE public.rbac_permissions
(
    id         BIGINT GENERATED ALWAYS AS IDENTITY
        CONSTRAINT rbac_permissions_pk
            PRIMARY KEY,
    name       VARCHAR(255)            NOT NULL,
    slug       VARCHAR(255)            NOT NULL,
    created_at TIMESTAMP DEFAULT NOW() NOT NULL,
    updated_at TIMESTAMP DEFAULT NOW() NOT NULL
);

CREATE UNIQUE INDEX rbac_permissions_slug_uindex
    ON public.rbac_permissions (slug);

CREATE TABLE public.rbac_role_permissions
(
    id            BIGINT GENERATED ALWAYS AS IDENTITY
        CONSTRAINT rbac_role_permissions_pk
            PRIMARY KEY,
    role_id       BIGINT                  NOT NULL,
    permission_id BIGINT                  NOT NULL,
    created_at    TIMESTAMP DEFAULT NOW() NOT NULL,
    updated_at    TIMESTAMP DEFAULT NOW() NOT NULL
);

CREATE UNIQUE INDEX rbac_role_permissions_unique
    ON public.rbac_role_permissions (role_id, permission_id);

ALTER TABLE public.rbac_role_permissions
    ADD CONSTRAINT rbac_role_permissions_roles_fk
        FOREIGN KEY (role_id) REFERENCES public.rbac_roles (id)
            ON DELETE CASCADE;

ALTER TABLE public.rbac_role_permissions
    ADD CONSTRAINT rbac_role_permissions_permissions_fk
        FOREIGN KEY (permission_id) REFERENCES public.rbac_permissions (id)
            ON DELETE CASCADE;

CREATE TABLE public.rbac_user_roles
(
    id         BIGINT GENERATED ALWAYS AS IDENTITY
        CONSTRAINT rbac_user_roles_pk
            PRIMARY KEY,
    user_id    BIGINT                  NOT NULL,
    role_id    BIGINT                  NOT NULL,
    created_at TIMESTAMP DEFAULT NOW() NOT NULL,
    updated_at TIMESTAMP DEFAULT NOW() NOT NULL
);

CREATE UNIQUE INDEX rbac_user_roles_unique
    ON public.rbac_user_roles (user_id, role_id);

ALTER TABLE public.rbac_user_roles
    ADD CONSTRAINT rbac_user_roles_roles_fk
        FOREIGN KEY (role_id) REFERENCES public.rbac_roles (id)
            ON DELETE CASCADE;

ALTER TABLE public.rbac_user_roles
    ADD CONSTRAINT rbac_user_roles_users_fk
        FOREIGN KEY (user_id) REFERENCES public.users (id)
            ON DELETE CASCADE;
