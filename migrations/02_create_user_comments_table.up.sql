CREATE TABLE user_comments
(
    user_id    UUID                                               NOT NULL REFERENCES users (user_id) ON DELETE CASCADE,
    comment_id UUID                                               NOT NULL REFERENCES comments (comment_id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT  user_comments_pkey PRIMARY KEY (user_id, comment_id)
);