CREATE TABLE IF NOT EXISTS categories(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    slug VARCHAR(100) UNIQUE NOT NULL,
    parent_id UUID REFERENCES categories(id),
    image_url TEXT,
    created_at timestamp(0) with time zone NOT NULL DEFAULt now()
 )