DO $$ BEGIN
    CREATE TYPE visibility_enum AS ENUM ('private', 'public', 'unlisted');
EXCEPTION
    WHEN duplicate_object THEN null;
END $$;

CREATE TABLE IF NOT EXISTS public.posts (
    id varchar(8) PRIMARY KEY NOT NULL,
    user_id UUID REFERENCES public.users(id) ,
    title varchar(255) NOT NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    password varchar(255),
    has_password BOOLEAN DEFAULT FALSE,
    visibility visibility_enum
);

CREATE INDEX idx_posts_title ON public.posts(title);
CREATE INDEX idx_posts_content ON public.posts(content);
CREATE INDEX idx_posts_created_at_id ON public.posts(created_at DESC, id DESC);
