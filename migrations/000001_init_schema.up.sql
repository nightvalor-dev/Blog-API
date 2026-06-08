CREATE TABLE users (
                       user_id       INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
                       username      VARCHAR(50)  NOT NULL,
                       email         VARCHAR(50)  UNIQUE NOT NULL,
                       phone         VARCHAR(20),
                       bio           VARCHAR(100),
                       password_hash VARCHAR(255) NOT NULL,
                       role          VARCHAR(20)  NOT NULL DEFAULT 'user',
                       created_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
                       updated_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE TABLE categories (
                            category_id   INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
                            category_name VARCHAR(100) NOT NULL UNIQUE,
                            description   TEXT,
                            created_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE TABLE tags (
                      tag_id     INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
                      tag_name   VARCHAR(100) NOT NULL UNIQUE,
                      created_at TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE TABLE blogs (
                       blog_id      INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
                       title        VARCHAR(255) NOT NULL,
                       content      TEXT         NOT NULL,
                       status       VARCHAR(20)  NOT NULL DEFAULT 'drafted'
                           CHECK (status IN ('drafted', 'published', 'archived')),
                       user_id      INT          NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
                       created_at   TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
                       updated_at   TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
                       published_at TIMESTAMPTZ
);

CREATE TABLE blog_categories (
                                 blog_id     INT NOT NULL REFERENCES blogs(blog_id)      ON DELETE CASCADE,
                                 category_id INT NOT NULL REFERENCES categories(category_id) ON DELETE CASCADE,
                                 PRIMARY KEY (blog_id, category_id)
);

CREATE TABLE blog_tags (
                           blog_id INT NOT NULL REFERENCES blogs(blog_id) ON DELETE CASCADE,
                           tag_id  INT NOT NULL REFERENCES tags(tag_id)   ON DELETE CASCADE,
                           PRIMARY KEY (blog_id, tag_id)
);

CREATE TABLE blog_likes (
                            blog_id INT REFERENCES blogs(blog_id) ON DELETE CASCADE,
                            user_id INT REFERENCES users(user_id) ON DELETE CASCADE,
                            created_at TIMESTAMPTZ DEFAULT NOW(),
                            PRIMARY KEY (blog_id, user_id)
);

CREATE TABLE comments (
                          comment_id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
                          content    TEXT        NOT NULL,
                          blog_id    INT         NOT NULL REFERENCES blogs(blog_id) ON DELETE CASCADE,
                          user_id    INT         NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
                          created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
                          updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE media (
                       media_id   SERIAL PRIMARY KEY,
                       blog_id    INT REFERENCES blogs(blog_id) ON DELETE CASCADE,
                       url        TEXT NOT NULL,
                       public_id  TEXT NOT NULL DEFAULT '',
                       created_at TIMESTAMPTZ DEFAULT NOW()
);