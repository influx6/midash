use midash;

create table users (
    hash varchar(255) NOT NULL,
    email varchar(255) NOT NULL,
    public_id varchar(255) PRIMARY KEY NOT NULL,
    private_id varchar(255) NOT NULL,
    updated_at timestamp NOT NULL,
    created_at timestamp NOT NULL
);

create table profiles (
    address text NOT NULL,
    user_public_id varchar(255) NOT NULL,
    public_id varchar(255) PRIMARY KEY NOT NULL,
    first_name varchar(255) NOT NULL,
    last_name varchar(255) NOT NULL,
    updated_at timestamp NOT NULL,
    created_at timestamp NOT NULL,
    INDEX user_id (user_public_id)
);

create table sessions (
    user_public_id varchar(255) NOT NULL,
    public_id varchar(255) PRIMARY KEY NOT NULL,
    token varchar(255) NOT NULL,
    expiration timestamp NOT NULL,
    created_at timestamp NOT NULL,
    INDEX user_id (user_public_id)
);