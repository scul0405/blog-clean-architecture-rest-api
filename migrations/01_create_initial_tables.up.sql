DROP TABLE IF EXISTS users CASCADE;
DROP TABLE IF EXISTS blogs CASCADE;
DROP TABLE IF EXISTS comments CASCADE;

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS CITEXT;

CREATE TABLE users
(
    user_id      UUID PRIMARY KEY                     DEFAULT uuid_generate_v4(),
    first_name   VARCHAR(32)                 NOT NULL CHECK ( first_name <> '' ),
    last_name    VARCHAR(32)                 NOT NULL CHECK ( last_name <> '' ),
    email        VARCHAR(64) UNIQUE          NOT NULL CHECK ( email <> '' ),
    password     VARCHAR(250)                NOT NULL CHECK ( octet_length(password) <> 0 ),
    role         VARCHAR(10)                 NOT NULL DEFAULT 'user',
    about        VARCHAR(1024)                        DEFAULT '',
    avatar       VARCHAR(512),
    phone_number VARCHAR(20),
    address      VARCHAR(250),
    city         VARCHAR(30),
    country      VARCHAR(30),
    gender       VARCHAR(20)                 NOT NULL DEFAULT 'male',
    postcode     INTEGER,
    birthday     DATE                                 DEFAULT NULL,
    created_at   TIMESTAMP WITH TIME ZONE    NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMP WITH TIME ZONE             DEFAULT CURRENT_TIMESTAMP,
    login_date   TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);



CREATE TABLE blogs
(
    blog_id    UUID PRIMARY KEY                  DEFAULT uuid_generate_v4(),
    author_id  UUID                     NOT NULL REFERENCES users (user_id),
    title      VARCHAR(250)             NOT NULL CHECK ( title <> '' ),
    content    TEXT                     NOT NULL CHECK ( content <> '' ),
    image_url  VARCHAR(1024) check ( image_url <> '' ),
    category   VARCHAR(250),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE          DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE comments
(
    comment_id UUID PRIMARY KEY         DEFAULT uuid_generate_v4(),
    author_id  UUID                                               NOT NULL REFERENCES users (user_id) ON DELETE CASCADE,
    blog_id    UUID                                               NOT NULL REFERENCES blogs (blog_id) ON DELETE CASCADE,
    message    VARCHAR(1024)                                      NOT NULL CHECK ( message <> '' ),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS blog_title_id_idx ON blogs (title);
