CREATE TABLE IF NOT EXISTS clientdb (
                                         id bigserial PRIMARY KEY,
                                         fname varchar(255),
                                         sname varchar(255),
                                         email citext UNIQUE NOT NULL,
                                         password_hash bytea NOT NULL,
                                         user_role varchar(50),
                                         activated bool NOT NULL,
                                         version integer NOT NULL DEFAULT 1
    );