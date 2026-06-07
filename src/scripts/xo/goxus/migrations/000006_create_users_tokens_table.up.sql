CREATE TABLE public.users_tokens
(
    id           BIGINT GENERATED ALWAYS AS IDENTITY,
    token        VARCHAR(255)            NOT NULL,
    user_id      BIGINT                  NOT NULL,
    last_used_at TIMESTAMP,
    created_at   TIMESTAMP DEFAULT NOW() NOT NULL,
    updated_at   TIMESTAMP DEFAULT NOW() NOT NULL,
    deleted_at   TIMESTAMP
);

-- primary key
ALTER TABLE public.users_tokens
    ADD CONSTRAINT users_tokens_pk PRIMARY KEY (id);

-- foreign key to users
ALTER TABLE public.users_tokens
    ADD CONSTRAINT users_tokens_users_fk FOREIGN KEY (user_id) REFERENCES public.users (id) ON DELETE CASCADE;

-- unique token for fast lookup
CREATE UNIQUE INDEX users_tokens_token_uindex
    ON public.users_tokens (token);

-- index for user token queries
CREATE INDEX users_tokens_user_id_index
    ON public.users_tokens (user_id);
