CREATE TABLE IF NOT EXISTS users (
                                         id bigserial PRIMARY KEY,
                                         created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
                                         updated_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
                                         fname varchar(255),
                                         sname varchar(255),
                                         email citext UNIQUE NOT NULL,
                                         password_hash bytea NOT NULL,
                                         user_role varchar(50),
                                         activated bool NOT NULL,
                                         version integer NOT NULL DEFAULT 1
);