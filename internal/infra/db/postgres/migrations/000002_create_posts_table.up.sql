CREATE TABLE IF NOT EXISTS public.posts (
    id UUID PRIMARY KEY NOT NULL,
    user_id UUID REFERENCES public.users(id) ,
    title varchar(255) NOT NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    password varchar(255),
    is_private BOOLEAN DEFAULT FALSE
);

CREATE INDEX idx_posts_title ON public.posts(title);
CREATE INDEX idx_posts_content ON public.posts(content);
CREATE INDEX idx_posts_created_at_id ON public.posts(created_at DESC, id DESC);
