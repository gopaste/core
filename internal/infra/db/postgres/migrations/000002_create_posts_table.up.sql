CREATE TABLE IF NOT EXISTS public.posts (
    post_id UUID PRIMARY KEY NOT NULL,
    user_id UUID REFERENCES public.users(id) ,
    title varchar(255) NOT NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
