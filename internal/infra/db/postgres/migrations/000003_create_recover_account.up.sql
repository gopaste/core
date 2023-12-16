CREATE TABLE IF NOT EXISTS public.password_reset (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    reset_token VARCHAR(255) NOT NULL,
    expiration_datetime TIMESTAMP NOT NULL,
    FOREIGN KEY (user_id) REFERENCES public.users(id)
);
